package block

type DestructibleType uint8

const (
	DestructibleTypeAll DestructibleType = iota + 1
)

type IDestructible interface {
	GetDestructibleType() DestructibleType
	GetHardness() uint8
}

type BaseDestructible struct {
	dtype    DestructibleType
	hardness uint8
}

var _ IDestructible = (*BaseDestructible)(nil)

func NewBaseDestructible(dtype DestructibleType, hardness uint8) IDestructible {
	d := new(BaseDestructible)
	d.dtype = dtype
	d.hardness = hardness
	return d
}

func (d *BaseDestructible) GetDestructibleType() DestructibleType {
	return d.dtype
}

func (d *BaseDestructible) GetHardness() uint8 {
	return d.hardness
}
