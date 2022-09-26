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

func (p PushMany) ToPushMany(f aq.PushmanyFactory) (fq.PushMany, error) {
	return f.
		WithPushMany(p.Build()).
		Build()
}

func (p PushMany) ToPushManyDefault() (fq.PushMany, error) {
	return p.ToPushMany(aq.PushmanyFactory{}.Default())
}

func (p PushMany) NewBuilder(tb aq.TempnameBuilder) aq.PushmanyBuilder {
	return fq.Compose(
		func(pm PushMany) aq.PushMany { return pm.Build() },
		aq.PushmanyBuilderNew(tb),
	)(p)
}

type item2tar func(t *tar.Writer) func(fq.Item) error

func (t item2tar) pushMany(ctx context.Context, tw *tar.Writer, items fq.Iter[fq.Item]) error {
	var wt func(fq.Item) error = t(tw)
	return items.TryForEach(wt)
}

func (t item2tar) toPushMany() PushMany { return t.pushMany }

type item struct {
	raw []byte
	hdr *tar.Header
}

func (i item) toWriter(w *tar.Writer) (int, error) {
	return fq.ComposeErr(
		func(hdr *tar.Header) ([]byte, error) { return i.raw, w.WriteHeader(hdr) },
		w.Write,
	)(i.hdr)
}

func writeItemBuilder(tw *tar.Writer, i fq.Item) func(*tar.Header) (int, error) {
	return fq.ComposeErr(
		func(hdr *tar.Header) ([]byte, error) { return i.Raw(), tw.WriteHeader(hdr) },
		tw.Write,
	)
}

func item2tarBuilderNew(gen headerGen) item2tar {
	var toItem func(fq.Item) (item, error) = gen.toItem
	return func(t *tar.Writer) func(fq.Item) error {
		writeItem := func(i item) (int, error) { return i.toWriter(t) }
		var i2t func(fq.Item) (int, error) = fq.ComposeErr(
			toItem,
			writeItem,
		)
		return fq.ErrOnly(i2t)
	}
}

type PushmanyBuilder func(mode int64) PushMany

func PushmanyBuilderNew(g NameGen) PushmanyBuilder {
	return func(mode int64) PushMany {
		return g.toPushMany(mode)
	}
}

var PushmanyBuilderUuidV4 PushmanyBuilder = PushmanyBuilderNew(NameGenUuidV4)

var PushmanyUuidV4Default PushMany = PushmanyBuilderUuidV4(FilemodeDefault)
