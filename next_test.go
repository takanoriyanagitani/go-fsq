package fsq

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

type opener func(string) (fs.File, error)

func (o opener) Open(name string) (fs.File, error) { return o(name) }

func TestNext(t *testing.T) {
	t.Parallel()

	t.Run("NextQueueI64", func(t *testing.T) {
		t.Parallel()

		t.Run("empty", func(t *testing.T) {
			t.Parallel()

			var nq NextQueue = NextQueueI64

			_, e := nq(context.Background(), "")
			t.Run("Must fail(empty)", check(nil != e, true))

			var next string = nq.UnwrapOrElse(
				context.Background(),
				"",
				func() string { return "1" },
			)
			t.Run("Must be same", check(next, "1"))
		})

		t.Run("non empty", func(t *testing.T) {
			t.Parallel()

			var nq NextQueue = NextQueueI64

			next, e := nq(context.Background(), "000000000000ff41")
			t.Run("Must not fail", check(nil == e, true))
			t.Run("Must be same", check(next, "000000000000ff42"))
		})

		t.Run("NextCheckBuilder", func(t *testing.T) {
			t.Parallel()

			ITEST_FSQ_DIRNAME := os.Getenv("ITEST_FSQ_DIRNAME")
			if len(ITEST_FSQ_DIRNAME) < 1 {
				t.Skip("skipping test using filesystem...")
			}

			var root string = filepath.Join(ITEST_FSQ_DIRNAME, "next/NextQueueI64/NextCheckBuilder")
			e := os.MkdirAll(root, 0755)
			t.Run("test dir created", check(nil == e, true))

			var open opener = ComposeErr(
				os.Open, // string -> *os.File, error
				func(f *os.File) (fs.File, error) { return f, nil },
			)

			var nc NextCheck = NextCheckBuilder(open)
			var nq NextQueue = NextQueueI64.WithoutDir().ToChecked(nc)

			t.Run("qfile not exists", func(t *testing.T) {
				t.Parallel()

				var pfilename string = filepath.Join(root, "not/exists/0123456789abcdef")
				var qfilename string = filepath.Join(root, "not/exists/0123456789abcdf0")
				e := os.RemoveAll(qfilename)
				t.Run("qfile must be absent", check(nil == e, true))

				_, e = nq(context.Background(), pfilename)
				t.Run("must not fail", check(nil == e, true))
			})

			t.Run("qfile already exists", func(t *testing.T) {
				t.Parallel()

				var pfilename string = filepath.Join(root, "exists/0123456789abcdef")
				var qfilename string = filepath.Join(root, "exists/0123456789abcdf0")

				var dirname string = filepath.Dir(pfilename)
				e := os.MkdirAll(dirname, 0755)
				t.Run("parent dir created", check(nil == e, true))

				f, e := os.Create(qfilename)
				t.Run("qfile created", check(nil == e, true))
				f.Close()

				_, e := nq(context.Background(), pfilename)
				t.Run("Must fail", check(nil != e, true))
			})
		})
	})
}
