package fsq

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestList(t *testing.T) {
	t.Parallel()

	t.Run("ListQueueBuilderUnlimited", func(t *testing.T) {
		t.Parallel()

		var ITEST_FSQ_DIRNAME string = os.Getenv("ITEST_FSQ_DIRNAME")
		if len(ITEST_FSQ_DIRNAME) < 1 {
			t.Skip("skipping tests using filesystem...")
		}

		var root string = filepath.Join(ITEST_FSQ_DIRNAME, "List/ListQueueBuilderUnlimited")

		t.Run("empty", func(t *testing.T) {
			t.Parallel()

			var dirname string = filepath.Join(root, "empty")
			e := os.MkdirAll(dirname, 0755)
			t.Run("empty dir created", check(nil == e, true))

			var lq ListQueue = ListQueueBuilderUnlimited(dirname)

			names, e := lq(context.Background(), -1)
			t.Run("names got", check(nil == e, true))

			var cnt int = names.Count()
			t.Run("empty names", check(cnt, 0))
		})

		t.Run("single empty queue", func(t *testing.T) {
			t.Parallel()

			var dirname string = filepath.Join(root, "single-empty-queue")
			e := os.MkdirAll(dirname, 0755)
			t.Run("dir created", check(nil == e, true))

			var filename string = filepath.Join(dirname, "empty.tmp")

			f, e := os.Create(filename)
			t.Run("empty file created", check(nil == e, true))
			f.Close()

			var lq ListQueue = ListQueueBuilderUnlimited(dirname)

			names, e := lq(context.Background(), -1)
			t.Run("names got", check(nil == e, true))

			name, hasValue := names()
			t.Run("name got", check(hasValue, true))
			t.Run("Must be same", check(name, filename))
		})

		t.Run("many empty queue files", func(t *testing.T) {
			t.Parallel()

			var dirname string = filepath.Join(root, "many-empty-queue-files")
			e := os.MkdirAll(dirname, 0755)
			t.Run("dir created", check(nil == e, true))

			create := func(basename string) error {
				var full string = filepath.Join(dirname, basename)
				f, e := os.Create(full)
				if nil != e {
					return e
				}
				return f.Close()
			}

			e = Err1st([]func() error{
				func() error { return create("sample1.empty.tmp") },
				func() error { return create("sample2.empty.txt") },
				func() error { return create("sample3.empty.txt") },
				func() error { return create("sample4.empty.tmp") },
				func() error { return create("sample5.empty.tmp") },
			})
			t.Run("test files created", check(nil == e, true))

			var lf ListFilter = ListFilterBuilderExt(".tmp").
				Negate()
			var lq ListQueue = ListQueueBuilderUnlimited(dirname).
				ToFiltered(lf)
			names, e := lq(context.Background(), -1)
			t.Run("names got", checkErr(e))

			chk := func(basename string) func(*testing.T) {
				return func(t *testing.T) {
					name, hasValue := names()
					t.Run("name got", check(hasValue, true))

					var full string = filepath.Join(dirname, basename)
					t.Run("Must be same", check(name, full))
				}
			}

			t.Run("sample2", chk("sample2.empty.txt"))
			t.Run("sample3", chk("sample3.empty.txt"))

			var cnt int = names.Count()
			t.Run("no more names", check(cnt, 0))
		})
	})
}
