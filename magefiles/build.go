// See https://magefile.org/

//go:build mage

// Build steps for the expect API:
package main

import (
	"github.com/magefile/mage/sh"
	"log"
	"os"
)

var Default = Build

func Build() error {
	if err := Tidy(); err != nil {
		return err
	}
	if err := Test(); err != nil {
		return err
	}
	if err := sh.RunV("gofmt", "-l", "-w", "-s", "."); err != nil {
		return err
	}
	if err := Install(); err != nil {
		return err
	}
	return nil
}

// runs go mod download & tidy
func Tidy() error {
	if err := sh.RunV("go", "mod", "download"); err != nil {
		return err
	}
	if err := sh.RunV("go", "mod", "tidy"); err != nil {
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
func Test() error {
	dirs := []string{".", "./clock", "./gregorian", "./timespan", "./view"}
	for _, pkg := range dirs {
		if err := sh.RunV("go", "test", "-v", "-covermode=count", "-coverprofile="+nameOf(pkg), pkg); err != nil {
			return err
		}
	}
	for _, pkg := range dirs {
		if err := sh.RunV("go", "tool", "cover", "-func="+nameOf(pkg)); err != nil {
			return err
		}
		if err := sh.Run("rm", nameOf(pkg)); err != nil {
			return err
		}
	}
	if err := sh.RunV("go", "vet", "./..."); err != nil {
		return err
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

func nameOf(pkg string) string {
	if pkg == "." {
		return "date.out"
	}
	return pkg + ".out"
}
