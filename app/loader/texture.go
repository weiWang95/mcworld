package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/g3n/engine/texture"
)

var sourcePath string
var textMap map[string]*texture.Texture2D

const defaultTex = "/images/blocks/default.jpg"

func init() {
	sourcePath = checkDirData("data")
	textMap = make(map[string]*texture.Texture2D)
}

func LoadBlockTexture(id uint64) *texture.Texture2D {
	return LoadTexture(fmt.Sprintf("/images/blocks/%d.jpg", id))
}

func LoadTexture(path string) *texture.Texture2D {
	if tex, ok := textMap[path]; ok {
		return tex
	}

	tex, err := texture.NewTexture2DFromImage(sourcePath + path)
	if err != nil {
		fmt.Printf("Error:%s loading texture:%s \n", err, path)
		tex = LoadTexture(defaultTex)
		textMap[path] = tex
		return tex
	}

	textMap[path] = tex

	return tex
}

func checkDirData(dirDataName string) string {
	// Check first if data directory is in the current directory
	if _, err := os.Stat(dirDataName); err != nil {
		panic(err)
	}
	dirData, err := filepath.Abs(dirDataName)
	if err != nil {
		panic(err)
	}
	return dirData
}
