package tarq

import (
	"archive/tar"
	"context"
	"io"

	fq "github.com/takanoriyanagitani/go-fsq"
	aq "github.com/takanoriyanagitani/go-fsq/pkg/archive"
)

type PushMany func(ctx context.Context, w *tar.Writer, items fq.Iter[fq.Item]) error

func (p PushMany) Build() aq.PushMany {
	return func(ctx context.Context, w io.Writer, items fq.Iter[fq.Item]) error {
		var tw *tar.Writer = tar.NewWriter(w)
		return fq.Err1st([]func() error{
			func() error { return p(ctx, tw, items) },
			func() error { return tw.Close() },
		})
	}
}

type item2tar func(t *tar.Writer) func(fq.Item) error

func (t item2tar) pushMany(ctx context.Context, tw *tar.Writer, items fq.Iter[fq.Item]) error {
	var wt func(fq.Item) error = t(tw)
	return items.TryForEach(wt)
}

type item struct {
	raw []byte
	hdr *tar.Header
}

func (i item) toWriter(w *tar.Writer) (int, error) {
	return fq.ComposeErr(
		func(tw *tar.Writer) (*tar.Writer, error) { return tw, tw.WriteHeader(i.hdr) },
		func(tw *tar.Writer) (int, error) { return tw.Write(i.raw) },
	)(w)
}

func writeItemBuilder(tw *tar.Writer, i fq.Item) func(*tar.Header) (int, error) {
	return fq.ComposeErr(
		func(hdr *tar.Header) (*tar.Writer, error) { return tw, tw.WriteHeader(hdr) },
		func(w *tar.Writer) (int, error) { return w.Write(i.Raw()) },
	)
}

func item2tarBuilderNew(gen headerGen) item2tar {
	var toItem func(fq.Item) (item, error) = gen.toItem
	return func(t *tar.Writer) func(fq.Item) error {
		var writeItem func(item) (int, error) = func(i item) (int, error) { return i.toWriter(t) }
		var i2t func(fq.Item) (int, error) = fq.ComposeErr(
			toItem,
			writeItem,
		)
		return fq.ErrOnly(i2t)
	}
}
