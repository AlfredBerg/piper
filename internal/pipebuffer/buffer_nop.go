//go:build !unix

package pipebuffer

func Set(fd uintptr, size int) error {
	return nil
}
