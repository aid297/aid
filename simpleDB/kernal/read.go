package kernal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"sort"
	"strings"

	"github.com/aid297/aid/simpleDB/plugin"
	json "github.com/json-iterator/go"
)

func (db *SimpleDB) Get(key string) ([]byte, bool, error) {
	var (
		err     error
		current entry
		exists  bool
	)
	if err = validateKey(key); err != nil {
		return nil, false, err
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if err = db.ensureOpen(); err != nil {
		return nil, false, err
	}

	if err = db.hydrateMemDiskLocked(); err != nil {
		return nil, false, err
	}

	if current, exists = db.index[key]; !exists || current.Deleted {
		return nil, false, nil
	}

	return cloneBytes(current.Value), true, nil
}

func (db *SimpleDB) Query(prefix string) (map[string][]byte, error) {
	var (
		err    error
		result map[string][]byte
	)

	db.mu.Lock()
	defer db.mu.Unlock()

	if err = db.ensureOpen(); err != nil {
		return nil, err
	}

	if err = db.hydrateMemDiskLocked(); err != nil {
		return nil, err
	}

	result = make(map[string][]byte)
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
	var (
		err  error
		keys []string
	)

	db.mu.Lock()
	defer db.mu.Unlock()

	if err = db.ensureOpen(); err != nil {
		return nil, err
	}

	if err = db.hydrateMemDiskLocked(); err != nil {
		return nil, err
	}

	keys = make([]string, 0, len(db.index))
	for key, current := range db.index {
		if current.Deleted {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys, nil
}

func readRecord(reader io.Reader, compressor plugin.Compressor, encryptor plugin.Encryptor) (logRecord, error) {
	var (
		err                    error
		header, checksum       [4]byte
		size, expected, actual uint32
		payload                []byte
		record                 logRecord
	)

	if _, err = io.ReadFull(reader, header[:]); err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return logRecord{}, io.EOF
		}
		return logRecord{}, err
	}

	if size = binary.BigEndian.Uint32(header[:]); size == 0 || size > maxRecordSize {
		return logRecord{}, fmt.Errorf("%w: invalid payload size %d", ErrCorruptedRecord, size)
	}

	payload = make([]byte, size)
	if _, err = io.ReadFull(reader, payload); err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return logRecord{}, io.EOF
		}
		return logRecord{}, err
	}

	if _, err = io.ReadFull(reader, checksum[:]); err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return logRecord{}, io.EOF
		}
		return logRecord{}, err
	}

	expected = binary.BigEndian.Uint32(checksum[:])
	actual = crc32.ChecksumIEEE(payload)
	if expected != actual {
		return logRecord{}, fmt.Errorf("%w: checksum mismatch", ErrCorruptedRecord)
	}

	if payload, err = encryptor.Decrypt(payload); err != nil {
		return logRecord{}, fmt.Errorf("%w: decrypt failed: %v", ErrCorruptedRecord, err)
	}
	if payload, err = compressor.Decompress(payload); err != nil {
		return logRecord{}, fmt.Errorf("%w: decompress failed: %v", ErrCorruptedRecord, err)
	}

	if err = json.NewDecoder(bytes.NewReader(payload)).Decode(&record); err != nil {
		return logRecord{}, fmt.Errorf("%w: %v", ErrCorruptedRecord, err)
	}
	if record.Key == "" {
		return logRecord{}, fmt.Errorf("%w: empty key", ErrCorruptedRecord)
	}

	return record, nil
}
