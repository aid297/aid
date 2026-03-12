package kernal

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"

	"github.com/aid297/aid/filesystem/filesystemV4"
	"github.com/aid297/aid/operation/operationV2"
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

type versionedEntry struct {
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
	versions   map[string][]versionedEntry
	lastMVCCAt int64
	schema     *TableSchema
	autoSeq    int64
	uniqueIdx  map[string]map[string]string
	indexIdx   map[string]map[string]map[string]struct{}
	closed     bool
	config     DatabaseConfig
}

func newSimpleDB(dbName, tableName string, attrs ...SchemaAttributer) (*SimpleDB, error) {
	if strings.HasPrefix(strings.TrimSpace(tableName), "_sys_") {
		dbName = systemDatabaseFor(dbName)
	}
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
		versions:   make(map[string][]versionedEntry),
		uniqueIdx:  make(map[string]map[string]string),
		indexIdx:   make(map[string]map[string]map[string]struct{}),
		config: DatabaseConfig{
			DefaultUUIDVersion:     DefaultUUIDVersion,
			DefaultCascadeMaxDepth: DefaultCascadeMaxDepth,
			DefaultUUIDWithHyphen:  boolPtr(DefaultUUIDWithHyphen),
			DefaultUUIDUppercase:   boolPtr(DefaultUUIDUppercase),
			MaxCPUCores:            detectSystemCPUCores(),
			MaxMemoryBytes:         detectSystemMemoryBytes(),
		},
	}

	db.setAttrs(attrs...)
	db.applyRuntimeResourceLimitsLocked()

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

	var (
		err      error
		firstErr error
	)

	if db.closed {
		return nil
	}
	db.closed = true

	if db.file != nil {
		if err = db.file.Close(); err != nil {
			captureFirstError(&firstErr, err)
		}
	}

	if db.lockFile != nil {
		if err = unlockFile(db.lockFile, db.lockMethod); err != nil {
			captureFirstError(&firstErr, err)
		}
		if err = db.lockFile.Close(); err != nil {
			captureFirstError(&firstErr, err)
		}
	}

	return firstErr
}

func (db *SimpleDB) setAttrs(attrs ...SchemaAttributer) {
	for idx := range attrs {
		attrs[idx].RegisterAttr(db)
	}
}

func (db *SimpleDB) SetAttrs(attrs ...SchemaAttributer) *SimpleDB {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.setAttrs(attrs...)

	db.applyRuntimeResourceLimitsLocked()

	return db
}

func (db *SimpleDB) GetConfig() DatabaseConfig {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var config DatabaseConfig = db.config

	if db.config.DefaultUUIDWithHyphen != nil {
		value := *db.config.DefaultUUIDWithHyphen
		config.DefaultUUIDWithHyphen = &value
	}
	if db.config.DefaultUUIDUppercase != nil {
		value := *db.config.DefaultUUIDUppercase
		config.DefaultUUIDUppercase = &value
	}

	return config
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

func (db *SimpleDB) getDefaultUUIDWithHyphen() bool {
	if db.config.DefaultUUIDWithHyphen != nil {
		return *db.config.DefaultUUIDWithHyphen
	}

	return DefaultUUIDWithHyphen
}

func (db *SimpleDB) getDefaultUUIDUppercase() bool {
	if db.config.DefaultUUIDUppercase != nil {
		return *db.config.DefaultUUIDUppercase
	}

	return DefaultUUIDUppercase
}

func (db *SimpleDB) getMaxCPUCores() int {
	if db.config.MaxCPUCores > 0 {
		return normalizeCPUCores(db.config.MaxCPUCores)
	}

	return detectSystemCPUCores()
}

func (db *SimpleDB) getMaxMemoryBytes() uint64 {
	if db.config.MaxMemoryBytes > 0 {
		return normalizeMemoryBytes(db.config.MaxMemoryBytes)
	}

	return detectSystemMemoryBytes()
}

func (db *SimpleDB) applyRuntimeResourceLimitsLocked() {
	var (
		cpuCores    int    = db.getMaxCPUCores()
		memoryBytes uint64 = db.getMaxMemoryBytes()
	)

	if cpuCores > 0 {
		runtime.GOMAXPROCS(cpuCores)
	}

	if memoryBytes == 0 {
		return
	}
	if memoryBytes > math.MaxInt64 {
		debug.SetMemoryLimit(math.MaxInt64)
		return
	}

	debug.SetMemoryLimit(int64(memoryBytes))
}

func normalizeCPUCores(requested int) int {
	var actual int = detectSystemCPUCores()

	if requested <= 0 {
		return actual
	}
	if requested > actual {
		return actual
	}

	return requested
}

func normalizeMemoryBytes(requested uint64) uint64 {
	var actual uint64 = detectSystemMemoryBytes()

	if actual == 0 {
		return requested
	}
	if requested == 0 || requested > actual {
		return actual
	}

	return requested
}

func detectSystemCPUCores() int {
	var cores int = runtime.NumCPU()

	if cores <= 0 {
		return 1
	}

	return cores
}

func detectSystemMemoryBytes() uint64 {
	switch runtime.GOOS {
	case "linux":
		return detectLinuxMemoryBytes()
	case "darwin":
		return detectDarwinMemoryBytes()
	default:
		return 0
	}
}

func detectLinuxMemoryBytes() uint64 {
	var (
		err error
		raw []byte
	)

	if raw, err = os.ReadFile("/proc/meminfo"); err != nil {
		return 0
	}

	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimSpace(line)

		if !strings.HasPrefix(line, "MemTotal:") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			return 0
		}

		value, parseErr := strconv.ParseUint(fields[1], 10, 64)
		if parseErr != nil {
			return 0
		}

		return value * 1024
	}

	return 0
}

func detectDarwinMemoryBytes() uint64 {
	var (
		output []byte
		err    error
		value  uint64
	)

	if output, err = exec.Command("sysctl", "-n", "hw.memsize").Output(); err != nil {
		return 0
	}

	if value, err = strconv.ParseUint(strings.TrimSpace(string(output)), 10, 64); err != nil {
		return 0
	}

	return value
}

func cloneBytes(value []byte) []byte {
	var cloned []byte

	if value == nil {
		return nil
	}

	cloned = make([]byte, len(value))
	copy(cloned, value)

	return cloned
}

func validateKey(key string) error {
	return operationV2.NewTernary(operationV2.TrueValue(ErrEmptyKey)).GetByValue(strings.TrimSpace(key) == "")
}

func captureFirstError(target *error, err error) {
	if target == nil || err == nil {
		return
	}

	if *target == nil {
		*target = err
	}
}
