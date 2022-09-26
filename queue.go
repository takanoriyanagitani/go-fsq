package fsq

type Item struct {
	raw []byte
}

func ItemNew(raw []byte) Item {
	return Item{raw}
}

func (i Item) Raw() []byte { return i.raw }
func (i Item) Size() int   { return len(i.raw) }
