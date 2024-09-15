package custom_handlers

import (
	"fmt"
	"strings"

	"github.com/HyPE-Network/vanilla-proxy/log"
	"github.com/HyPE-Network/vanilla-proxy/proxy"
	"github.com/HyPE-Network/vanilla-proxy/proxy/lang"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// cachedLangMap is a map of cached language translations.
var cachedLangMap map[string]map[string]string = make(map[string]map[string]string)

// getLanguageToTranslate returns the language to translate based on the players language code.
func getLanguageToTranslate(player human.Human) (string, error) {
	pokebedrockResourcePack := proxy.ProxyInstance.ResourcePacks[0]
	var languageToTranslate string = "en_US"

	// Read supported languages
	supportedLanguages, err := lang.GetSupportedLanguages(pokebedrockResourcePack)
	if err != nil {
		return languageToTranslate, fmt.Errorf("error while getting supported languages: %w", err)
	}

	// Check if players language is supported, if not default to `en_US`
	playersDesiredLanguage := player.GetSession().ClientData.LanguageCode
	for _, language := range supportedLanguages {
		if language == playersDesiredLanguage {
			languageToTranslate = playersDesiredLanguage
			break
		}
	}

	return languageToTranslate, nil
}

// getEntityTranslatedName returns the translated name of a entity type id based on the players language code.
func getEntityTranslatedName(player human.Human, entityTypeId string) (string, error) {
	pokebedrockResourcePack := proxy.ProxyInstance.ResourcePacks[0]
	// Get the language to translate
	languageToTranslate, err := getLanguageToTranslate(player)
	if err != nil {
		return "", fmt.Errorf("error while getting language to translate: %w", err)
	}

	// Check if the translation is already cached, if so return it
	if langMap, ok := cachedLangMap[languageToTranslate]; ok {
		return langMap["item.spawn_egg.entity."+entityTypeId+".name"], nil
	}

	// Get Lang file
	langMap, err := lang.GetLangTranslationMap(pokebedrockResourcePack, languageToTranslate)
	if err != nil {
		return "", fmt.Errorf("error while getting lang translation map: %w", err)
	}

	// Get the pokemon name
	translatedName := langMap["item.spawn_egg.entity."+entityTypeId+".name"]
	if translatedName == "" {
		return entityTypeId, fmt.Errorf("could not find translation for entity type id: %s", entityTypeId)
	}

	// Cache the translation
	cachedLangMap[languageToTranslate] = langMap

	return translatedName, nil
}

// getTranslatedNameTagOfSentOutPokemon returns the translated name tag of a sent out pokemon.
func getTranslatedNameTagOfSentOutPokemon(player human.Human, entityTypeId string, currentName string) (string, error) {
	if strings.HasPrefix(currentName, "§l§n§r") {
		// Name is nickname, do not translate
		return currentName, nil
	}

	// Get the translated name of the pokemon
	translatedName, err := getEntityTranslatedName(player, entityTypeId)
	if err != nil {
		return currentName, fmt.Errorf("error while getting translated name: %w", err)
	}

	// Read the current name, and replace the name with the translated name
	// Example: "§lImpidimp §eLvl 25\nOwner not found§r"
	// Example: "§lCharmander §eLvl 100\nSmell of curry§r"
	// Example: "§lTaillow\n§eLvl 21§r"
	// Example: "§lIron Treads\n§eLvl 21§r"
	// Example: "§lMr. Mime §eLvl 100\nSmell of curry§r"
	// Structure: "§l<Original Name> §eLvl <\d>\n<Other Information like Owner Name>§r"
	// Structure: "§l<Original Name>\n§eLvl <\d>§r"

	// Sanitize and process the current name string
	currentName = strings.TrimSpace(currentName)
	lines := strings.Split(currentName, "\n")

	if len(lines) > 0 {
		// Extract the level part to correctly replace the name
		if strings.Contains(lines[0], "§eLvl") {
			parts := strings.Split(lines[0], "§eLvl")
			lines[0] = "§l" + translatedName + " §eLvl" + parts[1]
		} else {
			lines[0] = "§l" + translatedName
		}

		// Rejoin the lines into a single string
		currentName = strings.Join(lines, "\n")
	} else {
		currentName = "§l" + translatedName
	}

	return currentName, nil
}

type AddActorNameTagHandler struct{}

func (h *AddActorNameTagHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.AddActor)

	// Check if the actor is a pokemon
	if !strings.HasPrefix(dataPacket.ActorType, "pokemon:") {
		return true, pk, nil
	}

	// Get the current name of the pokemon
	currentName, ok := dataPacket.ActorData[protocol.EntityDataKeyName].(string)
	if !ok {
		log.Logger.Warnln("Could not assert the current name of the pokemon:", dataPacket.ActorData[protocol.EntityDataKeyName])
		return true, pk, nil
	}

	// Get the translated name tag of the sent out pokemon
	translatedName, err := getTranslatedNameTagOfSentOutPokemon(player, dataPacket.ActorType, currentName)
	if err != nil {
		log.Logger.Warnln("Could not get the translated name tag of the sent out pokemon:", err)
		return true, pk, nil
	}

	// Update the name in the packet
	dataPacket.ActorData[protocol.EntityDataKeyName] = translatedName

	return true, pk, nil
}

type SetActorDataNameTagHandler struct{}

func (h *SetActorDataNameTagHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.SetActorData)

	// Get the current name of the pokemon
	currentName, ok := dataPacket.EntityMetadata[protocol.EntityDataKeyName].(string)
	if !ok {
		// NameTag is not being changed in this set data packet, ignore.
		return true, pk, nil
	}

	// Get the entity type id of the actor, to ensure its a pokemon
	entityID, ok := proxy.ProxyInstance.Entities.GetEntityFromRuntimeID(dataPacket.EntityRuntimeID)
	if !ok {
		// Entity was never being tracked.
		return true, pk, nil
	}

	actorTypeId, ok := proxy.ProxyInstance.Entities.GetEntityTypeID(entityID)
	if !ok {
		// AddActor was never called for this entity, so we don't know what type it is.
		return true, pk, nil
	}

	// Check if the actor is a pokemon
	if !strings.HasPrefix(actorTypeId, "pokemon:") {
		return true, pk, nil
	}

	// Get the translated name tag of the sent out pokemon
	translatedName, err := getTranslatedNameTagOfSentOutPokemon(player, actorTypeId, currentName)
	if err != nil {
		log.Logger.Warnln("Could not get the translated name tag of the sent out pokemon:", err)
		return true, pk, nil
	}

	// Update the name in the packet
	dataPacket.EntityMetadata[protocol.EntityDataKeyName] = translatedName

	return true, pk, nil
}
