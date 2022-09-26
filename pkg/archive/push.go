package aq

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	fq "github.com/takanoriyanagitani/go-fsq"
)

type PushMany func(ctx context.Context, w io.Writer, items fq.Iter[fq.Item]) error

type NameChecker func(unchecked string) (checked string)

var NameCheckerNoCheck NameChecker = fq.Identity[string]

type PushmanyBuilder func(NameChecker) fq.PushMany

type PushmanyFactory struct {
	PushMany
	TempnameBuilder
	NameChecker
}

func (f PushmanyFactory) Default() PushmanyFactory {
	return f.
		WithTempnameBuilder(TempnameBuilderSimple).
		WithNameChecker(NameCheckerNoCheck)
}

func (f PushmanyFactory) WithPushMany(p PushMany) PushmanyFactory {
	f.PushMany = p
	return f
}

func (f PushmanyFactory) WithTempnameBuilder(t TempnameBuilder) PushmanyFactory {
	f.TempnameBuilder = t
	return f
}

func (f PushmanyFactory) WithNameChecker(c NameChecker) PushmanyFactory {
	f.NameChecker = c
	return f
}

func (f PushmanyFactory) Build() (fq.PushMany, error) {
	var valid bool = fq.IterFromArr([]bool{
		nil != f.PushMany,
		nil != f.TempnameBuilder,
		nil != f.NameChecker,
	}).All(fq.Identity[bool])
	return fq.ErrFromBool(
		valid,
		func() fq.PushMany {
			return fq.Compose(
				func(p PushMany) PushmanyBuilder { return p.newBuilder(f.TempnameBuilder) },
				func(b PushmanyBuilder) fq.PushMany { return b(f.NameChecker) },
			)(f.PushMany)
		},
		func() error { return fmt.Errorf("Invalid builder") },
	)
}

func (p PushMany) newBuilder(tmp TempnameBuilder) PushmanyBuilder {
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

func PushmanyBuilderNew(b TempnameBuilder) func(PushMany) PushmanyBuilder {
	return func(p PushMany) PushmanyBuilder { return p.newBuilder(b) }
}

var PushmanyBuilderSimple func(PushMany) PushmanyBuilder = PushmanyBuilderNew(TempnameBuilderSimple)
