package resources

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var layout = "2006-01-02 15:04:05"

func getTime() time.Time {
	t, _ := time.Parse(layout, "2020-06-16 17:46:24")
	return t
}

func CreateComplexTagsLines() []*git.Line {
	originFileText, err := ioutil.ReadFile("../../tests/terraform/resources/complex_tags.tf")
	if err != nil {
		panic(err)
	}
	originLines := strings.Split(string(originFileText), "\n")
	lines := make([]*git.Line, 0)

	for _, line := range originLines {
		lines = append(lines, &git.Line{
			Author: "user@gmail.com",
			Text:   line,
			Date:   getTime(),
			Hash:   plumbing.NewHash("hash"),
		})
	}

	return lines
}
