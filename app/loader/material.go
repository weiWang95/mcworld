package loader

import (
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

func LoadBlockMaterial(id uint64) []material.IMaterial {
	mats := make([]material.IMaterial, 0, 6)

	for i := 0; i < 6; i++ {
		tex := LoadBlockTexture(id, i+1)
		if tex == nil {
			tex = LoadBlockTexture(id, 0)
		}

		mat := material.NewStandard(math32.NewColor("white"))
		mat.AddTexture(tex)
		mats = append(mats, mat)
	}

	return mats
}

func LoadBlockFaceMaterial(id uint64, idx int) material.IMaterial {
	tex := LoadBlockTexture(id, idx+1)
	if tex == nil {
		tex = LoadBlockTexture(id, 0)
	}

	mat := material.NewStandard(math32.NewColor("white"))
	mat.AddTexture(tex)

	return mat
}
