package aq

import (
	"bufio"
	"context"
	"io"
	"io/fs"

	fq "github.com/takanoriyanagitani/go-fsq"
)

type GetMany func(ctx context.Context, r io.Reader) (items fq.Iter[fq.Item], err error)

type GetManyBuilder func(NameChecker) fq.GetMany

// ToBuilder creates GetManyBuilder.
// Iterator will be copied(via ToArrayIter) before file close.
func (g GetMany) ToBuilder(f fs.FS) GetManyBuilder {
	return func(chk NameChecker) fq.GetMany {
		return func(ctx context.Context, filename string) (items fq.Iter[fq.Item], err error) {
			file, e := f.Open(filename)
			if nil != e {
				return nil, e
			}
			defer file.Close() // reading file -> ignore close error

			var br *bufio.Reader = bufio.NewReader(file)

			return fq.ComposeErr(
				func(r io.Reader) (fq.Iter[fq.Item], error) { return g(ctx, r) },
				func(i fq.Iter[fq.Item]) (fq.Iter[fq.Item], error) { return i.ToArrayIter(), nil },
			)(br)
		}
	}
}
