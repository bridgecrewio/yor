package main

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/logger"
	"bridgecrewio/yor/src/common/reports"
	"bridgecrewio/yor/src/common/runner"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {
	options := &common.Options{}
	cmd := &cobra.Command{
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       common.Version,
		Short:         fmt.Sprintf("\nYor, the IaC auto-tagger (v%v)", common.Version),
		RunE: func(cmd *cobra.Command, args []string) error {
			if options.Directory == "" {
				// If no flags supplied display the help menu and quit cleanly
				err := cmd.Help()
				if err == nil {
					os.Exit(0)
				}
				logger.Error(err.Error())
			}
			return run(options)
		},
	}
	tagCmd := &cobra.Command{
		Use:           "tag",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "Tag you IaC files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(options)
		},
	}
	addTagFlags(tagCmd.Flags(), options)
	cmd.AddCommand(tagCmd)
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

	if options.OutputJSONFile != "" {
		reportService.PrintJSONToFile(options.OutputJSONFile)
	}
	switch strings.ToLower(options.Output) {
	case "cli":
		reportService.PrintToStdout()
	case "json":
		reportService.PrintJSONToStdout()
	default:
		return
	}
}

func addTagFlags(flag *pflag.FlagSet, commands *common.Options) {
	flag.StringVarP(&commands.Directory, "directory", "d", "", "directory to tag")
	flag.StringVarP(&commands.Tag, "tag", "t", "", "run yor only with the specified tag")
	flag.StringVarP(&commands.SkipTag, "skip-tag", "s", "", "run yor without ths specified tag")
	flag.StringVarP(&commands.Output, "output", "o", "cli", "set output format")
	flag.StringVar(&commands.OutputJSONFile, "output-json-file", "", "json file path for output")
	flag.StringSliceVarP(&commands.CustomTaggers, "custom-taggers", "c", []string{}, "paths to custom taggers plugins")
	flag.StringVarP(&commands.ExtraTags, "extra-tags", "e", "{}", "json dictionary format of extra tags to add to all taggable resources")
	flag.StringSliceVar(&commands.SkipConfigurationPaths, "skip-configuration-paths", []string{}, "configuration paths to skip")
}
