package block

const (
	GoldBlock           = "minecraft:gold_block"
	DeepslateDiamondOre = "minecraft:deepslate_diamond_ore"
	DiamondOre          = "minecraft:diamond_ore"
	IronBlock           = "minecraft:iron_block"
	NetheriteBlock      = "minecraft:netherite_block"
	Chest               = "minecraft:chest"
)

func IsDiamondOre(rid int32) bool {
	return GetBlockName(rid) == DiamondOre || GetBlockName(rid) == DeepslateDiamondOre
}
