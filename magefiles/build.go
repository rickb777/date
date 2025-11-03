// See https://magefile.org/

//go:build mage

// Build steps for the date API:
package main

import (
	"github.com/magefile/mage/sh"
	"log"
	"os"
)

var Default = Install

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
	return nil
}

// installs datetool
func Install() error {
	if err := Build(); err != nil {
		return err
	}
	if err := sh.RunV("go", "install", "./datetool"); err != nil {
		return err
	}
	return nil
}

func Coverage() error {
	if err := sh.RunV("go", "test", "-cover", "./...", "-coverprofile", "coverage.out", "-coverpkg", "./..."); err != nil {
		return err
	}
	if err := sh.RunV("go", "tool", "cover", "-func", "coverage.out", "-o", "report.out"); err != nil {
		return err
	}
	return nil
}

// tests the module on both amd64 and i386 architectures for Linux and Windows
func CrossCompile() error {
	win := "build"
	linux := "test"
	if os.Getenv("GOOS") == "windows" {
		win = "test"
		linux = "build"
	}
	log.Printf("Testing on Windows\n")
	if err := sh.RunWithV(map[string]string{"GOOS": "windows"}, "go", win, "./..."); err != nil {
		return err
	}
	for _, arch := range []string{"amd64", "386"} {
		log.Printf("Testing on Linux/%s\n", arch)
		env := map[string]string{"GOOS": "linux", "GOARCH": arch}
		if _, err := sh.Exec(env, os.Stdout, os.Stderr, "go", linux, "./..."); err != nil {
			return err
		}
	}
	return nil
}
