//go:build darwin || linux || freebsd || netbsd || openbsd || dragonfly || solaris

package kernal

import (
	"errors"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

func lockFileExclusive(file *os.File) (fileLockMethod, error) {
	var (
		err  error
		lock unix.Flock_t
	)

	if file == nil {
		return fileLockMethodNone, ErrInitDB
	}

	if err = unix.Flock(int(file.Fd()), unix.LOCK_EX|unix.LOCK_NB); err == nil {
		return fileLockMethodFlock, nil
	} else if isLockBusy(err) {
		return fileLockMethodNone, ErrDatabaseLocked
	} else if !shouldFallbackToFcntl(err) {
		return fileLockMethodNone, err
	}

	lock = unix.Flock_t{Type: unix.F_WRLCK, Whence: int16(unix.SEEK_SET), Start: 0, Len: 0}

	if err = unix.FcntlFlock(file.Fd(), unix.F_SETLK, &lock); err != nil {
		if isLockBusy(err) {
			return fileLockMethodNone, ErrDatabaseLocked
		}
		return fileLockMethodNone, err
	}

	return fileLockMethodFcntl, nil
}

func unlockFile(file *os.File, method fileLockMethod) error {
	if file == nil {
		return nil
	}

	switch method {
	case fileLockMethodFlock:
		return unix.Flock(int(file.Fd()), unix.LOCK_UN)
	case fileLockMethodFcntl:
		lock := unix.Flock_t{Type: unix.F_UNLCK, Whence: int16(unix.SEEK_SET), Start: 0, Len: 0}
		return unix.FcntlFlock(file.Fd(), unix.F_SETLK, &lock)
	default:
		return nil
	}
}

func isLockBusy(err error) bool {
	return errors.Is(err, unix.EWOULDBLOCK) ||
		errors.Is(err, unix.EAGAIN) ||
		errors.Is(err, syscall.EWOULDBLOCK) ||
		errors.Is(err, syscall.EAGAIN)
}

func shouldFallbackToFcntl(err error) bool {
	return errors.Is(err, unix.ENOTSUP) ||
		errors.Is(err, unix.ENOSYS) ||
		errors.Is(err, syscall.ENOTSUP) ||
		errors.Is(err, syscall.ENOSYS)
}
