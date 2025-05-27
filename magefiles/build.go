// See https://magefile.org/

//go:build mage

// Build steps for the date API:
package main

import (
	"github.com/magefile/mage/sh"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var Default = Coverage

func Build() error {
	if err := sh.RunV("go", "test", "./..."); err != nil {
		return err
	}
	if err := sh.RunV("gofmt", "-l", "-w", "-s", "."); err != nil {
		return err
	}
	if err := sh.RunV("go", "vet", "./..."); err != nil {
		return err
	}
	if err := Install(); err != nil {
		return err
	}
	return nil
}

// installs datetool
func Install() error {
	if err := sh.RunV("go", "install", "./datetool"); err != nil {
		return err
	}
	return nil
}

// tests all the code and prints coverage information
func Coverage() error {
	if err := Build(); err != nil {
		return err
	}
	for _, dir := range listOfFoldersContainingTests() {
		if err := sh.RunV("go", "test", "-covermode=count", "-coverprofile="+dir+"test.out", packageName(dir)); err != nil {
			return err
		}
		if err := sh.RunV("go", "tool", "cover", "-func="+dir+"test.out"); err != nil {
			return err
		}
	}
	return nil
}

// tests the module on both amd64 and i386 architectures
func CrossCompile() error {
	for _, arch := range []string{"amd64", "386"} {
		log.Printf("Testing on %s\n", arch)
		env := map[string]string{"GOARCH": arch}
		if _, err := sh.Exec(env, os.Stdout, os.Stderr, "go", "test", "./..."); err != nil {
			return err
		}
		log.Printf("%s is good.\n\n", arch)
	}
	return nil
}

func listOfFoldersContainingTests() []string {
	root, _ := os.Getwd()
	fileSystem := os.DirFS(root)
	set := map[string]struct{}{}

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if strings.HasSuffix(path, "_test.go") {
			dir, _ := filepath.Split(path)
			set[dir] = struct{}{}
		}
		return nil
	})

	list := make([]string, 0, len(set))
	for dir := range set {
		list = append(list, dir)
	}
	sort.Strings(list)
	return list
}

func packageName(dir string) string {
	dir, _ = strings.CutSuffix(dir, "/")
	if dir == "" {
		return dir
	}
	return "./" + dir
}
