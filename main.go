package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bridgecrewio/yor/src/common"
	"github.com/bridgecrewio/yor/src/common/clioptions"
	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/reports"
	"github.com/bridgecrewio/yor/src/common/runner"
	"github.com/bridgecrewio/yor/src/common/tagging"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/bridgecrewio/yor/src/common/tagging/utils"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                   "yor",
		HelpName:               "",
		Usage:                  "enrich IaC files with tags automatically",
		Version:                common.Version,
		Description:            "Yor, the IaC auto-tagger",
		Compiled:               time.Time{},
		Authors:                []*cli.Author{{Name: "Bridgecrew", Email: "support@bridgecrew.io"}},
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			listTagsCommand(),
			listTagGroupsCommand(),
			tagCommand(),
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.Error(err.Error())
	}
}

func listTagGroupsCommand() *cli.Command {
	return &cli.Command{
		Name:  "list-tag-groups",
		Usage: "List the tag groups that will be applied by yor",
		Action: func(c *cli.Context) error {
			return listTagGroups()
		},
	}
}

func listTagsCommand() *cli.Command {
	tagGroupsArg := "tag-groups"
	return &cli.Command{
		Name:  "list-tags",
		Usage: "List the tags yor will create if possible",
		Action: func(c *cli.Context) error {
			listTagsOptions := clioptions.ListTagsOptions{
				// cli package doesn't split comma separated values
				TagGroups: c.StringSlice(tagGroupsArg),
			}

			listTagsOptions.Validate()
			return listTags(&listTagsOptions)
		},
		Flags: []cli.Flag{ // When adding flags, make sure they are supported in the GitHub action as well via entrypoint.sh
			&cli.StringSliceFlag{
				Name:        tagGroupsArg,
				Aliases:     []string{"g"},
				Usage:       "Filter results by specific tag group(s), comma delimited",
				Value:       cli.NewStringSlice(utils.GetAllTagGroupsNames()...),
				DefaultText: strings.Join(utils.GetAllTagGroupsNames(), ","),
			},
		},
		HideHelpCommand:        true,
		UseShortOptionHandling: true,
	}
}

