package tarq

import (
	"archive/tar"

	fq "github.com/takanoriyanagitani/go-fsq"
)

const FilemodeDefault = 0644

type headerGen func(fq.Item) (*tar.Header, error)

func (h headerGen) toItem(i fq.Item) (item, error) {
	return fq.ComposeErr(
		h,
		func(hdr *tar.Header) (item, error) {
			return item{
				raw: i.Raw(),
				hdr: hdr,
			}, nil
		},
	)(i)
}

func (h headerGen) toItem2tar() item2tar { return item2tarBuilderNew(h) }
func (h headerGen) toPushMany() PushMany { return h.toItem2tar().toPushMany() }

type headerGenBuilder func(mode int64) headerGen
