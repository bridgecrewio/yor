package reports

import (
	"bridgecrewio/yor/common/structure"
	"fmt"
	"strings"
)

type ReportService struct {
	report Report
}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

type Report struct {
	ScannedResources int
	NewResources     []structure.IBlock
	UpdatedResources []structure.IBlock
}

var ReportServiceInst *ReportService

func init() {
	ReportServiceInst = &ReportService{}
}

func (r *ReportService) GetReport() *Report {
	return &r.report
}

func (r *ReportService) CreateReport() *Report {
	changesAccumulator := TagChangeAccumulatorInstance
	r.report.ScannedResources = len(changesAccumulator.scannedBlocks)
	r.report.NewResources = changesAccumulator.newBlockTraces
	r.report.UpdatedResources = changesAccumulator.updatedBlockTraces
	return &r.report
}

func (r *ReportService) PrintToStdout() {
	fmt.Println(colorReset, "Yor Findings Summary")
	fmt.Println(colorReset, "Scanned Resources:\t\t", colorBlue, r.report.ScannedResources)
	fmt.Println(colorReset, "Updated Resources:\t\t", colorGreen, len(r.report.UpdatedResources))
	fmt.Println(colorReset, "New Resources Traced: \t", colorYellow, len(r.report.NewResources))

	fmt.Println()

	if len(r.report.NewResources) > 0 {
		fmt.Print(colorReset, fmt.Sprintf("New Resources Traced (%v):\n", len(r.report.NewResources)))
		fmt.Println(strings.Repeat("-", 80))
		for index, block := range r.report.NewResources {
			fmt.Printf("%v:\tResource: %v:%v\n", index+1, block.GetFilePath(), block.GetResourceId())
			fmt.Printf("\tOwner: %v\n", block.GetNewOwner())
			fmt.Printf("\tTrace ID: %v\n", block.GetTraceId())
			fmt.Println()
		}
	}

	fmt.Println()

	fmt.Println(strings.Repeat("-", 80))
	if len(r.report.UpdatedResources) > 0 {
		fmt.Print(colorReset, fmt.Sprintf("Updated Resource Traces (%v):\n", len(r.report.UpdatedResources)))
		for index, block := range r.report.UpdatedResources {
			fmt.Printf("%v:\tResource: %v:%v\n", index+1, block.GetFilePath(), block.GetResourceId())
			fmt.Printf("\tPrevious Owner: %v\n", block.GetPreviousOwner())
			fmt.Printf("\tNew Owner: %v\n", block.GetNewOwner())
			fmt.Printf("\tTrace ID: %v\n", block.GetTraceId())
			fmt.Println()
		}
	}
}
