package aq

import (
	"context"
	"io"

	fq "github.com/takanoriyanagitani/go-fsq"
)

type GetMany func(ctx context.Context, r io.Reader) (items fq.Iter[fq.Item], err error)

type NextQueue func(ctx context.Context, previous string) (next string, err error)
