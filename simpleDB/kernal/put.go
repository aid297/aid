package kernal

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"time"

	"github.com/aid297/aid/operation/operationV2"
	"github.com/aid297/aid/simpleDB/plugin"
	json "github.com/json-iterator/go"
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

func (db *SimpleDB) shouldPersistMemLocked() (bool, bool) {
	if db.schema == nil || db.schema.Engine != EngineMem || !db.schema.Disk {
		return false, false
	}

	// 阈值检查（内存上限）：清空内存模式
	threshold := db.config.Persistence.Threshold
	if threshold == 0 {
		threshold = db.schema.Persistence.Threshold
	}
	if threshold == 0 {
		threshold = 100 * 1024 * 1024 // 默认 100MB
	}
	if db.dirtyBytes >= threshold {
		return true, true
	}

	// 窗口期检查：落盘不代表清空内存
	windowSecs := db.config.Persistence.WindowSeconds
	if windowSecs == 0 {
		windowSecs = db.schema.Persistence.WindowSeconds
	}
	if windowSecs == 0 {
		windowSecs = 10
	}
	windowBytes := db.config.Persistence.WindowBytes
	if windowBytes == 0 {
		windowBytes = db.schema.Persistence.WindowBytes
	}
	if windowBytes == 0 {
		windowBytes = 10 * 1024 * 1024
	}

	if db.dirtyBytes >= windowBytes || time.Since(db.lastPersistAt).Seconds() >= float64(windowSecs) {
		return true, false
	}

	return false, false
}

func (db *SimpleDB) persistMemToDiskLocked(clearMemory bool) error {
	if db.isPersisting {
		return nil
	}
	db.isPersisting = true
	defer func() { db.isPersisting = false }()

	if len(db.memLog) == 0 {
		db.lastPersistAt = time.Now()
		return nil
	}

	// 确保文件已打开
	if db.file == nil {
		if err := db.ensureFileOpenLocked(); err != nil {
			return err
		}
	}

	// 批量写入记录
	for _, record := range db.memLog {
		if err := writeRecord(db.file, record, db.compressor, db.encryptor); err != nil {
			return err
		}
	}

	if err := db.file.Sync(); err != nil {
		return err
	}

	// 重置缓冲区
	db.memLog = make([]logRecord, 0)
	db.dirtyBytes = 0
	db.lastPersistAt = time.Now()

	if clearMemory {
		// 阈值模式：清空内存
		db.index = make(map[string]entry)
		db.versions = make(map[string][]versionedEntry)
		db.uniqueIdx = make(map[string]map[string]string)
		db.indexIdx = make(map[string]map[string]map[string]struct{})
		db.autoSeq = 0
		db.memCleared = true
	}

	return nil
}

func (db *SimpleDB) applyRecordToStateLocked(record logRecord) error {
	switch record.Operation {
	case opPut:
		db.index[record.Key] = entry{
			Value:     cloneBytes(record.Value),
			UpdatedAt: record.CreatedAt,
		}
		db.appendVersionLocked(record.Key, cloneBytes(record.Value), false, record.CreatedAt)
		return nil
	case opDelete:
		current := db.index[record.Key]
		current.Value = nil
		current.Deleted = true
		current.UpdatedAt = record.CreatedAt
		db.index[record.Key] = current
		db.appendVersionLocked(record.Key, nil, true, record.CreatedAt)
		return nil
	default:
		return fmt.Errorf("%w：%d", ErrCorruptedRecord, record.Operation)
	}
}

func (db *SimpleDB) hydrateMemDiskLocked() error {
	if db.schema == nil || db.schema.Engine != EngineMem || !db.schema.Disk || !db.memCleared {
		return nil
	}

	if db.file == nil {
		if err := db.ensureFileOpenLocked(); err != nil {
			return err
		}
	}

	memLogSnapshot := append([]logRecord(nil), db.memLog...)

	db.index = make(map[string]entry)
	db.versions = make(map[string][]versionedEntry)

	if err := db.load(); err != nil {
		return err
	}

	for _, record := range memLogSnapshot {
		if err := db.applyRecordToStateLocked(record); err != nil {
			return err
		}
	}

	db.memLog = memLogSnapshot

	if err := db.rebuildStructuredState(); err != nil {
		return err
	}

	db.memCleared = false
	return nil
}

func (db *SimpleDB) appendRecord(record logRecord) (err error) {
	if db.schema != nil && db.schema.Engine == EngineMem {
		if db.schema.Disk {
			// 内存引擎开启落盘：记录增量数据量并缓存日志
			payload, _ := json.Marshal(record)
			db.dirtyBytes += uint64(len(payload))
			db.memLog = append(db.memLog, record)

			// 检查是否触发落盘
			should, clear := db.shouldPersistMemLocked()
			if should {
				return db.persistMemToDiskLocked(clear)
			}
		}
		return nil
	}

	if db.file == nil {
		// 延迟创建文件（磁盘引擎但初始化时文件不存在）
		if err = db.ensureFileOpenLocked(); err != nil {
			return err
		}
	}

	if err = writeRecord(db.file, record, db.compressor, db.encryptor); err != nil {
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

	if bytes.Equal(current.Value, value) {
		return nil
	}
	return db.putRawLocked(key, value)
}

func (db *SimpleDB) ensureOpen() (err error) {
	return operationV2.NewTernary(operationV2.TrueValue(ErrDatabaseClosed)).GetByValue(db.closed)
}

func (db *SimpleDB) load() (err error) {
	if db.file == nil {
		return nil
	}

	var (
		reader *bufio.Reader
		record logRecord
	)

	if _, err = db.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	reader = bufio.NewReader(db.file)
	for {
		if record, err = readRecord(reader, db.compressor, db.encryptor); errors.Is(err, io.EOF) {
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

func writeRecord(writer io.Writer, record logRecord, compressor plugin.Compressor, encryptor plugin.Encryptor) error {
	var (
		err              error
		payload          []byte
		header, checksum [4]byte
	)

	if payload, err = json.Marshal(record); err != nil {
		return err
	}

	if payload, err = compressor.Compress(payload); err != nil {
		return err
	}
	if payload, err = encryptor.Encrypt(payload); err != nil {
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

	if err = db.hydrateMemDiskLocked(); err != nil {
		return err
	}

	// 先更新内存索引
	db.index[key] = entry{Value: cloneBytes(value), UpdatedAt: record.CreatedAt}
	db.appendVersionLocked(key, cloneBytes(value), false, record.CreatedAt)

	// 再追加记录（可能触发落盘和清空内存）
	if err = db.appendRecord(record); err != nil {
		return err
	}

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
	)

	if err = db.hydrateMemDiskLocked(); err != nil {
		return err
	}

	current := db.index[key]
	current.Value = nil
	current.Deleted = true
	current.UpdatedAt = record.CreatedAt
	db.index[key] = current
	db.appendVersionLocked(key, nil, true, record.CreatedAt)

	// 再追加记录（可能触发落盘和清空内存）
	if err = db.appendRecord(record); err != nil {
		return err
	}

	return nil
}
