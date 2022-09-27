package fsq

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestName(t *testing.T) {
	t.Parallel()

	var ITEST_FSQ_DIRNAME string = os.Getenv("ITEST_FSQ_DIRNAME")
	if len(ITEST_FSQ_DIRNAME) < 1 {
		t.Skip("Skipping tests using filesystem...")
	}

	var root string = filepath.Join(ITEST_FSQ_DIRNAME, "name")

	t.Run("QueueDirStatCheckerNewBySize", func(t *testing.T) {
		t.Parallel()

		t.Run("Valid next queue", func(t *testing.T) {
			t.Parallel()

			var dirname string = filepath.Join(root, "QueueDirStatCheckerNewBySize/Valid-Next-Q")
			e := os.RemoveAll(dirname)
			t.Run("test dir dropped", checkErr(e))

			e = os.MkdirAll(dirname, 0755)
			t.Run("test dir created", checkErr(e))

			var op opener = openerNewOrDefault(fsDefault)
			var nc NextCheck = NextCheckBuilder(op)
			var nq NextQueue = NextQueueI64.
				WithoutDir().
				ToChecked(nc)
			var q1st string = filepath.Join(dirname, "0123456789abcdef")

			fb, e := QueueFilenameBuilderNew(q1st, nq)
			t.Run("builder got", check(nil == e, true))

			t.Run("zero checker", func(t *testing.T) {
				t.Parallel()

				var sc QueueDirStatChecker = QueueDirStatCheckerNewBySize(0)
				var dc QueueDirChecker = QueueDirCheckerNewStat(sc)
				var fg QueueFilenameGenerator = fb.
					ToGenerator().
					ToChecked(dc)

				next, e := fg(context.Background())
				t.Run("next queue name got", checkErr(e))
				t.Run("Must be same", check(next, filepath.Join(dirname, "0123456789abcdf0")))

				f, e := os.Create(filepath.Join(dirname, "empty.dat"))
				t.Run("empty file created", checkErr(e))
				_ = f.Close()

				_, e = fg(context.Background())
				t.Run("Must fail", check(nil != e, true))
			})
		})

		t.Run("Invalid next queue", func(t *testing.T) {
			t.Parallel()

			var dirname string = filepath.Join(root, "QueueDirStatCheckerNewBySize/Valid-Next-Q")
			e := os.RemoveAll(dirname)
			t.Run("test dir dropped", checkErr(e))

			e = os.MkdirAll(dirname, 0755)
			t.Run("test dir created", checkErr(e))

			var q1st string = filepath.Join(dirname, "0123456789abcdef")

			_, e = QueueFilenameBuilderNew(q1st, nil)
			t.Run("Must fail", check(nil != e, true))
		})

	})
}
