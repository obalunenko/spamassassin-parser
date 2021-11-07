// Package fileutil provide functionality for working with files and directories.
package fileutil

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/obalunenko/logger"
)

// WriteFile creates file with passed name at passed dir and writes data.
func WriteFile(fname, dir, data string) error {
	if ext := filepath.Ext(filepath.Base(fname)); ext != "" {
		fname = strings.TrimSuffix(fname, ext)
	}

	outName := filepath.Join(dir, fmt.Sprintf("%s.json", fname))

	if err := createDirIfNotExist(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	resFile, err := os.Create(outName)
	if err != nil {
		return fmt.Errorf("failed to create result file [%s]: %w", outName, err)
	}

	if _, err = resFile.WriteString(data); err != nil {
		return fmt.Errorf("failed to write result to file: %w", err)
	}

	if err = resFile.Close(); err != nil {
		return fmt.Errorf("failed to close result file: %w", err)
	}

	return nil
}

// MoveFile moves file from base dir to target.
func MoveFile(ctx context.Context, name, sourceDir, destDir string) error {
	sourcePath := filepath.Join(sourceDir, name)

	inputFile, err := os.Open(filepath.Clean(sourcePath))
	if err != nil {
		return fmt.Errorf("couldn't open source file: %w", err)
	}

	defer func() {
		if err = inputFile.Close(); err != nil {
			log.WithError(ctx, err).Error("failed to close inputFile")
		}
	}()

	if err = createDirIfNotExist(destDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create destPath: %w", err)
	}

	destPath := filepath.Join(destDir, name)

	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("couldn't open dest file: %w", err)
	}

	defer func() {
		if err = outputFile.Close(); err != nil {
			log.WithError(ctx, err).Error("failed to close outputFile")
		}
	}()

	if _, err = io.Copy(outputFile, inputFile); err != nil {
		return fmt.Errorf("writing to output file failed: %w", err)
	}

	// The copy was successful, so now delete the original file
	if err = os.Remove(sourcePath); err != nil {
		return fmt.Errorf("failed removing original file: %w", err)
	}

	return nil
}

func createDirIfNotExist(dir string, mode os.FileMode) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, mode); err != nil {
			return fmt.Errorf("failed to create directory for [%s]: %w", dir, err)
		}
	}

	return nil
}

// PollResponse represents response of directory poll. Filename of new file if it appears, and error if smth went wrong.
type PollResponse struct {
	File string
	Err  error
}

func newPollResponse(file string, err error) PollResponse {
	return PollResponse{
		File: file,
		Err:  err,
	}
}

type poller struct {
	lastFile            string
	availableExtensions map[string]bool
}

func newPoller(availableExtensions map[string]bool) poller {
	return poller{
		lastFile:            "",
		availableExtensions: availableExtensions,
	}
}

func (p *poller) pollfiles(ctx context.Context, files []fs.FileInfo, respChan chan<- PollResponse) {
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if file.Name() == "" {
			continue
		}

		ext := filepath.Ext(file.Name())
		if ext != "" {
			ext = strings.TrimPrefix(ext, ".")
		}

		if file.Name() == p.lastFile {
			respChan <- newPollResponse("",
				fmt.Errorf("cannot porcess the last known file, hanged"))

			log.WithFields(ctx, log.Fields{
				"last_file": p.lastFile,
			}).Warn("internal/fileutil: Cannot process the last known file")
		}

		_, ok := p.availableExtensions[ext]
		if ok && file.Name() != p.lastFile {
			respChan <- newPollResponse(file.Name(), nil)

			p.lastFile = file.Name()

			break
		}
	}
}

// PollDirectory will pol for files that needs to be processed.
func PollDirectory(ctx context.Context, dir string, availableExtensions map[string]bool, respChan chan<- PollResponse) {
	const pollInterval int = 1

	ticker := time.NewTicker(time.Duration(pollInterval) * time.Second)

	defer func() {
		close(respChan)

		ticker.Stop()
	}()

	log.WithFields(ctx, log.Fields{
		"dir": dir,
	}).Info("started to watch directory for data files")

	p := newPoller(availableExtensions)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				respChan <- newPollResponse("",
					fmt.Errorf("internal/fileutil: readdir: %w", err))

				return
			}

			p.pollfiles(ctx, files, respChan)
		}
	}
}
