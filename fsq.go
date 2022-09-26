package fsq

import (
	"context"
)

type PushMany func(ctx context.Context, filename string, items Iter[Item]) error
type GetMany func(ctx context.Context, filename string) (items Iter[Item], err error)
type DelMany func(ctx context.Context, filename string) error

type NextQueue func(ctx context.Context, previous string) (next string, err error)
