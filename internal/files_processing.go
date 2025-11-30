package internal

import (
	"path/filepath"
	"strings"
)

type MappingEntity struct {
	Name       string
	Type       string
	Extensions []string
}

type FilesParams struct {
	FilesList []string
	Cla       *CommandLineArgs
	Mapping   []MappingEntity
}

func NewFilesParams(mapping []MappingEntity, cla *CommandLineArgs) *FilesParams {
	return &FilesParams{Cla: cla, Mapping: mapping}
}

func (fp *FilesParams) GetAllFiles(commitPointer, gitDir string) {
	gitTree, err := GitLsTree(commitPointer, gitDir)
	ProcessError(err, "GetAllFiles")

	filesInfo := strings.Split(gitTree, "\n")
	for _, file := range filesInfo {
		if file != "" {
			if !HasExtension(file, *fp.Cla.Extensions) {
				continue
			}
			if !IsAcceptableLanguage(GetLanguage(file, fp.Mapping), *fp.Cla.Languages) {
				continue
			}
			if len(*fp.Cla.Exclude) > 0 && MatchesPatterns(file, *fp.Cla.Exclude) {
				continue
			}
			if len(*fp.Cla.Restricted) > 0 && !MatchesPatterns(file, *fp.Cla.Restricted) {
				continue
			}
			fp.FilesList = append(fp.FilesList, file)
		}
	}
}

func HasExtension(path string, excludedExtensions []string) bool {
	if len(excludedExtensions) == 0 {
		return true
	}
	ext := filepath.Ext(path)
	for _, e := range excludedExtensions {
		if strings.EqualFold(ext, e) {
			return true
		}
	}
	return false
}

func IsAcceptableLanguage(fileLanguage string, languages []string) bool {
	if len(languages) == 0 {
		return true
	}

	if fileLanguage == "" {
		// Unknown Languages
		return false
	}

	for _, language := range languages {
		if strings.EqualFold(language, fileLanguage) {
			return true
		}
	}
	return false
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

func MatchesPatterns(filename string, patterns []string) bool {
	for _, pattern := range patterns {
		match, _ := filepath.Match(pattern, filename)
		if match {
			return true
		}
	}
	return false
}
