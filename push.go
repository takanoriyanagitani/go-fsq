package fsq

import (
	"context"
)

type PushMany func(ctx context.Context, filename string, items Iter[Item]) error

func (p PushMany) PushAuto(ctx context.Context, b QueueFilenameGenerator, items Iter[Item]) error {
	var f func(context.Context) (string, error) = ComposeErr(
		b, // context.Context -> string, error
		func(filename string) (string, error) { return "", p(ctx, filename, items) },
	)
	return ErrOnly(f)(ctx)
}
