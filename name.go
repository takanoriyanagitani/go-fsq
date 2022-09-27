package fsq

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type QueueFilenameGenerator func(context.Context) (next string, err error)

type QueueDirChecker func(ctx context.Context, dirname string) (string, error)

func (g QueueFilenameGenerator) ToChecked(checker QueueDirChecker) QueueFilenameGenerator {
	return func(ctx context.Context) (next string, err error) {
		return ComposeErr(
			g, // context.Context -> string, error
			func(filename string) (checked string, err error) {
				var dirname string = filepath.Dir(filename)
				return ComposeContext(
					checker, // context.Context, string -> string, error
					func(_c context.Context, _s string) (string, error) { return filename, nil },
				)(ctx, dirname)
			},
		)(ctx)
	}
}

type QueueDirStatChecker func(stat fs.FileInfo) (fs.FileInfo, error)

func QueueDirStatCheckerNewBySize(max int64) QueueDirStatChecker {
	return func(stat fs.FileInfo) (fs.FileInfo, error) {
		var ok bool = stat.Size() <= max
		return ErrFromBool(
			ok,
			func() fs.FileInfo { return stat },
			func() error { return fmt.Errorf("Too many queues: %v > %v", stat.Size(), max) },
		)
	}
}

func QueueDirCheckerNewStat(schk QueueDirStatChecker) QueueDirChecker {
	return func(_ context.Context, dirname string) (checked string, err error) {
		return ComposeErr(
			os.Stat, // string -> fs.FileInfo, error
			ComposeErr(
				schk, // fs.FileInfo -> fs.FileInfo, error
				IgnoreArg[fs.FileInfo](func() (string, error) { return dirname, nil }),
			),
		)(dirname)
	}
}

type QueueFilenameBuilder struct {
	prev string
	next NextQueue
}

func QueueFilenameBuilderNew(init string, next NextQueue) (QueueFilenameBuilder, error) {
	return ErrFromBool(
		nil != next,
		func() QueueFilenameBuilder {
			return QueueFilenameBuilder{
				prev: init,
				next: next,
			}
		},
		func() error { return fmt.Errorf("Invalid NextQueue") },
	)
}

func (b *QueueFilenameBuilder) Next(ctx context.Context) (next string, err error) {
	next, err = b.next(ctx, b.prev)
	_ = ErrTryForEach(next, err, func(nex string) error {
		b.prev = nex
		return nil
	})
	return
}

func (b *QueueFilenameBuilder) ToGenerator() QueueFilenameGenerator { return b.Next }
