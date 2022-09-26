package tarq

import (
	"archive/tar"

	fq "github.com/takanoriyanagitani/go-fsq"
)

type namedItem struct {
	name string
	item fq.Item
}

func (i namedItem) Name() string { return i.name }
func (i namedItem) Size() int64  { return int64(i.item.Size()) }

type NameGen func(fq.Item) (filename string, err error)

func (g NameGen) toNamedItem(i fq.Item) (namedItem, error) {
	return fq.ComposeErr(
		g,
		func(filename string) (namedItem, error) {
			return namedItem{
				name: filename,
				item: i,
			}, nil
		},
	)(i)
}

type namedItemGen func(fq.Item) (namedItem, error)

func (gen namedItemGen) newBuilder() headerGenBuilder {
	return func(mode int64) headerGen {
		return fq.ComposeErr(
			gen,
			func(ni namedItem) (*tar.Header, error) {
				return &tar.Header{
					Name: ni.Name(),
					Mode: mode,
					Size: ni.Size(),
				}, nil
			},
		)
	}
}
