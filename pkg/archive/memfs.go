package aq

import (
	"fmt"
	"io/fs"

	fq "github.com/takanoriyanagitani/go-fsq"
)

type MemFs struct {
	m map[string]fs.File
}

func MemFsNew() MemFs {
	return MemFs{
		m: make(map[string]fs.File),
	}
}

func pathErrNew(Op string) func(Err error) func(Path string) *fs.PathError {
	return func(Err error) func(string) *fs.PathError {
		return func(Path string) *fs.PathError {
			return &fs.PathError{
				Op,
				Path,
				Err,
			}
		}
	}
}

var openErrNew func(Err error) func(Path string) *fs.PathError = pathErrNew("open")

var notExistErr func(Path string) *fs.PathError = openErrNew(fs.ErrNotExist)
var invalidErr func(Path string) *fs.PathError = openErrNew(fs.ErrInvalid)

func (m MemFs) open(checkedName string) (fs.File, error) {
	f, found := m.m[checkedName]
	return fq.ErrFromBool(
		found,
		func() fs.File { return f },
		func() error { return notExistErr(checkedName) },
	)
}

func (m MemFs) getValidPath(unchecked string) (checked string, err error) {
	var valid bool = fs.ValidPath(unchecked)
	return fq.ErrFromBool(
		valid,
		func() string { return unchecked },
		func() error { return fmt.Errorf("Invalid path: %s", unchecked) },
	)
}

func (m MemFs) Open(unchecked string) (fs.File, error) {
	return fq.ComposeErr(
		m.getValidPath,
		m.open,
	)(unchecked)
}

//func (m MemFs) Upsert(nocheck string, f fs.File) {
//	m.m[nocheck] = f
//}
