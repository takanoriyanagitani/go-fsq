package tarq

import (
	"context"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	fq "github.com/takanoriyanagitani/go-fsq"
	aq "github.com/takanoriyanagitani/go-fsq/pkg/archive"
)

func BenchmarkPush(b *testing.B) {
	var ITEST_FSQ_DIRNAME string = os.Getenv("ITEST_FSQ_DIRNAME")
	if len(ITEST_FSQ_DIRNAME) < 1 {
		b.Skip("Skipping benchmark using filesystem...")
	}

	var root string = filepath.Join(ITEST_FSQ_DIRNAME, "pkg/archive/tar/std/push")

	b.Run("pushmany got", func(tp PushMany) func(*testing.B) {
		return func(b *testing.B) {
			var dirname string = filepath.Join(root, "pushmany-got")
			e := os.MkdirAll(dirname, 0755)
			mustNil(e)

			var ap aq.PushMany = tp.Build()

			p, e := aq.PushmanyFactory{}.
				Default().
				WithPushMany(ap).
				Build()

			mustNil(e)

			var nc fq.NextCheck = fq.NextCheckBuilder(nil)
			var nq fq.NextQueue = fq.NextQueueI64.
				WithoutDir().
				ToChecked(nc)

			var q1st string = filepath.Join(dirname, "0123456789abcdef")
			fb, e := fq.QueueFilenameBuilderNew(q1st, nq)
			mustNil(e)

			var sc fq.QueueDirStatChecker = fq.QueueDirStatCheckerNewBySize(1048576)
			var dc fq.QueueDirChecker = fq.QueueDirCheckerNewStat(sc)

			var fg fq.QueueFilenameGenerator = fb.
				ToGenerator().
				ToChecked(dc)

			type pushMany func(context.Context, fq.Iter[fq.Item]) error

			var push pushMany = func(ctx context.Context, items fq.Iter[fq.Item]) error {
				return p.PushAuto(ctx, fg, items)
			}

			refreshDir := func(dn string) error {
				e := os.RemoveAll(dn)
				mustNil(e)
				return os.MkdirAll(dn, 0755)
			}

			b.Run("push many got", func(pm pushMany) func(*testing.B) {
				return func(b *testing.B) {
					b.Run("empty", func(b *testing.B) {
						e := refreshDir(dirname)
						mustNil(e)

						b.ResetTimer()
						for i := 0; i < b.N; i++ {
							e := push(context.Background(), fq.IterEmpty[fq.Item]())
							mustNil(e)
						}
					})

					b.Run("single", func(b *testing.B) {
						e := refreshDir(dirname)
						mustNil(e)

						newBytes := func() []byte {
							var rb []byte = make([]byte, 8192)
							_, e = rand.Read(rb)
							mustNil(e)
							return rb
						}

						newItems := func() fq.Iter[fq.Item] {
							return fq.IterFromArr([]fq.Item{
								fq.ItemNew(newBytes()),
							})
						}

						b.ResetTimer()
						for i := 0; i < b.N; i++ {
							e := push(context.Background(), newItems())
							mustNil(e)
						}
					})

					newMulti := func(wait time.Duration, batch int) func(*testing.B) {
						return func(b *testing.B) {
							e := refreshDir(dirname)
							mustNil(e)

							newBytes := func() []byte {
								var rb []byte = make([]byte, 8192)
								_, e = rand.Read(rb)
								mustNil(e)
								return rb
							}

							newItems := func() fq.Iter[fq.Item] {
								var i fq.Iter[int] = fq.IterInts(0, batch)
								var items fq.Iter[fq.Item] = fq.IterMap(i, func(_ int) fq.Item {
									return fq.ItemNew(newBytes())
								})
								return items
							}

							b.ResetTimer()
							for i := 0; i < b.N; i++ {
								e := push(context.Background(), newItems())
								mustNil(e)
								time.Sleep(wait)
							}
						}
					}

					b.Run("multi 10", newMulti(100*time.Millisecond, 10))
					b.Run("multi 16", newMulti(100*time.Millisecond, 16))
					b.Run("multi 128", newMulti(100*time.Millisecond, 128))
				}
			}(push))

		}
	}(PushmanyUuidV4Default))
}
