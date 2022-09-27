package fsq

import (
	"context"
	"os"
)

type DelMany func(ctx context.Context, filename string) error

var DelManyFs DelMany = func(_ context.Context, filename string) error { return os.Remove(filename) }
