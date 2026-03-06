package simpleDBDriver

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	dbLogTitle      = "[SIMPLE-DB]"
	defaultDataFile = "data.db"
	tempDataFile    = "data.db.compact"
	maxRecordSize   = 32 * 1024 * 1024
)

type operation uint8

const (
	opPut operation = iota + 1
	opDelete
)

type logRecord struct {
	Operation operation `json:"op"`
	Key       string    `json:"key"`
	Value     []byte    `json:"value,omitempty"`
	CreatedAt int64     `json:"createdAt"`
}

type entry struct {
	Value     []byte
	Deleted   bool
	UpdatedAt int64
}

type SimpleDB struct {
	mu       sync.RWMutex
	dir      string
	dataPath string
	file     *os.File
	index    map[string]entry
	closed   bool
}

func newSimpleDB(path string) (*SimpleDB, error) {
	var (
		err      error
		dataPath string
		file     *os.File
		db       *SimpleDB
	)

	if path == "" {
		return nil, ErrDBPathEmpty
	}

	if err = os.MkdirAll(path, 0o755); err != nil {
		return nil, err
	}

	dataPath = filepath.Join(path, defaultDataFile)
	if file, err = os.OpenFile(dataPath, os.O_CREATE|os.O_RDWR, 0o644); err != nil {
		return nil, err
	}

	db = &SimpleDB{
		dir:      path,
		dataPath: dataPath,
		file:     file,
		index:    make(map[string]entry),
	}

	if err = db.load(); err != nil {
		_ = file.Close()
		return nil, err
	}

	if _, err = db.file.Seek(0, io.SeekEnd); err != nil {
		_ = db.file.Close()
		return nil, err
	}

	return db, nil
}

func (db *SimpleDB) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.closed {
		return nil
	}
	db.closed = true
	return db.file.Close()
}

func cloneBytes(value []byte) []byte {
	if value == nil {
		return nil
	}
	cloned := make([]byte, len(value))
	copy(cloned, value)
	return cloned
}

func validateKey(key string) error {
	if strings.TrimSpace(key) == "" {
		return ErrEmptyKey
	}
	return nil
}
