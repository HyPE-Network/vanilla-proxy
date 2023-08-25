package inventory

import "github.com/sandertv/gophertunnel/minecraft/protocol"

type Inventory interface {
	GetContent() []protocol.ItemInstance
	SetItem(int, protocol.ItemInstance)
}
