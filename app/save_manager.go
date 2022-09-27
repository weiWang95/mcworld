package app

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/vmihailenco/msgpack"
	"github.com/weiWang95/mcworld/app/block"
	"github.com/weiWang95/mcworld/lib/util"
)

type ISaveManager interface {
	SaveChunk(data *Chunk) error
	LoadChunk(pos ChunkPos) *ChunkData
}

type fileSaveManager struct {
	ch chan *ChunkData

	baseDir  string
	chunkDir string
}

func newFileSaveManager() ISaveManager {
	sm := new(fileSaveManager)
	sm.baseDir = util.AbsPath("userdata/save")
	sm.chunkDir = util.AbsPath(fmt.Sprintf("%s/world/w0", sm.baseDir))
	sm.ch = make(chan *ChunkData, 20)
	sm.Start()
	return sm
}

func ConvertChunk(c *Chunk) ChunkData {
	data := ChunkData{
		Pos:  *c.pos,
		Data: make(map[cPos]BlockData),
	}
	for y := 0; y < len(c.blocks); y++ {
		for x := 0; x < len(c.blocks[0]); x++ {
			for z := 0; z < len(c.blocks[0][0]); z++ {
				b := c.blocks[y][x][z]
				if b != nil {
					d := BlockData{
						Id:    b.GetId(),
						State: b.GetState(),
					}

					data.SetBlock(x, y, z, d)
				}
			}
		}
	}

	return data
}

func (sm *fileSaveManager) Start() {
	util.SafeGo(func() {
		for data := range sm.ch {
			if err := sm.saveChunk(data); err != nil {
				Instance().Log().Error("save chunk fail:", err.Error())
			}
		}
	})
}

func (sm *fileSaveManager) Stop() {
	close(sm.ch)
}

func (sm *fileSaveManager) SaveChunk(c *Chunk) error {
	data := ConvertChunk(c)
	sm.ch <- &data
	return nil
}

func (sm *fileSaveManager) saveChunk(data *ChunkData) error {
	file, err := os.OpenFile(fmt.Sprintf("%s/%d_%d.chunk", sm.chunkDir, data.Pos.X, data.Pos.Z), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	bs, err := msgpack.Marshal(data)
	if err != nil {
		return err
	}

	_, err = file.Write(bs)
	if err != nil {
		return err
	}

	return nil
}

func (sm *fileSaveManager) LoadChunk(pos ChunkPos) *ChunkData {
	chunkFileName := sm.chunkFileName(pos)
	if _, err := os.Stat(chunkFileName); err != nil {
		Instance().Log().Debug("%s not exist", chunkFileName)
		return nil
	}
	data, err := ioutil.ReadFile(chunkFileName)
	if err != nil {
		panic(err)
	}

	var chunk ChunkData
	if err := msgpack.Unmarshal(data, &chunk); err != nil {
		panic(err)
	}

	return &chunk
}

func (sm *fileSaveManager) chunkFileName(pos ChunkPos) string {
	return fmt.Sprintf("%s/%d_%d.chunk", sm.chunkDir, pos.X, pos.Z)
}

type cPos uint16

type ChunkData struct {
	Pos  ChunkPos
	Data map[cPos]BlockData
}

func (cd *ChunkData) GetBlock(x, y, z int) *BlockData {
	b, ok := cd.Data[cd.posToKey(x, y, z)]
	if ok {
		return &b
	}

	return nil
}

func (cd *ChunkData) SetBlock(x, y, z int, b BlockData) {
	cd.Data[cd.posToKey(x, y, z)] = b
}

func (cd *ChunkData) posToKey(x, y, z int) cPos {
	return cPos(y)<<8 + cPos(x)<<4 + cPos(z)
}

type BlockData struct {
	Id    block.BlockId
	State block.BlockState
}
