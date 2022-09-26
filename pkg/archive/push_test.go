package aq

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	fq "github.com/takanoriyanagitani/go-fsq"
)

func TestPush(t *testing.T) {
	t.Parallel()

	ITEST_FSQ_DIRNAME := os.Getenv("ITEST_FSQ_DIRNAME")

	if len(ITEST_FSQ_DIRNAME) < 1 {
		t.Skip("Skipping tests using filesystem...")
	}

	var root string = filepath.Join(ITEST_FSQ_DIRNAME, "pkg/archive/push.d")

	t.Run("PushmanyFactory", func(t *testing.T) {
		t.Parallel()

		t.Run("Default", func(t *testing.T) {
			t.Parallel()

			var pm PushMany = func(_ context.Context, w io.Writer, items fq.Iter[fq.Item]) error {
				return nil
			}
			fm, e := PushmanyFactory{}.
				Default().
				WithPushMany(pm).
				Build()
			t.Run("Must not fail(PushmanyFactory)", check(nil == e, true))

			var dirname string = filepath.Join(root, "NameCheckerNoCheck")

			e = os.MkdirAll(dirname, 0755)
			t.Run("dir created", check(nil == e, true))

			var filename string = filepath.Join(dirname, "empty.dat")

			e = fm(context.Background(), filename, fq.IterEmpty[fq.Item]())
			t.Run("Must not fail", check(nil == e, true))
		})

		t.Run("invalid", func(t *testing.T) {
			t.Parallel()

			_, e := PushmanyFactory{}.
				Default().
				Build()
			t.Run("Must fail", check(nil != e, true))
		})
	})
}
