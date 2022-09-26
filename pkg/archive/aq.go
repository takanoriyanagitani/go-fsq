package aq

import (
	"context"
)

type NextQueue func(ctx context.Context, previous string) (next string, err error)
