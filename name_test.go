package fsq

import (
	"context"
	"fmt"
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

	getMinDirSize := func(dirname string) int64 {
		e := os.RemoveAll(dirname)
		mustNil(e)
		e = os.MkdirAll(dirname, 0755)
		mustNil(e)
		df, e := os.Open(dirname)
		mustNil(e)
		defer df.Close()
		stat, e := df.Stat()
		mustNil(e)
		var sz int64 = stat.Size()
		e = os.RemoveAll(dirname)
		mustNil(e)
		return sz
	}

	var minDirSize int64 = getMinDirSize(filepath.Join(root, "getMinDirSize.d"))

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

			t.Run("sz minimum checker", func(t *testing.T) {
				t.Parallel()

				// - file system 1(github)
				//   - empty dir size: 4096
				//   - dir size after empty file creation = 4096
				// - file system 2(local btrfs)
				//   - empty dir size: 0
				//   - dir size after empty file creation > 0
				var sc QueueDirStatChecker = QueueDirStatCheckerNewBySize(minDirSize)
				var dc QueueDirChecker = QueueDirCheckerNewStat(sc)
				var fg QueueFilenameGenerator = fb.
					ToGenerator().
					ToChecked(dc)

				next, e := fg(context.Background())
				t.Run("next queue name got", checkErr(e))
				t.Run("Must be same", check(next, filepath.Join(dirname, "0123456789abcdf0")))

				var ints Iter[int] = IterInts(0, 16)
				e = ints.TryForEach(func(i int) error {
					var name string = filepath.Join(dirname, fmt.Sprintf("%04x.dat", i))
					f, e := os.Create(name)
					mustNil(e)
					return f.Close()
				})
				t.Run("empty files created", checkErr(e))

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

		t.Run("queue already exists", func(t *testing.T) {
			t.Parallel()

			var dirname string = filepath.Join(root, "QueueDirStatCheckerNewBySize/q-exists")
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
			var qnex string = filepath.Join(dirname, "0123456789abcdf0")

			f, e := os.Create(qnex)
			t.Run("next queue file created", checkErr(e))
			_ = f.Close()

			fb, e := QueueFilenameBuilderNew(q1st, nq)
			t.Run("builder got", check(nil == e, true))

			t.Run("sz minimum checker", func(t *testing.T) {
				t.Parallel()

				var sc QueueDirStatChecker = QueueDirStatCheckerNewBySize(minDirSize)
				var dc QueueDirChecker = QueueDirCheckerNewStat(sc)
				var fg QueueFilenameGenerator = fb.
					ToGenerator().
					ToChecked(dc)

				_, e := fg(context.Background())
				t.Run("Must fail", check(nil != e, true))
			})
		})
	})
}
