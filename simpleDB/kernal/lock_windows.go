//go:build windows

package kernal

import "os"

func lockFileExclusive(file *os.File) (fileLockMethod, error) {
	return fileLockMethodNone, ErrInitDB
}

func unlockFile(file *os.File, method fileLockMethod) error {
	return nil
}
