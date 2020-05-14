package main

import (
	"os"
	"strings"
)

type multiString []string

func (s *multiString) String() string {
	return strings.Join(*s, ", ")
}

func (s *multiString) Size() int {
	return len(*s)
}

func (s *multiString) Set(value string) error {
	*s = append(*s, value)
	return nil
}

type allowedFormats []string

func (a *allowedFormats) String() string {
	return strings.Join(*a, ", ")
}

func (a *allowedFormats) Exists(item string) bool {
	for _, b := range *a {
		if b == item {
			return true
		}
	}
	return false
}

type reportFile struct {
	Path   string
	Format string
}

func (f *reportFile) Exists() bool {
	info, err := os.Stat(f.Path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
