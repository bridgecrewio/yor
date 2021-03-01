package main

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/logger"
	"bridgecrewio/yor/src/common/reports"
	"bridgecrewio/yor/src/common/runner"
	"bridgecrewio/yor/src/common/tagging"
	"bridgecrewio/yor/src/common/tagging/code2cloud"
	"bridgecrewio/yor/src/common/tagging/gittag"
	"bridgecrewio/yor/src/common/tagging/simple"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {
	tagOptions := &common.TagOptions{}
	cmd := &cobra.Command{
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       common.Version,
		Short:         fmt.Sprintf("\nYor, the IaC auto-tagger (v%v)", common.Version),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.Help()
			if err == nil {
				os.Exit(0)
			}
			logger.Error(err.Error())
			return nil
		},
	}
	tagCmd := &cobra.Command{
		Use:           "tag",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "Tag you IaC files",
		RunE: func(cmd *cobra.Command, args []string) error {
			if tagOptions.Directory == "" {
				// If no flags supplied display the help menu and quit cleanly
				err := cmd.Help()
				if err == nil {
					os.Exit(0)
				}
				logger.Error(err.Error())
			}
			tagOptions.Validate()
			return tag(tagOptions)
		},
	}
	dtOptions := &common.DescribeTaggerOptions{}
	describeTaggerCmd := &cobra.Command{
		Use:           "describe-tagger",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "Get more details on each tagger",
		RunE: func(cmd *cobra.Command, args []string) error {
			if dtOptions.Tagger == "" {
				err := cmd.Help()
				if err == nil {
					os.Exit(0)
				}
				logger.Error(err.Error())
			}
			dtOptions.Validate()
			return describeTagger(dtOptions)
		},
	}
	listTaggersCmd := &cobra.Command{
		Use:           "list-taggers",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "List the available taggers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTaggers()
		},
	}
	addTagFlags(tagCmd.Flags(), tagOptions)
	addDescribeTaggerFlags(describeTaggerCmd.Flags(), dtOptions)
	cmd.AddCommand(tagCmd, describeTaggerCmd, listTaggersCmd)
	cmd.SetVersionTemplate(fmt.Sprintf("Yor version %s", cmd.Version))
	if err := cmd.Execute(); err != nil {
		logger.Error(err.Error())
	}
}

func tag(options *common.TagOptions) error {
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

func describeTagger(options *common.DescribeTaggerOptions) error {
	var tagger tagging.ITagger
	switch options.Tagger {
	case "simple":
		tagger = &simple.Tagger{}
	case "code2cloud":
		tagger = &code2cloud.Tagger{}
	case "git":
		tagger = &gittag.Tagger{}
	default:
		return fmt.Errorf("the tagger \"%s\" it not supported. Please try list-taggers to get a list of supported taggers", options.Tagger)
	}
	fmt.Println(tagger.GetDescription())
	return nil
}

func listTaggers() error {
	fmt.Println("Existing taggers:")
	for i, tagger := range common.SupportedTaggers {
		fmt.Printf("%d) %s\n", i+1, tagger)
	}
	return nil
}

func printReport(reportService *reports.ReportService, options *common.TagOptions) {
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

func addTagFlags(flag *pflag.FlagSet, commands *common.TagOptions) {
	flag.StringVarP(&commands.Directory, "directory", "d", "", "directory to tag")
	flag.StringVarP(&commands.Tag, "tag", "t", "", "run yor only with the specified tag")
	flag.StringVarP(&commands.SkipTag, "skip-tag", "s", "", "run yor without the specified tag")
	flag.StringVarP(&commands.Output, "output", "o", "cli", "set output format")
	flag.StringVar(&commands.OutputJSONFile, "output-json-file", "", "json file path for output")
	flag.StringSliceVarP(&commands.CustomTaggers, "custom-taggers", "c", []string{}, "paths to custom taggers plugins")
	flag.StringVarP(&commands.ExtraTags, "extra-tags", "e", "{}", "json dictionary format of extra tags to add to all taggable resources")
	flag.StringSliceVar(&commands.SkipConfigurationPaths, "skip-configuration-paths", []string{}, "configuration paths to skip")
}

func addDescribeTaggerFlags(flag *pflag.FlagSet, commands *common.DescribeTaggerOptions) {
	msg := "The tagger to be described. Valid values are "
	msg += "{" + strings.Join(common.SupportedTaggers, "|") + "}"
	flag.StringVarP(&commands.Tagger, "tagger", "t", "", msg)
}
