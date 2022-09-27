package fsq

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
)

type ListFilter func(fullpath string) (ok bool)

type ListQueue func(ctx context.Context, limit int) (filenames Iter[string], e error)

func (l ListQueue) ToFiltered(f ListFilter) ListQueue {
	return ComposeContext(
		l,
		func(_ context.Context, names Iter[string]) (Iter[string], error) {
			return names.Filter(f), nil
		},
	)
}

func dirent2name(dirent fs.DirEntry) string { return dirent.Name() }

func basename2fullBuilder(dirname string) func(basename string) (full string) {
	return func(basename string) (full string) {
		return filepath.Join(dirname, basename)
	}
}

func dirent2fullBuilder(dirname string) func(dirent fs.DirEntry) (full string) {
	return Compose(
		dirent2name,                   // fs.DirEntry -> string(basename)
		basename2fullBuilder(dirname), // string(basename) -> string(full)
	)
}

func ListQueueBuilderUnlimited(dirname string) ListQueue {
	var dirent2full func(fs.DirEntry) string = dirent2fullBuilder(dirname)
	return func(_c context.Context, _l int) (filenames Iter[string], e error) {
		return ComposeErr(
			os.ReadDir,
			func(items []fs.DirEntry) (Iter[string], error) {
				var di Iter[fs.DirEntry] = IterFromArr(items)
				var si Iter[string] = IterMap(di, dirent2full)
				return si, nil
			},
		)(dirname)
	}
}
