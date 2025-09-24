package authorization

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog/log"
)

type Roles struct {
	Modules map[string][]string `json:"modules"`
}

var roles Roles

func init() {
	file, err := os.Open("roles.json")
	if err != nil {
		log.Debug().Msgf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Read all content
	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Debug().Msgf("Failed to read file: %v", err)
	}

	// Unmarshal into struct
	if err := json.Unmarshal(bytes, &roles); err != nil {
		log.Debug().Msgf("Failed to unmarshal JSON: %v", err)
	}

	log.Info().Msgf("Loaded roles: %+v", roles)
}

func ModuleRoles(module string) ([]string, error) {
	roles, ok := roles.Modules[module]

	if !ok {
		return []string{}, fmt.Errorf("module %s not found", module)
	}
	return roles, nil
}
