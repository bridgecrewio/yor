package utils

import (
	"bytes"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

func ExtractDate(dateStr string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.000Z"
	parsedDate, err := time.Parse(layout, dateStr)
	return parsedDate, err
}

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
		log.SetOutput(os.Stderr)
	}()
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		_, _ = io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	writer.Close()
	return <-out
}
