package internal

import (
	"path/filepath"
	"strings"
)

type MappingEntity struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Extensions []string `json:"extensions"`
}

type FilesParams struct {
	FilesList        []string
	Cla              *CommandLineArgs
	Mapping          []MappingEntity
	LanguageFilter   map[string]struct{}
	UnknownLanguages []string
}

func NewFilesParams(mapping []MappingEntity, cla *CommandLineArgs) *FilesParams {
	languageFilter, unknownLanguages := buildLanguageFilter(mapping, cla.Languages)
	return &FilesParams{
		Cla:              cla,
		Mapping:          mapping,
		LanguageFilter:   languageFilter,
		UnknownLanguages: unknownLanguages,
	}
}

func buildLanguageFilter(mapping []MappingEntity, requestedLanguages []string) (map[string]struct{}, []string) {
	if len(requestedLanguages) == 0 {
		return nil, nil
	}

	knownLanguages := make(map[string]string, len(mapping))
	for _, mappingEntity := range mapping {
		if mappingEntity.Name == "" {
			continue
		}

		knownLanguages[strings.ToLower(mappingEntity.Name)] = mappingEntity.Name
	}

	filtered := make(map[string]struct{})
	unknown := make([]string, 0)
	for _, language := range requestedLanguages {
		normalized := strings.ToLower(language)
		if canonicalLanguage, ok := knownLanguages[normalized]; ok {
			filtered[strings.ToLower(canonicalLanguage)] = struct{}{}
		} else {
			unknown = append(unknown, language)
		}
	}

	if len(filtered) == 0 {
		return nil, unknown
	}

	return filtered, unknown
}

func (fp *FilesParams) GetAllFiles(commitPointer, gitDir string) error {
	gitTreeFiles, err := GitLsTree(commitPointer, gitDir)
	if err != nil {
		return err
	}

	for _, file := range gitTreeFiles {
		if file == "" {
			continue
		}

		if !HasExtension(file, fp.Cla.Extensions) {
			continue
		}

		fileLanguage := GetLanguage(file, fp.Mapping)
		if !IsAcceptableLanguage(fileLanguage, fp.LanguageFilter) {
			continue
		}

		if len(fp.Cla.Exclude) > 0 {
			matches, matchErr := MatchesPatterns(file, fp.Cla.Exclude)
			if matchErr != nil {
				return matchErr
			}

			if matches {
				continue
			}
		}

		if len(fp.Cla.Restricted) > 0 {
			matches, matchErr := MatchesPatterns(file, fp.Cla.Restricted)
			if matchErr != nil {
				return matchErr
			}

			if !matches {
				continue
			}
		}

		fp.FilesList = append(fp.FilesList, file)
	}

	return nil
}

func HasExtension(path string, allowedExtensions []string) bool {
	if len(allowedExtensions) == 0 {
		return true
	}

	ext := filepath.Ext(path)
	for _, allowedExtension := range allowedExtensions {
		if strings.EqualFold(ext, allowedExtension) {
			return true
		}
	}

	return false
}

func IsAcceptableLanguage(fileLanguage string, languageFilter map[string]struct{}) bool {
	if len(languageFilter) == 0 {
		return true
	}

	if fileLanguage == "" {
		return false
	}

	_, ok := languageFilter[strings.ToLower(fileLanguage)]
	return ok
}

func GetLanguage(path string, mapping []MappingEntity) string {
	fileExtension := filepath.Ext(path)
	for _, mappingEntity := range mapping {
		for _, extension := range mappingEntity.Extensions {
			if strings.EqualFold(fileExtension, extension) {
				return mappingEntity.Name
			}
		}
	}

	return ""
}

func MatchesPatterns(filename string, patterns []string) (bool, error) {
	for _, pattern := range patterns {
		match, err := filepath.Match(pattern, filename)
		if err != nil {
			return false, err
		}

		if match {
			return true, nil
		}
	}

	return false, nil
}
