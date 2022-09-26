package tarq

import (
	"archive/tar"

	"github.com/google/uuid"

	fq "github.com/takanoriyanagitani/go-fsq"
)

type namedItem struct {
	name string
	item fq.Item
}

func (i namedItem) Name() string { return i.name }
func (i namedItem) Size() int64  { return int64(i.item.Size()) }

type NameGen func(fq.Item) (filename string, err error)

var NameGenUuidV4 NameGen = fq.ComposeErr(
	fq.IgnoreArg[fq.Item](uuid.NewRandom),
	func(u uuid.UUID) (string, error) { return u.String(), nil },
)

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

func (g NameGen) toNamedItemGen() namedItemGen { return g.toNamedItem }

func (g NameGen) toHeaderGenBuilder() headerGenBuilder {
	return g.toNamedItemGen().toHeaderGenBuilder()
}

func (g NameGen) toHeaderGen(mode int64) headerGen { return g.toHeaderGenBuilder()(mode) }

func (g NameGen) toPushMany(mode int64) PushMany { return g.toHeaderGen(mode).toPushMany() }

type namedItemGen func(fq.Item) (namedItem, error)

func (gen namedItemGen) toHeaderGenBuilder() headerGenBuilder {
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
