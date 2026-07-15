package main

import (
	"fmt"
	"os"
	"syscall"
)

// acquireLock takes an exclusive, non-blocking file lock so that only one
// duckdbm process can run against the same database at a time.
func acquireLock(path string) (*os.File, error) {
	lockFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open lock file %s: %v", path, err)
	}
	if err := syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		_ = lockFile.Close()
		return nil, fmt.Errorf("another duckdbm instance is already running (lock: %s)", path)
	}
	return lockFile, nil
}

func releaseLock(lockFile *os.File) {
	_ = syscall.Flock(int(lockFile.Fd()), syscall.LOCK_UN)
	_ = lockFile.Close()
}
