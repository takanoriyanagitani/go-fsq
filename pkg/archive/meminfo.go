package aq

import (
	"bytes"
	"io/fs"
	"time"
)

var MemFilemodeDefault fs.FileMode = 0644
var TimeUnixEpoch time.Time = time.Unix(0, 0)

type MemInfo struct {
	name string
	size int64
	mode fs.FileMode
	date time.Time
}

type MemInfoBuilder struct {
	mode fs.FileMode
	tgen TimeProvider
}

var MemInfoBuilderDefault MemInfoBuilder = MemInfoBuilder{}.
	WithMode(MemFilemodeDefault).
	WithTimeProvider(TimeProviderUnixEpoch)

func (m MemInfoBuilder) WithMode(mode fs.FileMode) MemInfoBuilder {
	m.mode = mode
	return m
}

func (m MemInfoBuilder) WithTimeProvider(tgen TimeProvider) MemInfoBuilder {
	m.tgen = tgen
	return m
}

func (m MemInfoBuilder) Build(name string, size int64) MemInfo {
	return MemInfo{
		name: name,
		size: size,
		mode: m.mode,
		date: m.tgen(),
	}
}

func (m MemInfoBuilder) NewFile(name string, data []byte) MemFile {
	var mi MemInfo = m.Build(name, int64(len(data)))
	return MemFileNew(
		mi,
		bytes.NewReader(data),
	)
}

type TimeProvider func() time.Time

func TimeProviderConst(t time.Time) TimeProvider { return func() time.Time { return t } }

var TimeProviderUnixEpoch TimeProvider = TimeProviderConst(TimeUnixEpoch)

func (m MemInfo) Name() string       { return m.name }
func (m MemInfo) Size() int64        { return m.size }
func (m MemInfo) Mode() fs.FileMode  { return m.mode }
func (m MemInfo) ModTime() time.Time { return m.date }
func (m MemInfo) IsDir() bool        { return false }
func (m MemInfo) Sys() any           { return nil }