func tagCommand() *cli.Command {
	directoryArg := "directory"
	tagArg := "tags"
	skipTagsArg := "skip-tags"
	customTaggingArg := "custom-tagging"
	skipDirsArg := "skip-dirs"
	outputArg := "output"
	tagGroupArg := "tag-groups"
	outputJSONFileArg := "output-json-file"
	externalConfPath := "config-file"
	skipResourceTypesArg := "skip-resource-types"
	skipResourcesArg := "skip-resources"
	parsersArgs := "parsers"
	dryRunArgs := "dry-run"
	tagLocalModules := "tag-local-modules"
	tagPrefix := "tag-prefix"
	return &cli.Command{
		Name:                   "tag",
		Usage:                  "apply tagging across your directory",
		HideHelpCommand:        true,
		UseShortOptionHandling: true,
		Action: func(c *cli.Context) error {
			options := clioptions.TagOptions{
				Directory:         c.String(directoryArg),
				Tag:               c.StringSlice(tagArg),
				SkipTags:          c.StringSlice(skipTagsArg),
				CustomTagging:     c.StringSlice(customTaggingArg),
				SkipDirs:          c.StringSlice(skipDirsArg),
				Output:            c.String(outputArg),
				OutputJSONFile:    c.String(outputJSONFileArg),
				TagGroups:         c.StringSlice(tagGroupArg),
				ConfigFile:        c.String(externalConfPath),
				SkipResourceTypes: c.StringSlice(skipResourceTypesArg),
				SkipResources:     c.StringSlice(skipResourcesArg),
				Parsers:           c.StringSlice(parsersArgs),
				DryRun:            c.Bool(dryRunArgs),
				TagLocalModules:   c.Bool(tagLocalModules),
				TagPrefix:         c.String(tagPrefix),
			}

			options.Validate()

			return tag(&options)
		},
		Flags: []cli.Flag{ // When adding flags, make sure they are supported in the GitHub action as well via entrypoint.sh
			&cli.StringFlag{
				Name:        directoryArg,
				Aliases:     []string{"d"},
				Usage:       "directory to tag",
				Required:    true,
				DefaultText: "path/to/iac/root",
			},
			&cli.StringSliceFlag{
				Name:        tagArg,
				Aliases:     []string{"t"},
				Usage:       "run yor only with the specified tags",
				DefaultText: "yor_trace,git_repository",
			},
			&cli.StringSliceFlag{
				Name:        skipTagsArg,
				Aliases:     []string{"s"},
				Usage:       "run yor skipping the specified tags",
				Value:       cli.NewStringSlice(),
				DefaultText: "yor_trace",
			},
			&cli.StringFlag{
				Name:        outputArg,
				Aliases:     []string{"o"},
				Usage:       "set output format",
				Value:       "cli",
				DefaultText: "json",
			},
			&cli.StringFlag{
				Name:        outputJSONFileArg,
				Usage:       "json file path for output",
				DefaultText: "result.json",
			},
			&cli.StringSliceFlag{
				Name:        customTaggingArg,
				Aliases:     []string{"c"},
				Usage:       "paths to custom tag groups and tags plugins",
				Value:       cli.NewStringSlice(),
				DefaultText: "path/to/custom/yor/tagging",
			},
			&cli.StringSliceFlag{
				Name:        skipDirsArg,
				Aliases:     nil,
				Usage:       "configuration paths to skip",
				Value:       cli.NewStringSlice(),
				DefaultText: "path/to/skip,another/path/to/skip",
			},
			&cli.StringSliceFlag{
				Name:        tagGroupArg,
				Aliases:     []string{"g"},
				Usage:       "Narrow down results to the matching tag groups",
				Value:       cli.NewStringSlice(utils.GetAllTagGroupsNames()...),
				DefaultText: "git,code2cloud",
			},
			&cli.StringFlag{
				Name:        externalConfPath,
				Usage:       "external tag group configuration file path",
				DefaultText: "/path/to/conf/file/ (.yml/.yaml extension)",
			},
			&cli.StringSliceFlag{
				Name:        skipResourceTypesArg,
				Usage:       "skip resource types for tagging",
				Value:       cli.NewStringSlice(),
				DefaultText: "aws_rds_instance,AWS::S3::Bucket",
			},
			&cli.StringSliceFlag{
				Name:        skipResourcesArg,
				Usage:       "skip resources for tagging",
				Value:       cli.NewStringSlice(),
				DefaultText: "aws_s3_bucket.test-bucket,EC2InstanceResource0",
			},
			&cli.StringSliceFlag{
				Name:        parsersArgs,
				Aliases:     []string{"i"},
				Usage:       "IAC types to tag",
				Value:       cli.NewStringSlice("Terraform", "CloudFormation", "Serverless"),
				DefaultText: "Terraform,CloudFormation,Serverless",
			},
			&cli.BoolFlag{
				Name:        dryRunArgs,
				Usage:       "skip resource tagging",
				Value:       false,
				DefaultText: "false",
			},
			&cli.BoolFlag{
				Name:        tagLocalModules,
				Usage:       "Always tag local modules",
				Value:       false,
				DefaultText: "false",
			},
			&cli.StringFlag{
				Name:        tagPrefix,
				Usage:       "Add prefix to all the tags",
				DefaultText: "",
			},
		},
	}
}

func listTagGroups() error {
	for _, tagGroup := range utils.GetAllTagGroupsNames() {
		fmt.Println(tagGroup)
	}
	return nil
}

func listTags(options *clioptions.ListTagsOptions) error {
	var tagGroup tagging.ITagGroup
	tagsByGroup := make(map[string][]tags.ITag)
	for _, group := range options.TagGroups {
		tagGroup = utils.TagGroupsByName(utils.TagGroupName(group))
		if tagGroup == nil {
			return fmt.Errorf("tag group %v is not supported", group)
		}
		tagGroup.InitTagGroup("", nil, nil)
		tagsByGroup[group] = tagGroup.GetTags()
	}
	reports.ReportServiceInst.PrintTagGroupTags(tagsByGroup)
	return nil
}

func tag(options *clioptions.TagOptions) error {
	yorRunner := new(runner.Runner)
	logger.Info(fmt.Sprintf("Setting up to tag the directory %v\n", options.Directory))
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

func printReport(reportService *reports.ReportService, options *clioptions.TagOptions) {
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
