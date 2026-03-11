package kernal

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"time"

	"github.com/aid297/aid/operation/operationV2"
)

func (db *SimpleDB) Put(key string, value []byte) (err error) {
	if err = validateKey(key); err != nil {
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if err = db.ensureOpen(); err != nil {
		return
	}

	return db.putRawLocked(key, value)
}

func (db *SimpleDB) Delete(key string) (err error) {
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

	_ = current
	return db.deleteRawLocked(key)
}

func (db *SimpleDB) appendRecord(record logRecord) (err error) {
	if err = writeRecord(db.file, record); err != nil {
		return
	}

	return db.file.Sync()
}

func (db *SimpleDB) Update(key string, value []byte) (err error) {
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
		Operation: opPut,
		Key:       key,
		Value:     cloneBytes(value),
		CreatedAt: time.Now().UnixNano(),
	}

	_ = record
	return db.putRawLocked(key, value)
}

func (db *SimpleDB) ensureOpen() (err error) {
	return operationV2.NewTernary(operationV2.TrueValue(ErrDatabaseClosed)).GetByValue(db.closed)
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
			return nil
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
			db.appendVersionLocked(record.Key, cloneBytes(record.Value), false, record.CreatedAt)
		case opDelete:
			current := db.index[record.Key]
			current.Value = nil
			current.Deleted = true
			current.UpdatedAt = record.CreatedAt
			db.index[record.Key] = current
			db.appendVersionLocked(record.Key, nil, true, record.CreatedAt)
		default:
			return fmt.Errorf("%w：%d", ErrCorruptedRecord, record.Operation)
		}
	}
}

func writeRecord(writer io.Writer, record logRecord) error {
	var (
		err              error
		payload          []byte
		header, checksum [4]byte
	)

	if payload, err = json.Marshal(record); err != nil {
		return err
	}
	if len(payload) > maxRecordSize {
		return fmt.Errorf("simpleDB: record too large: %d", len(payload))
	}

	binary.BigEndian.PutUint32(header[:], uint32(len(payload)))
	if _, err = writer.Write(header[:]); err != nil {
		return err
	}
	if _, err = writer.Write(payload); err != nil {
		return err
	}

	binary.BigEndian.PutUint32(checksum[:], crc32.ChecksumIEEE(payload))
	if _, err = writer.Write(checksum[:]); err != nil {
		return err
	}

	return nil
}

func (db *SimpleDB) putRawLocked(key string, value []byte) error {
	return db.putRawAtLocked(key, value, time.Now().UnixNano())
}

func (db *SimpleDB) putRawAtLocked(key string, value []byte, at int64) error {
	var (
		err    error
		record = logRecord{
			Operation: opPut,
			Key:       key,
			Value:     cloneBytes(value),
			CreatedAt: at,
		}
	)

	if err = db.appendRecord(record); err != nil {
		return err
	}

	db.index[key] = entry{Value: cloneBytes(value), UpdatedAt: record.CreatedAt}
	db.appendVersionLocked(key, cloneBytes(value), false, record.CreatedAt)

	return nil
}

func (db *SimpleDB) deleteRawLocked(key string) error {
	return db.deleteRawAtLocked(key, time.Now().UnixNano())
}

func (db *SimpleDB) deleteRawAtLocked(key string, at int64) error {
	var (
		err    error
		record = logRecord{
			Operation: opDelete,
			Key:       key,
			CreatedAt: at,
		}
		current entry
	)

	if err = db.appendRecord(record); err != nil {
		return err
	}

	current = db.index[key]
	current.Value = nil
	current.Deleted = true
	current.UpdatedAt = record.CreatedAt
	db.index[key] = current
	db.appendVersionLocked(key, nil, true, record.CreatedAt)

	return nil
}
