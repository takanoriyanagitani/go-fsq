package aq

import (
	"bufio"
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

			t.Run("empty", func(t *testing.T) {
				var pm PushMany = func(_c context.Context, w io.Writer, _i fq.Iter[fq.Item]) error {
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

			t.Run("many non-empty queues", func(t *testing.T) {
				var pm PushMany = func(_c context.Context, w io.Writer, i fq.Iter[fq.Item]) error {
					var bw *bufio.Writer = bufio.NewWriter(w)
					return i.TryForEach(fq.ErrOnly(fq.ComposeErr(
						func(itm fq.Item) (int, error) { return bw.Write(itm.Raw()) },
						func(_ int) (int, error) { return bw.Write([]byte("\n")) },
					)))
				}
				fm, e := PushmanyFactory{}.
					Default().
					WithPushMany(pm).
					Build()
				t.Run("Must not fail(PushmanyFactory)", check(nil == e, true))

				var dirname string = filepath.Join(root, "NameCheckerNoCheck")

				e = os.MkdirAll(dirname, 0755)
				t.Run("dir created", check(nil == e, true))

				var filename string = filepath.Join(dirname, "lines.txt")

				e = fm(context.Background(), filename, fq.IterFromArr([]fq.Item{
					fq.ItemNew([]byte("hw")),
					fq.ItemNew([]byte("hh")),
				}))
				t.Run("Must not fail", check(nil == e, true))

				var gm GetMany = func(_ context.Context, r io.Reader) (fq.Iter[fq.Item], error) {
					var br *bufio.Scanner = bufio.NewScanner(r)
					var i fq.Iter[fq.Item] = func() (i fq.Item, hasValue bool) {
						hasValue = br.Scan()
						return fq.ItemNew(br.Bytes()), hasValue
					}
					return i.ToArrayIter(), nil
				}
				var gmb GetManyBuilder = gm.ToBuilder(RealFs{})
				var fg fq.GetMany = gmb(NameCheckerNoCheck)

				items, e := fg(context.Background(), filename)
				t.Run("Items got", check(nil == e, true))

				chk := func(expected []byte) func(*testing.T) {
					return func(t *testing.T) {
						itm, hasValue := items()
						t.Run("itm got", check(hasValue, true))
						t.Run("Must be same", checkBytes(itm.Raw(), expected))
					}
				}

				t.Run("item 1", chk([]byte("hw")))
				t.Run("item 2", chk([]byte("hh")))
			})

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
