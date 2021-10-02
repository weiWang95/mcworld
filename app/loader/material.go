package loader

import (
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

var matMap map[uint64]material.IMaterial

func init() {
	matMap = make(map[uint64]material.IMaterial)
}

func LoadBlockMaterial(id uint64) material.IMaterial {
	if mat, ok := matMap[id]; ok {
		return mat
	}

	tex := LoadBlockTexture(id)
	mat := material.NewStandard(&math32.Color{1, 1, 1})
	mat.AddTexture(tex)

	matMap[id] = mat

	return mat
}
