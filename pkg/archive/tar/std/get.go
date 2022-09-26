package tarq

import (
	"archive/tar"
	"context"
	"io"

	fq "github.com/takanoriyanagitani/go-fsq"
	aq "github.com/takanoriyanagitani/go-fsq/pkg/archive"
)

type GetMany func(ctx context.Context, r *tar.Reader) (items fq.Iter[fq.Item], err error)

type Reader2Bytes func(io.Reader) ([]byte, error)

var Reader2BytesUnlimited Reader2Bytes = io.ReadAll

func (g GetMany) ToGetMany() aq.GetMany {
	return func(ctx context.Context, r io.Reader) (items fq.Iter[fq.Item], err error) {
		return g(ctx, tar.NewReader(r))
	}
}

func GetManyBuilderNew(r2b Reader2Bytes) GetMany {
	var tar2bytes func(*tar.Reader) ([]byte, error) = fq.ComposeErr(
		func(r *tar.Reader) (io.Reader, error) {
			_, e := r.Next()
			return r, e
		},
		r2b,
	)

	var tar2item func(*tar.Reader) (fq.Item, error) = fq.ComposeErr(
		tar2bytes,
		fq.ErrFuncGen(fq.ItemNew),
	)
	return func(_ context.Context, r *tar.Reader) (items fq.Iter[fq.Item], err error) {
		return func() (i fq.Item, hasValue bool) {
			itm, e := tar2item(r)
			return itm, nil == e
		}, nil
	}
}

var GetManyUnlimited GetMany = GetManyBuilderNew(Reader2BytesUnlimited)
