package utils

import (
	"encoding/json"
	"os"

	"github.com/humooo/urlshortener/internal"
)

func LoadMapping() []internal.MappingEntity {
	curDir, err := os.Getwd()
	internal.ProcessError(err, "LoadMapping")

	rootPath := FindRoot(curDir)
	mappingData, err := os.Open(rootPath + "/gitfame/configs/language_extensions.json")
	defer func(mappingData *os.File) {
		err = mappingData.Close()
		internal.ProcessError(err, "LoadMapping")
	}(mappingData)
	internal.ProcessError(err, "LoadMapping")

	var mapping []internal.MappingEntity
	err = json.NewDecoder(mappingData).Decode(&mapping)
	internal.ProcessError(err, "LoadMapping")

	return mapping
}
