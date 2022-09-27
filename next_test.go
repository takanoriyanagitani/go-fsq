package fsq

import (
	"context"
	"testing"
)

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
	})
}
