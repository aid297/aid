package simpleDBDriver

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"time"
)

func (db *SimpleDB) Put(key string, value []byte) (err error) {
	var record logRecord

	if err = validateKey(key); err != nil {
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if err = db.ensureOpen(); err != nil {
		return
	}

	record = logRecord{
		Operation: opPut,
		Key:       key,
		Value:     cloneBytes(value),
		CreatedAt: time.Now().UnixNano(),
	}

	if err = db.appendRecord(record); err != nil {
		return
	}

	db.index[key] = entry{Value: cloneBytes(value), UpdatedAt: record.CreatedAt}

	return
}

func (db *SimpleDB) Delete(key string) error {
	if err := validateKey(key); err != nil {
		return err
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if err := db.ensureOpen(); err != nil {
		return err
	}

	current, exists := db.index[key]
	if !exists {
		return ErrKeyNotFound
	}
	if current.Deleted {
		return ErrKeyDeleted
	}

	record := logRecord{
		Operation: opDelete,
		Key:       key,
		CreatedAt: time.Now().UnixNano(),
	}

	if err := db.appendRecord(record); err != nil {
		return err
	}

	current.Value = nil
	current.Deleted = true
	current.UpdatedAt = record.CreatedAt
	db.index[key] = current
	return nil
}

func (db *SimpleDB) appendRecord(record logRecord) error {
	if _, err := db.file.Seek(0, io.SeekEnd); err != nil {
		return err
	}
	if err := writeRecord(db.file, record); err != nil {
		return err
	}
	return db.file.Sync()
}

func (db *SimpleDB) Update(key string, value []byte) error {
	if err := validateKey(key); err != nil {
		return err
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if err := db.ensureOpen(); err != nil {
		return err
	}

	current, exists := db.index[key]
	if !exists {
		return ErrKeyNotFound
	}
	if current.Deleted {
		return ErrKeyDeleted
	}

	record := logRecord{
		Operation: opPut,
		Key:       key,
		Value:     cloneBytes(value),
		CreatedAt: time.Now().UnixNano(),
	}

	if err := db.appendRecord(record); err != nil {
		return err
	}

	db.index[key] = entry{Value: cloneBytes(value), UpdatedAt: record.CreatedAt}
	return nil
}

func (db *SimpleDB) ensureOpen() error {
	if db.closed {
		return ErrDatabaseClosed
	}
	return nil
}

func (db *SimpleDB) load() error {
	if _, err := db.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	reader := bufio.NewReader(db.file)
	for {
		record, err := readRecord(reader)
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}

		switch record.Operation {
		case opPut:
			db.index[record.Key] = entry{
				Value:     cloneBytes(record.Value),
				UpdatedAt: record.CreatedAt,
			}
		case opDelete:
			current := db.index[record.Key]
			current.Value = nil
			current.Deleted = true
			current.UpdatedAt = record.CreatedAt
			db.index[record.Key] = current
		default:
			return fmt.Errorf("%w: unknown operation %d", ErrCorruptedRecord, record.Operation)
		}
	}
}

func writeRecord(writer io.Writer, record logRecord) error {
	payload, err := json.Marshal(record)
	if err != nil {
		return err
	}
	if len(payload) > maxRecordSize {
		return fmt.Errorf("simpleDB: record too large: %d", len(payload))
	}

	var header [4]byte
	binary.BigEndian.PutUint32(header[:], uint32(len(payload)))
	if _, err = writer.Write(header[:]); err != nil {
		return err
	}
	if _, err = writer.Write(payload); err != nil {
		return err
	}

	var checksum [4]byte
	binary.BigEndian.PutUint32(checksum[:], crc32.ChecksumIEEE(payload))
	if _, err = writer.Write(checksum[:]); err != nil {
		return err
	}
	return nil
}
