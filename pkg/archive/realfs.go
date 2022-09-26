package aq

import (
	"io/fs"
	"os"

	fq "github.com/takanoriyanagitani/go-fsq"
)

type RealFs struct{}

func (r RealFs) Open(name string) (fs.File, error) {
	return fq.ComposeErr(
		os.Open, // string -> *os.File, error
		func(f *os.File) (fs.File, error) { return f, nil },
	)(name)
}
