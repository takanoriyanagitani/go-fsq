package tarq

import (
	"archive/tar"

	fq "github.com/takanoriyanagitani/go-fsq"
)

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

type headerGenBuilder func(mode int64) headerGen
