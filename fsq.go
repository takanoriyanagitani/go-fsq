package fsq

import (
	"context"
)

type PushMany func(ctx context.Context, filename string, items Iter[Item]) error
type GetMany func(ctx context.Context, filename string) (items Iter[Item], err error)
