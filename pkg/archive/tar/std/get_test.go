package tarq

import (
	"archive/tar"
	"bytes"
	"context"
	"testing"

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
	})
}
