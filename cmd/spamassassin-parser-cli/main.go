// spamassassin-parser-cli is a command line tool that shows how processing of reports works.
package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/oleg-balunenko/spamassassin-parser/internal/fileutil"
	"github.com/oleg-balunenko/spamassassin-parser/pkg/models"
	"github.com/oleg-balunenko/spamassassin-parser/pkg/processor"
	"github.com/oleg-balunenko/spamassassin-parser/pkg/utils"
)

var (
	inputDir = flag.String("input_dir", "input",
		"Path to directory where files for proccession are located")
	outputDir = flag.String("output_dir", "output",
		"Path to directory where parserd results will be stored")
	processedDir = flag.String("processed_dir", "archive",
		"Path to dir where processed files will be moved for history")
)

func main() {
	defer log.Println("Exit...")

	printVersion()

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := processor.NewConfig()
	cfg.Receive.Errors = true
	pr := processor.NewProcessor(cfg)

	go pr.Process(ctx)

	fileChan := make(chan string)
	go fileutil.PollDirectory(ctx, *inputDir, availableExtensions, fileChan)

	go func(ctx context.Context, fileChan chan string) {
		for {
			select {
			case <-ctx.Done():
				return

			case reportFile := <-fileChan:
				file, err := os.Open(filepath.Clean(filepath.Join(*inputDir, reportFile)))
				if err != nil {
					log.Fatal(errors.Wrap(err, "failed to open file with report"))
				}

				go func() {
					pr.Input() <- models.NewProcessorInput(file, filepath.Base(file.Name()))
				}()
			}
		}
	}(ctx, fileChan)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	var wg sync.WaitGroup

	waitRoutinesNum := 1

	wg.Add(waitRoutinesNum)

	go process(ctx, &wg, pr, *outputDir, *processedDir)

	s := <-stopChan
	log.Infof("Signal [%s] received", s.String())

	cancel()

	wg.Wait()
}

func process(ctx context.Context, wg *sync.WaitGroup, pr processor.Processor, outDir string, processedDir string) {
	defer wg.Done()

	for {
		select {
		case res := <-pr.Results():
			if res != nil {
				s, err := utils.PrettyPrint(res.Report, "", "\t")
				if err != nil {
					log.Error(errors.Wrap(err, "failed to print report"))
				}

				log.Printf("[TestID: %s] processed: \n %s \n",
					res.TestID, s)

				if err = fileutil.WriteFile(res.TestID, outDir, s); err != nil {
					log.Error(errors.Wrap(err, "failed to write file"))
				}

				log.Infof("Moving file %s to archive folder: %s", res.TestID, processedDir)

				if err = fileutil.MoveFileToFolder(res.TestID, *inputDir, processedDir); err != nil {
					log.Error(errors.Wrap(err, "failed to move processed file"))
				}

				log.Info("File moved")
			}

		case err := <-pr.Errors():
			if err != nil {
				log.Error(err)
			}

		case <-ctx.Done():
			log.Println("context canceled")

			pr.Close()

			return
		}
	}
}

var availableExtensions = map[string]bool{
	"txt": true,
}
