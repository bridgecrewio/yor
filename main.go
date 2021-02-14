package main

import (
	"bridgecrewio/yor/common/reports"
	"bridgecrewio/yor/common/tagging/tags"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

func main() {
	fmt.Println("Welcome to Yor!")
	loadExternalTags("main")
}

func parseArgs(args ...interface{}) {
	// TODO
}

func printReport(report *reports.Report) {
	// TODO
}

func createExtraTags(extraTagsFromArgs map[string]string) []tags.ITag {
	extraTags := make([]tags.ITag, len(extraTagsFromArgs))
	index := 0
	for key := range extraTagsFromArgs {
		newTag := tags.Init(key, extraTagsFromArgs[key])
		extraTags[index] = newTag
		index++
	}

	return extraTags
}

func loadExternalTags(tagsPath string) {
	plugins := []string{}
	err := filepath.Walk(tagsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".so") {
			plugins = append(plugins, path)
		}
		return nil
	})
	//goPath, err := exec.LookPath("go")
	if err != nil {
		panic(err)
	}
	//cmd := exec.Command(goPath, "build", "main/yor_tags.go")
	////cmd := exec.Command(goPath, "build", "buildmode", "plugin", "o", "/Users/rotemavni/BridgeCrew/yor_tags/yor_tags.so", "/Users/rotemavni/BridgeCrew/yor_tags/yor_tags.go")
	////cmd := exec.Command(goPath, "build", "-gcflags", "\"all=-N -l\"", " -buildmode", "plugin", "-o", "/Users/rotemavni/BridgeCrew/yor_tags/yor_tags.so", "/Users/rotemavni/BridgeCrew/yor_tags/yor_tags.go")
	//msg, err := cmd.Output()
	//if err != nil {
	//	panic(string(msg))
	//}
	//

	for _, pluginPath := range plugins {
		plug, err := plugin.Open(pluginPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		//f, err := plug.Lookup("GetTag")
		//if err != nil {
		//	fmt.Println(err)
		//	os.Exit(1)
		//}
		//var symGreeter plugin.Symbol
		//symGreeter = f.(func() interface{})()

		// 2. look up a symbol (an exported function or variable)
		// in this case, variable Greeter
		symGreeter, err := plug.Lookup("ITag")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 3. Assert that loaded symbol is of a desired type
		// in this case interface type Greeter (defined above)
		var iTag tags.ITag
		iTag, ok := symGreeter.(tags.ITag)
		if !ok {
			fmt.Println("unexpected type from module symbol")
			os.Exit(1)
		}

		iTag.Init()
		iTag.CalculateValue(nil)

		key := iTag.GetKey()
		value := iTag.GetValue()
		// 4. use the module
		fmt.Printf("key: %s, value: %s\n", key, value)
	}

}
