package main

import (
	"flag"
	"fmt"
	cb "github.com/brianbianco/configbundle"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	bundleDir := os.Getenv("BUNDLE_DIR")

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		log.Fatal("No bundle name given as argument")
	}

	name := args[0]

	if len(bundleDir) <= 0 {
		bundleDir = "./bundles"
	}

	bundleDir = strings.TrimRight(bundleDir, "/")

	fmt.Printf("Bundle directory: %v\n", bundleDir)
	b, err := cb.CreateBundle(name, "./bundles")
	if err != nil {
		log.Fatal(err)
	}
	bparent := strings.TrimRight(b, name)
	od, _ := os.Getwd()
	os.Chdir(bparent)
	target := filepath.Join(os.TempDir(), name+".tgz")

	err = cb.TgzDir(name, target)
	if err != nil {
		log.Fatal(err)
	}
	os.Chdir(od)
	os.Exit(0)
}
