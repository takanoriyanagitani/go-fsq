package aq

import (
	"bytes"
	"io/fs"
)

type MemFile struct {
	info fs.FileInfo
	data *bytes.Reader
}

func (m MemFile) Stat() (fs.FileInfo, error) { return m.info, nil }
func (m MemFile) Read(b []byte) (int, error) { return m.data.Read(b) }
func (m MemFile) Close() error               { return nil } // Nothing to close
