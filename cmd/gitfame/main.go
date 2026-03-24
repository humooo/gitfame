package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/humooo/gitfame/internal"
	"github.com/humooo/gitfame/internal/utils"
	"github.com/humooo/gitfame/internal/utils/constants"
	flag "github.com/spf13/pflag"
)

func run(args []string) error {
	cla := internal.NewCommandLineArgs()
	if err := cla.GetCommandLineArgs(args); err != nil {
		return err
	}

	mapping, err := utils.LoadMapping()
	if err != nil {
		return err
	}

	fp := internal.NewFilesParams(mapping, cla)
	if len(fp.UnknownLanguages) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, "warning: unknown languages are ignored: %s\n", strings.Join(fp.UnknownLanguages, ", "))
	}

	if err := fp.GetAllFiles(cla.CommitPointer, cla.RepositoryPath); err != nil {
		return err
	}

	stats, err := internal.CountStatistics(fp)
	if err != nil {
		return err
	}

	stats.SortResults(cla.SortOrderKey)
	switch cla.OutputFormat {
	case constants.Tabular:
		return stats.PrintTabular()
	case constants.CSV:
		return stats.PrintCSV()
	case constants.SimpleJSON:
		return stats.PrintJSON()
	case constants.JSONLines:
		return stats.PrintJSONLines()
	default:
		return fmt.Errorf("unsupported output format: %s", cla.OutputFormat)
	}
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}

		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
