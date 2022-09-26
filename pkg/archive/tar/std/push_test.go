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
				})
			}
		}(PushmanyUuidV4Default))
	})
}
