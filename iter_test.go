package fsq

import (
	"testing"
)

func TestIter(t *testing.T) {
	t.Parallel()

	t.Run("IterInts", func(t *testing.T) {
		t.Parallel()

		t.Run("empty", func(t *testing.T) {
			t.Parallel()

			var i Iter[int] = IterInts(0, 0)
			var cnt int = i.Count()
			t.Run("Must be empty", check(cnt, 0))
		})

		t.Run("0,1,2, ..., 9", func(t *testing.T) {
			t.Parallel()

			var i Iter[int] = IterInts(0, 10)
			var tot int = i.Reduce(0, AddInt)
			t.Run("Must be same", check(tot, 45))

			var i2 Iter[int] = IterInts(0, 10)
			var cnt int = i2.Count()
			t.Run("Must be same(count)", check(cnt, 10))
		})
	})
}
