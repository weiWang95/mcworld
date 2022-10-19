package blockv2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/g3n/engine/texture"
	"github.com/g3n/engine/util/logger"
)

type BlockManager struct {
	log     *logger.Logger
	baseDir string
	texDir  string

	texMap   map[string]*texture.Texture2D
	blockMap map[BlockId]BlockAttr
}

func NewBlockManager(log *logger.Logger, baseDir string) *BlockManager {
	m := new(BlockManager)
	m.log = log
	m.baseDir = baseDir
	m.texDir = fmt.Sprintf("%s/images/blocks", m.baseDir)

	m.texMap = make(map[string]*texture.Texture2D)
	m.blockMap = make(map[BlockId]BlockAttr)

	m.init()

	return m
}

func (m *BlockManager) init() {
	m.initBlocks()
}

func (m *BlockManager) NewBlock(id BlockId) *Block {
	attr, ok := m.blockMap[id]
	if !ok {
		attr = m.blockMap[0]
	}

	b := new(Block)
	b.BlockAttr = &attr

	mesh := NewCube()
	mesh.SetTextures(m.loadTextures(&attr))
	b.IDrawable = mesh

	return b
}

func (m *BlockManager) initBlocks() {
	for _, item := range m.loadBlockAttrs() {
		m.blockMap[item.Id] = item
	}
}

func (m *BlockManager) loadBlockAttrs() []BlockAttr {
	bytes, err := ioutil.ReadFile(fmt.Sprintf("%s/config/block.json", m.baseDir))
	if err != nil {
		m.log.Warn("missing blocks data, %v", err)
		return nil
	}

	var data []BlockAttr
	if err := json.Unmarshal(bytes, &data); err != nil {
		m.log.Warn("unmarshal blocks data fail, %v", err)
		return nil
	}

	m.log.Info("success, %v block loaded", len(data))

	return data
}

func (m *BlockManager) loadTextures(attr *BlockAttr) []texture.Texture2D {
	texs := make([]texture.Texture2D, 0, len(attr.Textures))
	for _, item := range attr.Textures {
		tex := m.loadTexture(item)
		if tex == nil {
			tex = m.defaultTexture()
		}
		texs = append(texs, *tex)
	}
	return texs
}

func (m *BlockManager) defaultTexture() *texture.Texture2D {
	return m.loadTexture("default.jpg")
}

func (m *BlockManager) loadTexture(name string) *texture.Texture2D {
	tex, ok := m.texMap[name]
	if ok {
		return tex
	}

	tex, err := texture.NewTexture2DFromImage(fmt.Sprintf("%s/%s", m.texDir, name))
	if err != nil {
		m.log.Warn("missing texture:%s", name)
		m.texMap[name] = nil
		return nil
	}

	m.texMap[name] = tex

	return tex
}
