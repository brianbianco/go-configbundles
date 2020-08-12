package configbundle

import (
	"io"
	"os"
)

func SliceIncludes(c []string, value string) (bool, error) {
	for _, i := range c {
		if i == value {
			return true, nil
		}
	}
	return false, nil
}

func CompactSlices(l []string, r []string) ([]string, error) {
	uniques := make([]string, len(l), len(l)+len(r))
	copy(uniques, l)
	for _, v := range r {
		if b, _ := SliceIncludes(uniques, v); b {
			continue
		} else {
			uniques = append(uniques, v)
		}
	}
	return uniques, nil
}

func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	return
}

func IsDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	return fi.IsDir(), err
}

func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}
