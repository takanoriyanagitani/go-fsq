package fsq

import (
	"context"
)

type GetMany func(ctx context.Context, filename string) (items Iter[Item], err error)
