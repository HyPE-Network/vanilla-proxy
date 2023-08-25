package chest

import "github.com/sandertv/gophertunnel/minecraft/protocol"

type Chest struct {
	Content []protocol.ItemInstance
}

func NewChest() *Chest {
	return &Chest{
		Content: make([]protocol.ItemInstance, 27),
	}
}

func (chest *Chest) GetContent() []protocol.ItemInstance {
	return chest.Content
}

func (chest *Chest) SetItem(index int, item protocol.ItemInstance) {
	chest.Content[index] = item
}
