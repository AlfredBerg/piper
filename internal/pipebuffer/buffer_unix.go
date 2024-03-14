//go:build unix

package pipebuffer

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/unix"
)

func Set(fd uintptr, size int) error {
	_, _, errno := unix.Syscall(syscall.SYS_FCNTL, fd, syscall.F_SETPIPE_SZ, uintptr(size))
	if errno != 0 {
		return fmt.Errorf(errno.Error())
	}
	return nil
}
