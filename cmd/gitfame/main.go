package main

import (
	"github.com/humooo/urlshortener/internal"
	"github.com/humooo/urlshortener/internal/utils"
	"github.com/humooo/urlshortener/internal/utils/constants"
)

type FameHandler struct {
	fp      *internal.FilesParams
	cla     *internal.CommandLineArgs
	mapping []internal.MappingEntity
	stats   *internal.Stats
}

func NewFameHandler() *FameHandler {
	cla := internal.NewCommandLineArgs()
	err := cla.GetCommandLineArgs()
	internal.ProcessError(err, "NewFameHandler")

	mapping := utils.LoadMapping()
	fp := internal.NewFilesParams(mapping, cla)
	fp.GetAllFiles(*fp.Cla.CommitPointer, *fp.Cla.RepositoryPath)
	return &FameHandler{fp: fp, cla: cla, mapping: mapping}
}

func (handler *FameHandler) GitFame() {
	stats := internal.CountStatistics(handler.fp)
	handler.stats = &stats
}

func (handler *FameHandler) PrintResults() {
	handler.stats.SortResults(handler.cla.SortOrderKey)
	switch handler.cla.OutputFormat {
	case constants.Tabular:
		handler.stats.PrintTabular()
	case constants.CSV:
		handler.stats.PrintCSV()
	case constants.SimpleJSON:
		handler.stats.PrintJSON()
	case constants.JSONLines:
		handler.stats.PrintJSONLines()
	}
}

func main() {
	handler := NewFameHandler()
	handler.GitFame()
	handler.PrintResults()
}
