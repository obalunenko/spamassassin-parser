package fileutil

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// WriteFile creates file with passed name at passed dir and writes data.
func WriteFile(fname, dir string, data string) error {
	if ext := filepath.Ext(filepath.Base(fname)); ext != "" {
		fname = strings.TrimSuffix(fname, ext)
	}

	outName := filepath.Join(dir, fmt.Sprintf("%s.json", fname))

	if err := createDirIfNotExist(dir, os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to create directory")
	}

	resFile, err := os.Create(outName)
	if err != nil {
		return errors.Wrapf(err, "failed to create result file [%s]", outName)
	}

	if _, err = resFile.WriteString(data); err != nil {
		return errors.Wrap(err, "failed to write result to file")
	}

	if err = resFile.Close(); err != nil {
		return errors.Wrap(err, "failed to close result file")
	}

	return nil
}

// MoveFileToFolder moves file from base dir to target
func MoveFileToFolder(filename, fromDir, targetDir string) error {
	filename = filepath.Base(filename)

	if err := createDirIfNotExist(targetDir, os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to create  dir")
	}

	targetFilename := filepath.Join(targetDir, filename+".tmp")
	currentFilename := filepath.Join(fromDir, filename)

	if err := os.Rename(currentFilename, targetFilename); err != nil {
		log.Warnf("Failed to move %s to folder: %v, does the desired folder exist?", filename, err)
		return err
	}

	return nil
}

func createDirIfNotExist(dir string, mode os.FileMode) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, mode); err != nil {
			return errors.Wrapf(err, "failed to create directory for [%s]", dir)
		}
	}

	return nil
}

// PollDirectory will pol for files that needs to be processed.
func PollDirectory(ctx context.Context, dir string, availableExtensions map[string]bool, fileChan chan string) {
	const pollInterval int = 1

	ticker := time.NewTicker(time.Duration(pollInterval) * time.Second)

	defer func() {
		close(fileChan)
		ticker.Stop()
	}()

	log.Printf("started to watch directory <%s> for data files", dir)

	var lastFile string

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				log.Error(err)

				return
			}

			for _, file := range files {
				if !file.IsDir() {
					ext := filepath.Ext(file.Name())
					if ext != "" {
						ext = strings.TrimPrefix(ext, ".")
					}

					if file.Name() == lastFile {
						log.Warnf("Cannot process the last known file: %s", lastFile)
					}

					_, ok := availableExtensions[ext]
					if ok && file.Name() != lastFile {
						fileChan <- file.Name()

						lastFile = file.Name()

						break
					}
				}
			}
		}
	}
}
