package block

type ILumable interface {
	GetBlockLum() uint8
}

type BaseLumable struct {
	lum uint8
}

func NewBaseLumable(lum uint8) ILumable {
	return &BaseLumable{lum: lum}
}

func (l *BaseLumable) GetBlockLum() uint8 {
	return l.lum
}
