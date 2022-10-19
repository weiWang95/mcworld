package blockv2

type DigType uint8

const (
	DigTypeNone DigType = iota
	DigTypeAll
)

type Diggable struct {
	DigType  DigType `json:"dig_type"`
	DigLevel uint8   `json:"dig_level"`
}

func (b *Diggable) Diggable() bool {
	return b.DigType != DigTypeNone
}
