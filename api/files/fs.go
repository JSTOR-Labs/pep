package files

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func CopyFile(src, dst string) error {
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	return err
}

func CopyRecursive(src, dst string) error {
	if s, err := os.Stat(src); err == nil && s.IsDir() {
		_ = os.Mkdir(fmt.Sprintf("%s/%s", dst, s.Name()), s.Mode())
		contents, err := ioutil.ReadDir(src)
		if err != nil {
			return err
		}

		for _, f := range contents {
			err = CopyRecursive(fmt.Sprintf("%s/%s", src, f.Name()), fmt.Sprintf("%s/%s", dst, s.Name()))
			if err != nil {
				return err
			}
		}
	} else if err == nil {
		err = CopyFile(src, fmt.Sprintf("%s/%s", dst, s.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}
