package app

type Luminance uint8

const LUM_BLOCK uint8 = 0b00001111
const LUM_SUN uint8 = 0b11110000
const MAX_LUM uint8 = 15

func NewLuminance(sunLum, blockLum uint8) Luminance {
	return Luminance(sunLum<<4 | blockLum)
}

func (l Luminance) SunLum() uint8 {
	return uint8(l) >> 4
}

func (l Luminance) BlockLum() uint8 {
	return uint8(l) & LUM_BLOCK
}

func (l Luminance) SetSunLum(sunLum uint8) Luminance {
	return Luminance(uint8(l)&LUM_BLOCK | (sunLum << 4))
}

func (l Luminance) SetBlockLum(blockLum uint8) Luminance {
	return Luminance(uint8(l)&LUM_SUN | blockLum)
}

func (l Luminance) Lum() uint8 {
	sun := l.SunLum()
	block := l.BlockLum()

	if sun > block {
		return sun
	}

	return block
}

func (l Luminance) CurLum(night bool) uint8 {
	if night {
		return l.BlockLum()
	}

	return l.Lum()
}

func MaxLum(l1, l2 Luminance) Luminance {
	var max Luminance
	if l1.SunLum() > l2.SunLum() {
		max = max.SetSunLum(l1.SunLum())
	} else {
		max = max.SetSunLum(l2.SunLum())
	}

	if l1.BlockLum() > l2.BlockLum() {
		max = max.SetBlockLum(l1.BlockLum())
	} else {
		max = max.SetBlockLum(l2.BlockLum())
	}

	return max
}
