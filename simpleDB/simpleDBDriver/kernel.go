package simpleDBDriver

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/aid297/aid/filesystem/filesystemV4"
)

const (
	dbLogTitle = "[SIMPLE-DB]"
	// defaultDataFile = "data.db"
	tempDataFile    = ".compact"
	lockFileEx      = ".lock"
	defaultDBFileEx = ".tbl"
	maxRecordSize   = 32 * 1024 * 1024
)

type fileLockMethod uint8

const (
	fileLockMethodNone fileLockMethod = iota
	fileLockMethodFlock
	fileLockMethodFcntl
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
	mu         sync.RWMutex
	database   string
	table      string
	dir        string
	dataPath   string
	file       *os.File
	lockPath   string
	lockFile   *os.File
	lockMethod fileLockMethod
	index      map[string]entry
	schema     *TableSchema
	autoSeq    int64
	uniqueIdx  map[string]map[string]string
	indexIdx   map[string]map[string]map[string]struct{}
	closed     bool
	config     DatabaseConfig
}

func newSimpleDB(dbName, tableName string) (*SimpleDB, error) {
	var (
		err        error
		dbPath     string
		lockPath   string
		file       *os.File
		lockFile   *os.File
		lockMethod fileLockMethod
		db         *SimpleDB
		dir        = filesystemV4.NewDir(filesystemV4.Rel(dbName, tableName))
	)

	if dbName == "" || tableName == "" {
		return nil, ErrDBPathEmpty
	}

	if dir.Create(filesystemV4.Mode(0o755)).GetError() != nil {
		return nil, fmt.Errorf("%w: 无法创建数据库目录", ErrInitDB)
	}

	dbPath = filesystemV4.NewFile(filesystemV4.Abs(dir.GetFullPath(), tableName+defaultDBFileEx)).GetFullPath()
	lockPath = dbPath + lockFileEx
	if lockFile, err = os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o644); err != nil {
		return nil, err
	}

	if lockMethod, err = lockFileExclusive(lockFile); err != nil {
		_ = lockFile.Close()
		return nil, err
	}

	if file, err = os.OpenFile(dbPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o644); err != nil {
		_ = unlockFile(lockFile, lockMethod)
		_ = lockFile.Close()
		return nil, err
	}

	db = &SimpleDB{
		database:   dbName,
		table:      tableName,
		dir:        dir.GetFullPath(),
		dataPath:   dbPath,
		file:       file,
		lockPath:   lockPath,
		lockFile:   lockFile,
		lockMethod: lockMethod,
		index:      make(map[string]entry),
		uniqueIdx:  make(map[string]map[string]string),
		indexIdx:   make(map[string]map[string]map[string]struct{}),
		config: DatabaseConfig{
			DefaultUUIDVersion:     DefaultUUIDVersion,
			DefaultCascadeMaxDepth: DefaultCascadeMaxDepth,
		},
	}

	if err = db.load(); err != nil {
		_ = file.Close()
		_ = unlockFile(lockFile, lockMethod)
		_ = lockFile.Close()
		return nil, err
	}

	if err = db.rebuildStructuredState(); err != nil {
		_ = file.Close()
		_ = unlockFile(lockFile, lockMethod)
		_ = lockFile.Close()
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

	var firstErr error
	if db.file != nil {
		if err := db.file.Close(); err != nil {
			captureFirstError(&firstErr, err)
		}
	}
	if db.lockFile != nil {
		if err := unlockFile(db.lockFile, db.lockMethod); err != nil {
			captureFirstError(&firstErr, err)
		}
		if err := db.lockFile.Close(); err != nil {
			captureFirstError(&firstErr, err)
		}
	}

	return firstErr
}

func (db *SimpleDB) SetConfig(config DatabaseConfig) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if config.DefaultUUIDVersion >= 1 && config.DefaultUUIDVersion <= 8 {
		db.config.DefaultUUIDVersion = config.DefaultUUIDVersion
	}
	if config.DefaultCascadeMaxDepth > 0 && config.DefaultCascadeMaxDepth <= HardCascadeMaxDepthLimit {
		db.config.DefaultCascadeMaxDepth = config.DefaultCascadeMaxDepth
	}
}

func (db *SimpleDB) GetConfig() DatabaseConfig {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.config
}

func (db *SimpleDB) getDefaultUUIDVersion() int {
	if db.config.DefaultUUIDVersion >= 1 && db.config.DefaultUUIDVersion <= 8 {
		return db.config.DefaultUUIDVersion
	}
	return DefaultUUIDVersion
}

func (db *SimpleDB) getDefaultCascadeMaxDepth() int {
	if db.config.DefaultCascadeMaxDepth > 0 && db.config.DefaultCascadeMaxDepth <= HardCascadeMaxDepthLimit {
		return db.config.DefaultCascadeMaxDepth
	}
	return DefaultCascadeMaxDepth
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

func captureFirstError(target *error, err error) {
	if target == nil || err == nil {
		return
	}
	if *target == nil {
		*target = err
	}
}
