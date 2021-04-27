package copy

import (
	"fmt"
	"io"
	"os"
)

func CopyFile(from, to string) (err error) {
	src, err := os.Open(from)
	if err != nil  {
		return err
	}
	defer func() {
		if cerr := src.Close(); cerr != nil {
			if err == nil {
				err = cerr
			}
		}
	}()
	stats, err := src.Stat()
	if err != nil {
		return err
	}
	dst, err := os.Create(to)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := dst.Close(); cerr != nil {
			if err == nil {
				err = cerr
			}
		}
	}()
	written , err := io.Copy(dst, src)
	if err != nil {
		return err
	}
	if written != stats.Size() {
		return fmt.Errorf("copied size: %d, original size: %d", written, stats.Size())
	}
	return nil

}
