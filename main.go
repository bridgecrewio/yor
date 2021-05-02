package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bridgecrewio/yor/src/common"
	"github.com/bridgecrewio/yor/src/common/cli"
	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/reports"
	"github.com/bridgecrewio/yor/src/common/runner"
	"github.com/bridgecrewio/yor/src/common/tagging"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/bridgecrewio/yor/src/common/tagging/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {
	tagOptions := &cli.TagOptions{}
	cmd := &cobra.Command{
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       common.Version,
		Short:         fmt.Sprintf("\nYor, the IaC auto-tagger (v%v)", common.Version),
		RunE: func(cmd *cobra.Command, args []string) error {
			if tagOptions.Directory == "" {
				// If no flags supplied display the help menu and quit cleanly
				err := cmd.Help()
				if err == nil {
					os.Exit(0)
				}
				logger.Error(err.Error())
			}
			return run(tagOptions)
		},
	}
	tagCmd := &cobra.Command{
		Use:           "tag",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "Tag you IaC files",
		RunE: func(cmd *cobra.Command, args []string) error {
			tagOptions.Validate()
			return run(tagOptions)
		},
	}
	addTagFlags(tagCmd.Flags(), tagOptions)

	listTagsOptions := &cli.ListTagsOptions{}
	listTagsCmd := &cobra.Command{
		Use:           "list-tags",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "List the tags supported by Yor",
		RunE: func(cmd *cobra.Command, args []string) error {
			listTagsOptions.Validate()
			return listTags(listTagsOptions)
		},
	}
	addListTagsFlags(listTagsCmd.Flags(), listTagsOptions)

	listTagGroups := &cobra.Command{
		Use:           "list-tag-groups",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "List the Tag Groups supported by Yor",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTagGroups()
		},
	}
	cmd.AddCommand(tagCmd, listTagsCmd, listTagGroups)

	cmd.SetVersionTemplate(fmt.Sprintf("Yor version %s", cmd.Version))
	if err := cmd.Execute(); err != nil {
		logger.Error(err.Error())
	}
}

func listTagGroups() error {
	for _, tagGroup := range utils.GetAllTagGroupsNames() {
		fmt.Println(tagGroup)
	}
	return nil
}

func listTags(options *cli.ListTagsOptions) error {
	var tagGroup tagging.ITagGroup
	tagsByGroup := make(map[string][]tags.ITag)
	for _, group := range options.TagGroups {
		tagGroup = utils.TagGroupsByName(utils.TagGroupName(group))
		if tagGroup == nil {
			return fmt.Errorf("tag group %v is not supported", group)
		}
		tagGroup.InitTagGroup("", nil)
		tagsByGroup[group] = tagGroup.GetTags()
	}
	reports.ReportServiceInst.PrintTagGroupTags(tagsByGroup)
	return nil
}

func run(options *cli.TagOptions) error {
	yorRunner := new(runner.Runner)
	logger.Info("Setting up to tag the directory %v\n", options.Directory)
	err := yorRunner.Init(options)
	if err != nil {
		logger.Error(err.Error())
	}
	reportService, err := yorRunner.TagDirectory()
	if err != nil {
		logger.Error(err.Error())
	}
	printReport(reportService, options)

	return nil
}

func printReport(reportService *reports.ReportService, options *cli.TagOptions) {
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

func addTagFlags(flag *pflag.FlagSet, options *cli.TagOptions) {
	flag.StringVarP(&options.Directory, "directory", "d", "", "directory to tag")
	flag.StringVarP(&options.Tag, "tag", "t", "", "run yor only with the specified tag")
	flag.StringSliceVarP(&options.SkipTags, "skip-tags", "s", []string{}, "run yor without ths specified tag")
	flag.StringVarP(&options.Output, "output", "o", "cli", "set output format")
	flag.StringVar(&options.OutputJSONFile, "output-json-file", "", "json file path for output")
	flag.StringSliceVarP(&options.CustomTagging, "custom-tagging", "c", []string{}, "paths to custom tag groups and tags plugins")
	flag.StringSliceVar(&options.SkipDirs, "skip-dirs", []string{}, "configuration paths to skip")
	flag.StringSliceVarP(&options.TagGroups, "tag-groups", "g", utils.GetAllTagGroupsNames(), "Narrow down results to the matching tag groups")
}

func addListTagsFlags(flag *pflag.FlagSet, options *cli.ListTagsOptions) {
	flag.StringSliceVarP(&options.TagGroups, "tag-groups", "", utils.GetAllTagGroupsNames(), "Narrow down results to the matching tag groups")
}
