package main

import (
	"bridgecrewio/yor/common"
	"bridgecrewio/yor/common/logger"
	"bridgecrewio/yor/common/reports"
	"bridgecrewio/yor/common/runner"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
)

func main() {
	fmt.Println("Welcome to Yor!")
	options := &common.Options{}
	cmd := &cobra.Command{
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       common.Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(options)
		},
	}
	parseCommands(cmd.PersistentFlags(), options)
	options.Validate()
	cmd.SetVersionTemplate(fmt.Sprintf("Yor version %s", cmd.Version))
	if err := cmd.Execute(); err != nil {
		logger.Error(err.Error())
	}
}

func run(options *common.Options) error {
	yorRunner := new(runner.Runner)
	err := yorRunner.Init(options)
	if err != nil {
		logger.Error(err.Error())
	}
	reportService, err := yorRunner.TagDirectory(options.Directory)
	if err != nil {
		logger.Error(err.Error())
	}
	printReport(reportService, options)

	return nil
}

func printReport(reportService *reports.ReportService, options *common.Options) {
	reportService.CreateReport()
	switch strings.ToLower(options.Output) {
	case "cli":
		reportService.PrintToStdout()
	default:
		return
	}
}

func parseCommands(flag *pflag.FlagSet, commands *common.Options) {
	flag.StringVarP(&commands.Directory, "directory", "d", "", "directory to tag")
	flag.StringVarP(&commands.Tag, "tag", "t", "", "run yor only with the specified tag")
	flag.StringVarP(&commands.SkipTag, "skip-tag", "s", "", "run yor without ths specified tag")
	flag.StringVarP(&commands.Output, "output", "o", "cli", "set output format")
	flag.StringVar(&commands.OutputJsonFile, "output-json-file", "", "json file path for output")
	flag.StringSliceVarP(&commands.CustomTaggers, "custom-taggers", "c", []string{}, "paths to custom taggers plugins")
	flag.StringVarP(&commands.ExtraTags, "extra-tags", "e", "{}", "json dictionary format of extra tags to add to all taggable resources")
	flag.StringSliceVar(&commands.SkipConfigurationPaths, "skip-configuration-paths", []string{}, "configuration paths to skip")
}
