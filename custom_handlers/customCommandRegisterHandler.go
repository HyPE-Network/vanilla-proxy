package custom_handlers

import (
	"encoding/json"
	"log"
	"math"
	"strings"

	"github.com/HyPE-Network/vanilla-proxy/proxy"
	"github.com/HyPE-Network/vanilla-proxy/proxy/command"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type IEngineResponseCommand struct {
	BaseCommand       string                        `json:"baseCommand"`
	Name              string                        `json:"name"`
	Description       string                        `json:"description"`
	Aliases           []string                      `json:"aliases,omitempty"`
	Type              string                        `json:"type"`
	AllowedTypeValues []string                      `json:"allowedTypeValues,omitempty"`
	Children          []IEngineResponseCommandChild `json:"children"`
	CanBeCalled       bool                          `json:"canBeCalled"`
	RequiresOp        bool                          `json:"requiresOp"`
}

type IEngineResponseCommandChild struct {
	IEngineResponseCommand
	Parent string `json:"parent"`
	Depth  int    `json:"depth"`
}

type IMinecraftRawText struct {
	Text string `json:"text"`
}

type IMinecraftTextMessage struct {
	RawText []IMinecraftRawText `json:"rawtext"`
}

type commandEnum struct {
	Type    string
	Options []string
	Dynamic bool
}

// getParamTypeAndOptional maps custom command types to protocol.CommandParameter types
func valueToParamType(child IEngineResponseCommandChild) (t uint32, enum commandEnum) {
	switch child.Type {
	case "literal":
		return 0, commandEnum{
			Type:    "SubCommand" + child.Name,
			Options: []string{child.Name},
		}
	case "string":
		return protocol.CommandArgTypeString, enum
	case "int":
		return protocol.CommandArgTypeInt, enum
	case "float":
		return protocol.CommandArgTypeFloat, enum
	case "location":
		return protocol.CommandArgTypePosition, enum
	case "boolean":
		return 0, commandEnum{
			Type:    "bool",
			Options: []string{"true", "1", "false", "0"},
		}
	case "player":
		return protocol.CommandArgTypeTarget, enum
	case "target":
		return protocol.CommandArgTypeTarget, enum
	case "array":
		return 0, commandEnum{
			Type:    "array",
			Options: child.AllowedTypeValues,
		}
	case "duration":
		return protocol.CommandArgTypeString, enum
	case "playerName":
		return protocol.CommandArgTypeString, enum
	default:
		return protocol.CommandArgTypeString, enum
	}
}

// sendAvailableCommands sends all available commands of the server. Once sent, they will be visible in the
// /help list and will be auto-completed.
func formatAvailableCommands(commands map[string]IEngineResponseCommand) packet.AvailableCommands {
	pk := &packet.AvailableCommands{}
	var enums []commandEnum
	enumIndices := map[string]uint32{}

	var dynamicEnums []commandEnum
	dynamicEnumIndices := map[string]uint32{}

	for alias, c := range commands {
		if c.Name != alias {
			// Don't add duplicate entries for aliases.
			continue
		}

		params := c.Children
		overloads := make([]protocol.CommandOverload, len(params))

		aliasesIndex := uint32(math.MaxUint32)
		if len(c.Aliases) > 0 {
			aliasesIndex = uint32(len(enumIndices))
			enumIndices[c.Name+"Aliases"] = aliasesIndex
			enums = append(enums, commandEnum{Type: c.Name + "Aliases", Options: c.Aliases})
		}

		for i, param := range params {
			// if param.RequiresOp && !player.IsOP() {
			// 	continue
			// }
			t, enum := valueToParamType(param)
			t |= protocol.CommandArgValid

			opt := byte(0)
			if param.Type == "bool" {
				opt |= protocol.ParamOptionCollapseEnum
			}
			if len(enum.Options) > 0 || enum.Type != "" {
				if !enum.Dynamic {
					index, ok := enumIndices[enum.Type]
					if !ok {
						index = uint32(len(enums))
						enumIndices[enum.Type] = index
						enums = append(enums, enum)
					}
					t |= protocol.CommandArgEnum | index
				} else {
					index, ok := dynamicEnumIndices[enum.Type]
					if !ok {
						index = uint32(len(dynamicEnums))
						dynamicEnumIndices[enum.Type] = index
						dynamicEnums = append(dynamicEnums, enum)
					}
					t |= protocol.CommandArgSoftEnum | index
				}
			}
			overloads[i].Parameters = append(overloads[i].Parameters, protocol.CommandParameter{
				Name:     strings.ToLower(param.Name),
				Type:     t,
				Optional: false,
				Options:  opt,
			})
		}
		pk.Commands = append(pk.Commands, protocol.Command{
			Name:          strings.ToLower(c.Name),
			Description:   c.Description,
			AliasesOffset: aliasesIndex,
			Overloads:     overloads,
		})
	}
	pk.DynamicEnums = make([]protocol.DynamicEnum, 0, len(dynamicEnums))
	for _, e := range dynamicEnums {
		pk.DynamicEnums = append(pk.DynamicEnums, protocol.DynamicEnum{Type: e.Type, Values: e.Options})
	}

	enumValueIndices := make(map[string]uint32, len(enums)*3)
	pk.EnumValues = make([]string, 0, len(enumValueIndices))

	pk.Enums = make([]protocol.CommandEnum, 0, len(enums))
	for _, enum := range enums {
		protoEnum := protocol.CommandEnum{Type: enum.Type}
		for _, opt := range enum.Options {
			index, ok := enumValueIndices[opt]
			if !ok {
				index = uint32(len(pk.EnumValues))
				enumValueIndices[opt] = index
				pk.EnumValues = append(pk.EnumValues, opt)
			}
			protoEnum.ValueIndices = append(protoEnum.ValueIndices, uint(index))
		}
		pk.Enums = append(pk.Enums, protoEnum)
	}
	return *pk
}

type CustomCommandRegisterHandler struct {
}

func (CustomCommandRegisterHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.Text)
	if dataPacket.TextType != packet.TextTypeObject {
		return true, dataPacket, nil
	}

	// Parse the raw message
	var messageData IMinecraftTextMessage
	if err := json.Unmarshal([]byte(dataPacket.Message), &messageData); err != nil {
		log.Println("Failed to parse message JSON:", err)
		return false, dataPacket, err
	}

	// Extract the text field
	message := messageData.RawText[0].Text

	if !strings.HasPrefix(message, "[PROXY_SYSTEM][COMMANDS]=") {
		// not a update commands message
		return true, dataPacket, nil
	}

	// Server has sent commands for this player to register.
	commandsRaw := strings.TrimPrefix(message, "[PROXY_SYSTEM][COMMANDS]=")
	var commands (map[string]IEngineResponseCommand)

	err := json.Unmarshal([]byte(commandsRaw), &commands)
	if err != nil {
		log.Println("Failed to unmarshal commands:", err)
		return false, dataPacket, err
	}

	// Prepare the AvailableCommands packet.
	availableCommands := formatAvailableCommands(commands)

	// Merge the commands here, with the existing commands.
	bdsSentCommands := proxy.ProxyInstance.Worlds.BDSAvailableCommands
	availableCommands = command.MergeAvailableCommands(availableCommands, bdsSentCommands)

	// Send the AvailableCommands packet to the player.
	player.DataPacket(&availableCommands)

	return false, dataPacket, nil
}
