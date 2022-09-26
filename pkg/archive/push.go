package aq

import (
	"bufio"
	"context"
	"io"
	"os"

	fq "github.com/takanoriyanagitani/go-fsq"
)

type PushMany func(ctx context.Context, w io.Writer, items fq.Iter[fq.Item]) error

type NameChecker func(unchecked string) (checked string)

type PushmanyBuilder func(NameChecker) fq.PushMany

func (p PushMany) NewBuilder(tmp TempnameBuilder) PushmanyBuilder {
	return func(chk NameChecker) fq.PushMany {
		return func(ctx context.Context, filename string, items fq.Iter[fq.Item]) error {
			tmpname := tmp(filename)
			f, e := os.Create(chk(tmpname))
			if nil != e {
				return e
			}
			defer func() {
				_ = f.Close() // ignore close error after rename
			}()

			var bw *bufio.Writer = bufio.NewWriter(f)

			return fq.Err1st([]func() error{
				func() error { return p(ctx, bw, items) },
				func() error { return bw.Flush() },
				func() error { return f.Sync() },
				func() error { return os.Rename(tmpname, filename) },
			})
		}
	}
}

func (p PushMany) BuildSimple(chk NameChecker) fq.PushMany {
	var bldr PushmanyBuilder = p.NewBuilder(TempnameBuilderSimple)
	return bldr(chk)
}
