package tarq

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"testing"

	fq "github.com/takanoriyanagitani/go-fsq"
	aq "github.com/takanoriyanagitani/go-fsq/pkg/archive"
)

func TestPush(t *testing.T) {
	t.Parallel()

	t.Run("PushmanyUuidV4Default", func(t *testing.T) {
		t.Parallel()

		t.Run("Pushmany got", func(pm PushMany) func(*testing.T) {
			return func(t *testing.T) {
				t.Parallel()

				t.Run("ToPushManyDefault", func(t *testing.T) {
					t.Parallel()

					_, e := pm.ToPushManyDefault()
					t.Run("Pushmany got", check(nil == e, true))
				})

				t.Run("Build", func(t *testing.T) {
					t.Parallel()

					var am aq.PushMany = pm.Build()

					t.Run("empty", func(t *testing.T) {
						t.Parallel()

						var buf bytes.Buffer

						e := am(context.Background(), &buf, fq.IterEmpty[fq.Item]())
						t.Run("Must not fail(empty)", check(nil == e, true))

						var rdr *bytes.Reader = bytes.NewReader(buf.Bytes())

						var tr *tar.Reader = tar.NewReader(rdr)

						_, e = tr.Next()
						t.Run("Must be eof", check(e == io.EOF, true))
					})

					t.Run("single empty queue", func(t *testing.T) {
						t.Parallel()

						var buf bytes.Buffer

						e := am(context.Background(), &buf, fq.IterFromArr([]fq.Item{
							fq.ItemNew(nil),
						}))
						t.Run("Must not fail(empty)", check(nil == e, true))

						var rdr *bytes.Reader = bytes.NewReader(buf.Bytes())

						var tr *tar.Reader = tar.NewReader(rdr)

						th, e := tr.Next()
						t.Run("Must not fail(Next)", check(nil == e, true))

						t.Run("Must be empty(size)", check(0, th.Size))

						ba, e := io.ReadAll(tr)
						t.Run("Must not fail(ReadAll)", check(nil == e, true))

						t.Run("Must be empty(bytes)", check(0, len(ba)))
					})

					t.Run("many non-empty queues", func(t *testing.T) {
						t.Parallel()

						var buf bytes.Buffer

						e := am(context.Background(), &buf, fq.IterFromArr([]fq.Item{
							fq.ItemNew([]byte("hw")),
							fq.ItemNew([]byte("hh")),
						}))
						t.Run("Must not fail(non empty)", check(nil == e, true))

						var rdr *bytes.Reader = bytes.NewReader(buf.Bytes())

						var tr *tar.Reader = tar.NewReader(rdr)

						chk := func(expected []byte) func(*testing.T) {
							return func(t *testing.T) {
								th, e := tr.Next()
								t.Run("Must not fail(Next)", check(nil == e, true))

								t.Run("size check", check(th.Size, int64(len(expected))))

								ba, e := io.ReadAll(tr)
								t.Run("Must not fail(ReadAll)", check(nil == e, true))

								t.Run("Must be same", checkBytes(ba, expected))
							}
						}

						t.Run("item 1", chk([]byte("hw")))
						t.Run("item 2", chk([]byte("hh")))

					})
				})
			}
		}(PushmanyUuidV4Default))
	})
}
