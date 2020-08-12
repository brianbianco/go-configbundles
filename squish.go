package configbundle

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func AddBundlePath(names []string, bundleDir string) []string {
	withPaths := make([]string, len(names), len(names))
	for i, v := range names {
		withPaths[i] = filepath.Join(bundleDir, v)
	}
	return withPaths
}

func CreateBundle(name string, bundleDir string) (string, error) {
	var visited []string

	b, _, err := RecurseBundles(name, bundleDir, visited)
	if err != nil {
		return "", err
	}
	b = append(b, name)
	b = AddBundlePath(b, bundleDir)
	fmt.Println("Bundles being merged:", b)

	tmpdir, err := ioutil.TempDir("", "")
	tmpdir = filepath.Join(tmpdir, name)
	if err != nil {
		return "", err
	}
	d, err := MergeBundles(b, tmpdir)
	if err != nil {
		return d, err
	}
	return d, nil
}

func MergeBundles(bundles []string, dir string) (string, error) {
	for _, d := range bundles {
		CopyBundle(d, dir)
	}
	return dir, nil
}

func CopyBundle(source string, dest string) error {
	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		basename := strings.TrimPrefix(path, source)
		if !(len(basename) < 1) {
			dest := filepath.Join(dest, basename)
			fmt.Println("copying", path, "to", dest)
			if d, err := IsDir(path); err != nil {
				return err
			} else if d {
				EnsureDir(dest)
			} else {
				EnsureDir(filepath.Dir(dest))
				CopyFile(path, dest)
			}
		}
		return nil
	})
	return nil
}

func RecurseBundles(bundle string, bundleDir string, visited []string) ([]string, []string, error) {
	var bundles []string

	if s, _ := SliceIncludes(visited, bundle); s {
		return bundles, visited, nil
	}
	visited = append(visited, bundle)

	ib, err := IncludedBundles(bundle, bundleDir)
	if err != nil {
		return bundles, visited, err
	}

	if len(ib) <= 0 {
		return bundles, visited, nil
	}

	bundles = make([]string, len(ib))
	copy(bundles, ib)
	for _, b := range ib {
		list, tmpvisited, _ := RecurseBundles(b, bundleDir, visited)
		bundles, _ = CompactSlices(bundles, list)
		visited, _ = CompactSlices(visited, tmpvisited)
	}

	return bundles, visited, nil
}

func IncludedBundles(bundle string, bundleDir string) ([]string, error) {
	f, err := os.Open(filepath.Join(bundleDir, bundle, "bundles.txt"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	defer f.Close()

	bundles := make([]string, 0, 20)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		bundles = append(bundles, scanner.Text())
	}
	return bundles, nil
}
