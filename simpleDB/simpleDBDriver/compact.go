package simpleDBDriver

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
)

func (db *SimpleDB) Compact() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := db.ensureOpen(); err != nil {
		return err
	}

	tempPath := filepath.Join(db.dir, tempDataFile)
	tempFile, err := os.OpenFile(tempPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(tempFile)
	keys := make([]string, 0, len(db.index))
	for key, current := range db.index {
		if current.Deleted {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		current := db.index[key]
		record := logRecord{
			Operation: opPut,
			Key:       key,
			Value:     cloneBytes(current.Value),
			CreatedAt: current.UpdatedAt,
		}
		if err = writeRecord(writer, record); err != nil {
			_ = tempFile.Close()
			_ = os.Remove(tempPath)
			return err
		}
	}

	if err = writer.Flush(); err != nil {
		_ = tempFile.Close()
		_ = os.Remove(tempPath)
		return err
	}
	if err = tempFile.Sync(); err != nil {
		_ = tempFile.Close()
		_ = os.Remove(tempPath)
		return err
	}
	if err = tempFile.Close(); err != nil {
		_ = os.Remove(tempPath)
		return err
	}

	if err = db.file.Close(); err != nil {
		_ = os.Remove(tempPath)
		return err
	}
	if err = os.Rename(tempPath, db.dataPath); err != nil {
		_ = os.Remove(tempPath)
		return err
	}

	db.file, err = os.OpenFile(db.dataPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}

	return nil
}
