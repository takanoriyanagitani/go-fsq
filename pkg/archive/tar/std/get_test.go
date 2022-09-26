package tarq

import (
	"archive/tar"
	"bytes"
	"context"
	"testing"

	fq "github.com/takanoriyanagitani/go-fsq"
	aq "github.com/takanoriyanagitani/go-fsq/pkg/archive"
)

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("GetManyUnlimited", func(t *testing.T) {
		t.Parallel()

		t.Run("empty", func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			var tw *tar.Writer = tar.NewWriter(&buf)
			defer tw.Close()

			var rdr *bytes.Reader = bytes.NewReader(buf.Bytes())
			var tm GetMany = GetManyUnlimited
			var am aq.GetMany = tm.ToGetMany()

			items, e := am(context.Background(), rdr)
			t.Run("empty items got", check(nil == e, true))

			var cnt int = items.Count()
			t.Run("Must be empty", check(0, cnt))
		})

		t.Run("single empty queue", func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			var tw *tar.Writer = tar.NewWriter(&buf)
			defer tw.Close()

			e := tw.WriteHeader(&tar.Header{
				Name: "empty.dat",
				Mode: 0644,
				Size: 0,
			})
			t.Run("header wrote", check(nil == e, true))

			_, e = tw.Write(nil)
			t.Run("empty bytes wrote", check(nil == e, true))

			var rdr *bytes.Reader = bytes.NewReader(buf.Bytes())
			var tm GetMany = GetManyUnlimited
			var am aq.GetMany = tm.ToGetMany()

			items, e := am(context.Background(), rdr)
			t.Run("items got", check(nil == e, true))

			itm, hasValue := items()
			t.Run("Must not be empty", check(hasValue, fq.OptHasValue))

			t.Run("Must be empty(content)", check(itm.Size(), 0))
		})

		t.Run("many non-empty queues", func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			var tw *tar.Writer = tar.NewWriter(&buf)
			defer tw.Close()

			files := []struct {
				name string
				body []byte
			}{
				{"file1.txt", []byte("hw")},
				{"file2.txt", []byte("hh")},
			}

			for _, file := range files {
				e := tw.WriteHeader(&tar.Header{
					Name: file.name,
					Mode: 0644,
					Size: int64(len(file.body)),
				})
				t.Run("header wrote", check(nil == e, true))

				_, e = tw.Write(file.body)
				t.Run("empty bytes wrote", check(nil == e, true))
			}

			var rdr *bytes.Reader = bytes.NewReader(buf.Bytes())
			var tm GetMany = GetManyUnlimited
			var am aq.GetMany = tm.ToGetMany()

			items, e := am(context.Background(), rdr)
			t.Run("items got", check(nil == e, true))

			chk := func(name string, expected []byte) func(*testing.T) {
				return func(t *testing.T) {
					item, hasValue := items()
					t.Run("Must not be empty", check(hasValue, true))
					t.Run("Must not be empty(size)", check(item.Size(), len(expected)))

					t.Run("Must be same", checkBytes(item.Raw(), expected))
				}
			}

			t.Run("file1", chk("file1.txt", []byte("hw")))
			t.Run("file2", chk("file2.txt", []byte("hh")))

			_, hasValue := items()
			t.Run("Must be empty", check(hasValue, fq.OptEmpty))
		})
	})
}
