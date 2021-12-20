package main

import (
	"fmt"
	"os"
	"syscall"
)

func mount(source, target, fstype string, flags uintptr) error {
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	if err := syscall.Mount(source, target, fstype, flags, ""); err != nil {
		if sce, ok := err.(syscall.Errno); ok && sce == syscall.EBUSY {
			// already mounted
		} else {
			return fmt.Errorf("%v: %v", target, err)
		}
	}

	return nil
}
