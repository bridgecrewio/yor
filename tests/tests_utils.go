package tests

import (
	"bufio"
	"github.com/go-git/go-git/v5"
	"io/ioutil"
	"log"
	"os"
)

func CloneRepo(repoPath string) string {
	dir, err := ioutil.TempDir("", "temp-repo")
	if err != nil {
		log.Fatal(err)
	}

	// Clones the repository into the given dir, just as a normal git clone does
	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL: repoPath,
	})

	if err != nil {
		log.Fatal(err)
	}

	return dir
}

func ReadFileLines(filePath string) ([]string, error) {
	expectedBlameFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	fileScanner := bufio.NewScanner(expectedBlameFile)
	fileScanner.Split(bufio.ScanLines)
	fileTextLines := []string{}
	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}
	expectedBlameFile.Close()

	return fileTextLines, nil
}
