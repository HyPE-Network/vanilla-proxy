package custom_handlers

import (
	"log"

	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/HyPE-Network/vanilla-proxy/utils"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type VectorXZ struct {
	X float32 `json:"x"`
	Z float32 `json:"z"`
}

type Location struct {
	Dimension string   `json:"dimension"`
	Pos1      VectorXZ `json:"pos1"`
	Pos2      VectorXZ `json:"pos2"`
}

type IPlayerClaim struct {
	ClaimId    string   `json:"claimId"`
	PlayerXUID string   `json:"playerXUID"`
	Location   Location `json:"location"`
	Trusts     []string `json:"trusts"`
}

var RegisteredClaims (map[string]IPlayerClaim)

func FetchClaims() error {
	claims, err := utils.FetchDatabase[IPlayerClaim]("claims")
	if err != nil {
		return err
	}

	RegisteredClaims = claims
	log.Printf("Loaded %d claims from the database\n", len(RegisteredClaims))

	return nil
}

// Dimension is the ID of the dimension that the player spawns in. It is a value from 0-2, with 0 being
// the overworld, 1 being the nether and 2 being the end.
func ClaimDimensionToInt(dimension string) int32 {
	if dimension == "minecraft:overworld" {
		return 0
	} else if dimension == "minecraft:nether" {
		return 1
	} else if dimension == "minecraft:end" {
		return 2
	} else {
		return -1
	}
}

// PlayerInsideClaim checks if a player is inside a claim
func PlayerInsideClaim(playerData minecraft.GameData, claim IPlayerClaim) bool {
	dimensionInt := ClaimDimensionToInt(claim.Location.Dimension)
	if dimensionInt != playerData.Dimension {
		return false
	}
	Pos1X, Pos1Z := float32(claim.Location.Pos1.X), float32(claim.Location.Pos1.Z)
	Pos2X, Pos2Z := float32(claim.Location.Pos2.X), float32(claim.Location.Pos2.Z)

	if playerData.PlayerPosition.X() >= Pos1X && playerData.PlayerPosition.X() <= Pos2X {
		if playerData.PlayerPosition.Z() >= Pos1Z && playerData.PlayerPosition.Z() <= Pos2Z {
			return true
		}
	}

	return false
}

func getClaimAt(dimension int32, x, z int32) IPlayerClaim {
	for _, claim := range RegisteredClaims {
		if ClaimDimensionToInt(claim.Location.Dimension) == dimension {
			Pos1X, Pos1Z := int32(claim.Location.Pos1.X), int32(claim.Location.Pos1.Z)
			Pos2X, Pos2Z := int32(claim.Location.Pos2.X), int32(claim.Location.Pos2.Z)

			if x >= Pos1X && x <= Pos2X {
				if z >= Pos1Z && z <= Pos2Z {
					return claim
				}
			}
		}
	}

	return IPlayerClaim{}
}

func canPreformActionInClaim(player human.Human, claim IPlayerClaim, action string) bool {
	// if player.GetData().GameData.PlayerPermissions == 2 {
	// 	return true
	// }

	if action == "interactWithBlock" && claim.PlayerXUID == "*" {
		// Players can interact with blocks in admin claims
		return true
	}

	if action == "interactWithEntity" && claim.PlayerXUID == "*" {
		// Players can interact with entities in admin claims
		return true
	}

	playerXuid := player.GetSession().IdentityData.XUID

	if claim.PlayerXUID == playerXuid || utils.StringInSlice(playerXuid, claim.Trusts) {
		return true
	}

	return false
}

type ClaimPlayerAuthInputHandler struct {
}

func (ClaimPlayerAuthInputHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	if player.IsOP() {
		return true, pk, nil
	}
	dataPacket := pk.(*packet.PlayerAuthInput)
	playerData := player.GetData().GameData

	// Handles Player movement and breaking of blocks.

	if len(dataPacket.BlockActions) == 0 {
		return true, pk, nil
	}

	// Something is being placed or broken, check if it is in a claim

	for _, ba := range dataPacket.BlockActions {
		log.Println(`Player break block`)
		if ba.Action == protocol.PlayerActionCrackBreak { // continue break
			continue
		}

		claim := getClaimAt(player.GetData().GameData.Dimension, ba.BlockPos.X(), ba.BlockPos.Z())
		if claim.ClaimId == "" {
			continue
		}

		if canPreformActionInClaim(player, claim, "breakBlock") {
			// Player is allowed to do action here
			continue
		}

		// Player does not own the claim or is not trusted, cancel the action
		player.SendMessage("§cYou cannot perform actions in this claim!")
		player.PlaySound("note.bass", playerData.PlayerPosition, 1, 1)
		return false, pk, nil
	}

	return true, pk, nil
}

type ClaimInventoryTransactionHandler struct {
}

func (ClaimInventoryTransactionHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	if player.IsOP() {
		return true, pk, nil
	}

	dataPacket := pk.(*packet.InventoryTransaction)
	playerData := player.GetData().GameData

	switch td := dataPacket.TransactionData.(type) {
	case *protocol.UseItemTransactionData:
		// Can stop players from using items like GUI, Pokedex, Potions, etc
		// So use this to stop players from throwing projectiles
		// This cannot stop interactions with entities or blocks.

		if td.HeldItem.Stack.ItemType.NetworkID == 0 {
			// Using hand
			return true, pk, nil
		}

		claim := getClaimAt(player.GetData().GameData.Dimension, int32(td.Position.X()), int32(td.Position.Z()))
		if claim.ClaimId == "" {
			return true, pk, nil
		}
		if canPreformActionInClaim(player, claim, "useItem") {
			return true, pk, nil
		}

		item := player.GetItemEntry(td.HeldItem.Stack.ItemType.NetworkID)
		if item == nil {
			// Item not sent over, most likely a minecraft item
			return true, pk, nil
		}
		itemComponents := player.GetItemComponentEntry(item.Name)
		if itemComponents == nil {
			// Item does not have any components, most-likely a minecraft item
			return true, pk, nil
		}

		if components, ok := itemComponents.Data["components"].(map[string]any); ok {
			if _, ok := components["minecraft:throwable"]; ok {
				// Item is throwable, stop the player from using it
				player.SendMessage("§cYou cannot use throwable items in this claim!")
				player.PlaySound("note.bass", playerData.PlayerPosition, 1, 1)
				return false, pk, nil
			}
		}
	}

	return true, pk, nil
}
