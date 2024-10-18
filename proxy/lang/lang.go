package lang

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/HyPE-Network/vanilla-proxy/utils"
	"github.com/sandertv/gophertunnel/minecraft/resource"
)

// GetSupportedLanguages returns a list of supported languages from a resource pack.
func GetSupportedLanguages(pack *resource.Pack) ([]string, error) {
	supportedLanguageBytes, err := pack.ReadFile("texts/languages.json")
	if err != nil {
		return nil, fmt.Errorf("error while reading languages.json: %w", err)
	}
	var supportedLanguages []string
	supportedLanguagesJSON, err := utils.ParseCommentedJSON([]byte(supportedLanguageBytes))
	if err != nil {
		return nil, fmt.Errorf("error while parsing languages.json: %w", err)
	}
	err = json.Unmarshal(supportedLanguagesJSON, &supportedLanguages)
	if err != nil {
		return nil, fmt.Errorf("error while un-marshalling languages.json: %w", err)
	}
	return supportedLanguages, nil
}

// GetLangTranslationMap returns a map of translations for a specific language from a resource pack.
func GetLangTranslationMap(pack *resource.Pack, language string) (map[string]string, error) {
	languageBytes, err := pack.ReadFile("texts/" + language + ".lang")
	if err != nil {
		return nil, fmt.Errorf("error while reading language file: %w", err)
	}

	langMap := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(string(languageBytes)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			// Skip empty lines or comments
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			langMap[key] = value
		} else {
			// Line has more than one `=` sign, or none at all, throw an error
			return nil, fmt.Errorf("error while parsing line in .lang file: %s", line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error while scanning .lang file: %w", err)
	}
	return langMap, nil
}
