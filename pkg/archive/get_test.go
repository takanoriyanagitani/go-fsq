package aq

import (
	"context"
	"io"
	"testing"

	fq "github.com/takanoriyanagitani/go-fsq"
)

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("ToBuilder", func(t *testing.T) {
		t.Parallel()

		t.Run("not exist", func(t *testing.T) {
			t.Parallel()

			var gm GetMany = func(_ context.Context, _ io.Reader) (fq.Iter[fq.Item], error) {
				return fq.IterEmpty[fq.Item](), nil
			}

			var mf MemFs = MemFsNew()

			var bldr GetManyBuilder = gm.ToBuilder(mf)

			var fm fq.GetMany = bldr(NameCheckerNoCheck)

			_, e := fm(context.Background(), "path/to/file/not-exist.dat")
			t.Run("Must fail", check(nil != e, true))
		})

		t.Run("empty file", func(t *testing.T) {
			t.Parallel()

			var gm GetMany = func(_ context.Context, r io.Reader) (fq.Iter[fq.Item], error) {
				data, e := io.ReadAll(r)
				if nil != e {
					return nil, e
				}
				items := fq.IterFromArr([]fq.Item{
					fq.ItemNew(data),
				})
				return items, nil
			}

			var mf MemFs = MemFsNew()

			var mib MemInfoBuilder = MemInfoBuilderDefault

			mf.Upsert(mib, "path/to/file/empty.dat", nil)

			var bldr GetManyBuilder = gm.ToBuilder(mf)

			var fm fq.GetMany = bldr(NameCheckerNoCheck)

			items, e := fm(context.Background(), "path/to/file/empty.dat")
			t.Run("Must not fail", check(nil == e, true))

			item, hasValue := items()
			t.Run("Must not be empty", check(hasValue, true))

			t.Run("Must be empty(data)", check(item.Size(), 0))
		})

		t.Run("invalid filename", func(t *testing.T) {
			t.Parallel()

			var gm GetMany = func(_ context.Context, _ io.Reader) (fq.Iter[fq.Item], error) {
				return fq.IterEmpty[fq.Item](), nil
			}

			var mf MemFs = MemFsNew()

			var mib MemInfoBuilder = MemInfoBuilderDefault

			mf.Upsert(mib, "path/to/file/empty.dat", nil)

			var bldr GetManyBuilder = gm.ToBuilder(mf)

			var fm fq.GetMany = bldr(NameCheckerNoCheck)

			_, e := fm(context.Background(), "")
			t.Run("Must fail", check(nil != e, true))
		})
	})
}
