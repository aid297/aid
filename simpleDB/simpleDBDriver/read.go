package simpleDBDriver

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"sort"
	"strings"
)

func (db *SimpleDB) Get(key string) ([]byte, bool, error) {
	if err := validateKey(key); err != nil {
		return nil, false, err
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	if err := db.ensureOpen(); err != nil {
		return nil, false, err
	}

	current, exists := db.index[key]
	if !exists || current.Deleted {
		return nil, false, nil
	}

	return cloneBytes(current.Value), true, nil
}

func (db *SimpleDB) Query(prefix string) (map[string][]byte, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if err := db.ensureOpen(); err != nil {
		return nil, err
	}

	result := make(map[string][]byte)
	for key, current := range db.index {
		if current.Deleted {
			continue
		}
		if prefix != "" && !strings.HasPrefix(key, prefix) {
			continue
		}
		result[key] = cloneBytes(current.Value)
	}

	return result, nil
}

func (db *SimpleDB) Keys() ([]string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if err := db.ensureOpen(); err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(db.index))
	for key, current := range db.index {
		if current.Deleted {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys, nil
}

func readRecord(reader io.Reader) (logRecord, error) {
	var header [4]byte
	if _, err := io.ReadFull(reader, header[:]); err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return logRecord{}, io.EOF
		}
		return logRecord{}, err
	}

	size := binary.BigEndian.Uint32(header[:])
	if size == 0 || size > maxRecordSize {
		return logRecord{}, fmt.Errorf("%w: invalid payload size %d", ErrCorruptedRecord, size)
	}

	payload := make([]byte, size)
	if _, err := io.ReadFull(reader, payload); err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return logRecord{}, io.EOF
		}
		return logRecord{}, err
	}

	var checksum [4]byte
	if _, err := io.ReadFull(reader, checksum[:]); err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return logRecord{}, io.EOF
		}
		return logRecord{}, err
	}

	expected := binary.BigEndian.Uint32(checksum[:])
	actual := crc32.ChecksumIEEE(payload)
	if expected != actual {
		return logRecord{}, fmt.Errorf("%w: checksum mismatch", ErrCorruptedRecord)
	}

	var record logRecord
	if err := json.NewDecoder(bytes.NewReader(payload)).Decode(&record); err != nil {
		return logRecord{}, fmt.Errorf("%w: %v", ErrCorruptedRecord, err)
	}
	if record.Key == "" {
		return logRecord{}, fmt.Errorf("%w: empty key", ErrCorruptedRecord)
	}

	return record, nil
}
