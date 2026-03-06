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

func (db *SimpleDB) Delete(key string) (err error) {
	var (
		current entry
		exists  bool
		record  logRecord
	)

	if err = validateKey(key); err != nil {
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if err = db.ensureOpen(); err != nil {
		return
	}

	if current, exists = db.index[key]; !exists {
		return ErrKeyNotFound
	}

	if current.Deleted {
		return ErrKeyDeleted
	}

	record = logRecord{
		Operation: opDelete,
		Key:       key,
		CreatedAt: time.Now().UnixNano(),
	}

	if err = db.appendRecord(record); err != nil {
		return
	}

	current.Value = nil
	current.Deleted = true
	current.UpdatedAt = record.CreatedAt
	db.index[key] = current

	return
}

func (db *SimpleDB) appendRecord(record logRecord) (err error) {
	if _, err = db.file.Seek(0, io.SeekEnd); err != nil {
		return
	}
	if err = writeRecord(db.file, record); err != nil {
		return
	}

	return db.file.Sync()
}

func (db *SimpleDB) Update(key string, value []byte) (err error) {
	var (
		current entry
		exists  bool
	)

	if err = validateKey(key); err != nil {
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if err = db.ensureOpen(); err != nil {
		return
	}

	if current, exists = db.index[key]; !exists {
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

	if err = db.appendRecord(record); err != nil {
		return
	}

	db.index[key] = entry{Value: cloneBytes(value), UpdatedAt: record.CreatedAt}

	return nil
}

func (db *SimpleDB) ensureOpen() (err error) {
	if db.closed {
		return ErrDatabaseClosed
	}
	return
}

func (db *SimpleDB) load() (err error) {
	var (
		reader *bufio.Reader
		record logRecord
	)

	if _, err = db.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	reader = bufio.NewReader(db.file)
	for {
		if record, err = readRecord(reader); errors.Is(err, io.EOF) {
			return
		}
		if err != nil {
			return
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
			return fmt.Errorf("%w：%d", ErrCorruptedRecord, record.Operation)
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
