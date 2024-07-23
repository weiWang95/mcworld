package app

import "github.com/weiWang95/mcworld/app/blockv2"

type InventoryItem struct {
	blockId blockv2.BlockId
	count   uint8
}

func NewInventoryItem(blockId blockv2.BlockId, count uint8) *InventoryItem {
	return &InventoryItem{
		blockId: blockId,
		count:   count,
	}
}

func NewInventoryItems(blockId blockv2.BlockId, count uint8) []InventoryItem {
	max := Instance().bm.GetMaxStack(blockId)
	items := make([]InventoryItem, 0)
	for {
		if count <= max {
			items = append(items, *NewInventoryItem(blockId, count))
			break
		}

		items = append(items, *NewInventoryItem(blockId, max))
		count -= max
	}

	return items
}

type PlayerInventory struct {
	itemMap  map[blockv2.BlockId][]*InventoryItem
	bag      [4][10]*InventoryItem
	quickbar [10]*InventoryItem
}

func NewPlayerInventory() *PlayerInventory {
	p := new(PlayerInventory)
	p.init()
	return p
}

func (p *PlayerInventory) init() {
	p.reindexItems()
}

func (p *PlayerInventory) AddItem(blockId blockv2.BlockId, count uint8) uint8 {
	max := Instance().bm.GetMaxStack(blockId)

	if items, ok := p.itemMap[blockId]; ok {
		for i, _ := range items {
			if items[i].count == max {
				continue
			}

			if items[i].count+count <= max {
				items[i].count += count
				count = 0
				break
			} else {
				count = count + items[i].count - max
				items[i].count = max
			}
		}
		return count
	}

	return p.newItem(blockId, count)
}

func (p *PlayerInventory) newItem(blockId blockv2.BlockId, count uint8) uint8 {
	items := NewInventoryItems(blockId, count)
	cur := 0

	for i, _ := range p.quickbar {
		if p.quickbar[i] != nil {
			continue
		}

		p.quickbar[i] = &items[cur]
		p.reindexItem(p.quickbar[i])
		if cur == len(items)-1 {
			return 0
		}

		cur += 1
	}

	for i, _ := range p.bag {
		for j, _ := range p.bag[i] {
			if p.bag[i][j] != nil {
				continue
			}

			p.bag[i][j] = &items[cur]
			p.reindexItem(p.bag[i][j])
			if cur == len(items)-1 {
				return 0
			}

			cur += 1
		}
	}

	var remand uint8
	for {
		remand += items[cur].count

		if cur == len(items)-1 {
			break
		}
		cur += 1
	}

	return remand
}

func (p *PlayerInventory) reindexItems() {
	p.itemMap = make(map[blockv2.BlockId][]*InventoryItem)
	for i, _ := range p.quickbar {
		p.reindexItem(p.quickbar[i])
	}

	for i, _ := range p.bag {
		for j, _ := range p.bag[i] {
			p.reindexItem(p.bag[i][j])
		}
	}
}

func (p *PlayerInventory) reindexItem(item *InventoryItem) {
	if item == nil {
		return
	}

	if _, ok := p.itemMap[item.blockId]; !ok {
		p.itemMap[item.blockId] = []*InventoryItem{item}
		return
	}

	p.itemMap[item.blockId] = append(p.itemMap[item.blockId], item)
}
