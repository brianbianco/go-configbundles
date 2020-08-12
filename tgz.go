package configbundle

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func TgzDir(source string, targetname string) error {
	fmt.Println("Creating", targetname)
	target, err := os.Create(targetname)
	if err != nil {
		return err
	}
	defer target.Close()

	gw := gzip.NewWriter(target)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	filepath.Walk(source,
		func(path string, fi os.FileInfo, err error) error {
			header, err := tar.FileInfoHeader(fi, fi.Name())
			if err != nil {
				return err
			}
			header.Name = path
			fmt.Println("Adding", header.Name)
			if err := tw.WriteHeader(header); err != nil {
				return err
			}

			if fi.IsDir() {
				return nil
			}
			f, err := os.Open(path)
			if err != nil {
				return err
			}

			defer f.Close()
			_, err = io.Copy(tw, f)
			if err != nil {
				return err
			}

			return nil
		})

	return nil
}

func GenerateTgzName(source string) (string, error) {
	split := strings.Split(source, string(os.PathSeparator))
	name := split[len(split)-1]
	target := filepath.Join("/tmp", name+".tgz")
	return target, nil
}
