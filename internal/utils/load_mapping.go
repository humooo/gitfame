package utils

import (
	"encoding/json"
	"fmt"

	configs "github.com/humooo/gitfame/configs"
	"github.com/humooo/gitfame/internal"
)

func LoadMapping() ([]internal.MappingEntity, error) {
	var mapping []internal.MappingEntity
	if err := json.Unmarshal(configs.LanguageExtensionsJSON, &mapping); err != nil {
		return nil, fmt.Errorf("could not decode language mapping: %w", err)
	}

	return mapping, nil
}
