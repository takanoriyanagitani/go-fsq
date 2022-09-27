package fsq

import (
	"context"
	"fmt"
	"strconv"
)

type NextQueue func(ctx context.Context, previous string) (next string, err error)

func (n NextQueue) UnwrapOrElse(ctx context.Context, prev string, alt func() string) (next string) {
	return ErrUnwrapOrElse(
		func(p string) (string, error) { return n(ctx, p) },
		func(_ error) string { return alt() },
	)(prev)
}

func (n NextQueue) ToChecked(checker func(next string) error) NextQueue {
	return ComposeContext(
		func(ctx context.Context, prev string) (next string, err error) { return n(ctx, prev) },
		func(ctx context.Context, next string) (string, error) { return next, checker(next) },
	)
}

var NextQueueI64 NextQueue = ComposeContext(
	func(_ context.Context, prev string) (int64, error) { return strconv.ParseInt(prev, 16, 64) },
	func(_ context.Context, p int64) (string, error) { return fmt.Sprintf("%016x", p+1), nil },
)
