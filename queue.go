package fsq

type Item struct {
	raw []byte
}

func (i Item) Raw() []byte { return i.raw }
func (i Item) Size() int   { return len(i.raw) }
