package kernal

import (
	"time"
)

type txWrite struct {
	value   []byte
	deleted bool
}

type Tx struct {
	db       *SimpleDB
	snapshot int64
	readSet  map[string]int64
	writeSet map[string]txWrite
	closed   bool
}

func (db *SimpleDB) appendVersionLocked(key string, value []byte, deleted bool, updatedAt int64) {
	db.versions[key] = append(db.versions[key], versionedEntry{
		Value:     cloneBytes(value),
		Deleted:   deleted,
		UpdatedAt: updatedAt,
	})
	if updatedAt > db.lastMVCCAt {
		db.lastMVCCAt = updatedAt
	}
}

func (db *SimpleDB) latestVersionTimestampLocked(key string) int64 {
	versions := db.versions[key]
	if len(versions) == 0 {
		return 0
	}
	return versions[len(versions)-1].UpdatedAt
}

func (db *SimpleDB) readValueAtSnapshotLocked(key string, snapshot int64) ([]byte, bool, int64) {
	versions := db.versions[key]
	for idx := len(versions) - 1; idx >= 0; idx-- {
		current := versions[idx]
		if current.UpdatedAt > snapshot {
			continue
		}
		if current.Deleted {
			return nil, false, current.UpdatedAt
		}
		return cloneBytes(current.Value), true, current.UpdatedAt
	}
	return nil, false, 0
}

func (db *SimpleDB) BeginTx() (*Tx, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if err := db.ensureOpen(); err != nil {
		return nil, err
	}

	snapshot := db.lastMVCCAt
	if snapshot == 0 {
		snapshot = time.Now().UnixNano()
	}

	return &Tx{
		db:       db,
		snapshot: snapshot,
		readSet:  make(map[string]int64),
		writeSet: make(map[string]txWrite),
	}, nil
}

func (tx *Tx) ensureOpen() error {
	if tx == nil || tx.closed {
		return ErrTxClosed
	}
	return nil
}

func (tx *Tx) Get(key string) ([]byte, bool, error) {
	if err := tx.ensureOpen(); err != nil {
		return nil, false, err
	}
	if err := validateKey(key); err != nil {
		return nil, false, err
	}

	if pending, exists := tx.writeSet[key]; exists {
		if pending.deleted {
			tx.readSet[key] = tx.snapshot
			return nil, false, nil
		}
		tx.readSet[key] = tx.snapshot
		return cloneBytes(pending.value), true, nil
	}

	tx.db.mu.RLock()
	defer tx.db.mu.RUnlock()

	if err := tx.db.ensureOpen(); err != nil {
		return nil, false, err
	}

	value, exists, versionAt := tx.db.readValueAtSnapshotLocked(key, tx.snapshot)
	tx.readSet[key] = versionAt
	if !exists {
		return nil, false, nil
	}
	return value, true, nil
}

func (tx *Tx) Put(key string, value []byte) error {
	if err := tx.ensureOpen(); err != nil {
		return err
	}
	if err := validateKey(key); err != nil {
		return err
	}

	tx.writeSet[key] = txWrite{value: cloneBytes(value)}
	return nil
}

func (tx *Tx) Delete(key string) error {
	if err := tx.ensureOpen(); err != nil {
		return err
	}
	if err := validateKey(key); err != nil {
		return err
	}

	tx.writeSet[key] = txWrite{deleted: true}
	return nil
}

func (tx *Tx) Commit() error {
	if err := tx.ensureOpen(); err != nil {
		return err
	}

	tx.db.mu.Lock()
	defer tx.db.mu.Unlock()

	if err := tx.db.ensureOpen(); err != nil {
		return err
	}

	for key := range tx.readSet {
		if tx.db.latestVersionTimestampLocked(key) > tx.snapshot {
			return ErrTxConflict
		}
	}

	if len(tx.writeSet) == 0 {
		tx.closed = true
		return nil
	}

	commitAt := time.Now().UnixNano()
	if commitAt <= tx.db.lastMVCCAt {
		commitAt = tx.db.lastMVCCAt + 1
	}

	for key, write := range tx.writeSet {
		var err error
		if write.deleted {
			err = tx.db.deleteRawAtLocked(key, commitAt)
		} else {
			err = tx.db.putRawAtLocked(key, write.value, commitAt)
		}
		if err != nil {
			return err
		}
	}

	tx.closed = true
	return nil
}

func (tx *Tx) Rollback() error {
	if err := tx.ensureOpen(); err != nil {
		return err
	}
	tx.closed = true
	tx.writeSet = nil
	tx.readSet = nil
	return nil
}
