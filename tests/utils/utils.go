package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

const TerragoatURL = "https://github.com/bridgecrewio/terragoat.git"

func CaptureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(stderr)
	}()
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)

	f()
	_ = writer.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, reader)
	_ = reader.Close()
	return buf.String()
}

func CloneRepo(repoPath string, commitHash string) string {
	dir, err := ioutil.TempDir("", "temp-repo")
	if err != nil {
		log.Fatal(err)
	}

	// Clones the repository into the given dir, just as a normal git clone does
	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: repoPath,
	})

	if commitHash != "" {
		wt, _ := repo.Worktree()

		commitRef := plumbing.NewHash(commitHash)
		_ = wt.Checkout(&git.CheckoutOptions{Hash: commitRef})
	}

	if err != nil {
		log.Fatal(err)
	}

	return dir
}

func CopyDirToTempDir(originDir string) string {
	tempDir, err := ioutil.TempDir("", "temp-repo")
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("cp", "-R", originDir, tempDir)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return tempDir
}
