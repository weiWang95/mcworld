package blockv2

type Lumable struct {
	Lum uint8 `json:"lum"`
}

func (b *Lumable) GetLumable() bool {
	return b.Lum > 0
}

func (b *Lumable) GetBlockLum() uint8 {
	return b.Lum
}
