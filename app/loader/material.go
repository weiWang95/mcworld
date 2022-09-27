package loader

import (
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

func LoadBlockMaterial(id uint64) []material.IMaterial {
	tex := LoadBlockTexture(id)
	mats := make([]material.IMaterial, 0, 6)

	for i := 0; i < 6; i++ {
		mat := material.NewStandard(math32.NewColor("white"))
		mat.AddTexture(tex)
		mats = append(mats, mat)
	}

	return mats
}
