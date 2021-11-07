// spamassassin-parser is a service that shows how processing of reports works.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	log "github.com/obalunenko/logger"
	"github.com/obalunenko/version"

	"github.com/obalunenko/spamassassin-parser/cmd/spamassassin-parser/internal/config"
	"github.com/obalunenko/spamassassin-parser/internal/fileutil"
	"github.com/obalunenko/spamassassin-parser/internal/processor"
	"github.com/obalunenko/spamassassin-parser/pkg/utils"
)

func main() {
	ctx := context.Background()
	printVersion(ctx)

	config.Load()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Init(ctx, log.Params{
		Level:  config.LogLevel(),
		Format: config.LogFormat(),
		SentryParams: log.SentryParams{
			Enabled:      config.LogSentryEnabled(),
			DSN:          config.LogSentryDSN(),
			TraceEnabled: config.LogSentryTraceEnabled(),
			TraceLevel:   config.LogSentryTraceLevel(),
			Tags: map[string]string{
				"app_name":     version.GetAppName(),
				"go_version":   version.GetGoVersion(),
				"version":      version.GetVersion(),
				"build_date":   version.GetBuildDate(),
				"short_commit": version.GetShortCommit(),
			},
		},
	})

	defer log.Info(ctx, "Exit...")

	prccfg := processor.NewConfig()
	prccfg.Receive.Errors = config.ReceiveErrors()

	pr := processor.New(prccfg)

	go pr.Process(ctx)

	pollChan := make(chan fileutil.PollResponse)
	fileChan := make(chan string)

	go fileutil.PollDirectory(ctx, config.InputDir(), availableExtensions, pollChan)

	go putToProcessing(ctx, pr, config.InputDir(), fileChan)

	go func() {
		defer func() {
			close(fileChan)
		}()

		for resp := range pollChan {
			if resp.Err != nil {
				log.WithError(ctx, resp.Err).Error("poll error")
				cancel()

				return
			}

			fileChan <- resp.File
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	var wg sync.WaitGroup

	waitRoutinesNum := 1

	wg.Add(waitRoutinesNum)

	go process(ctx, &wg, pr, processParams{
		resultDir:  config.ResultDir(),
		archiveDir: config.ArchiveDir(),
		inputDir:   config.InputDir(),
	})

	select {
	case s := <-stopChan:
		log.WithField(ctx, "signal", s.String()).Info("Signal received")

		cancel()

	case <-ctx.Done():
		log.Info(ctx, "context canceled")
	}

	wg.Wait()
}

func putToProcessing(ctx context.Context, proc processor.Processor, inputDir string, fileChan <-chan string) {
	var finish bool

	for !finish {
		select {
		case <-ctx.Done():
			finish = true

		case reportFile, ok := <-fileChan:
			if !ok {
				log.Warn(ctx, "cmd: filechan is closed")

				finish = true
			}

			if reportFile == "" {
				log.Info(ctx, "cmd: empty report filename - wait for new")

				time.Sleep(1 * time.Second)

				continue
			}

			file, err := os.Open(filepath.Clean(filepath.Join(inputDir, reportFile)))
			if err != nil {
				log.WithError(ctx, err).Error("cmd: failed to open file with report")

				break
			}

			go func() {
				proc.Input() <- processor.NewInput(file, filepath.Base(file.Name()))
			}()
		}
	}
}

type processParams struct {
	resultDir  string
	archiveDir string
	inputDir   string
}

func process(ctx context.Context, wg *sync.WaitGroup, pr processor.Processor, params processParams) {
	defer wg.Done()

	for {
		select {
		case res := <-pr.Results():
			if res != nil {
				s, err := utils.PrettyPrint(res.Report, "", "\t")
				if err != nil {
					log.WithError(ctx, err).Error("failed to print report")
				}

				log.WithFields(ctx, log.Fields{
					"test_id": res.TestID,
				}).Info(fmt.Sprintf("Archive: \n %s \n", s))

				if err = fileutil.WriteFile(res.TestID, params.resultDir, s); err != nil {
					log.WithError(ctx, err).Error("failed to write file")

					continue
				}

				log.WithFields(ctx, log.Fields{
					"test_id":     res.TestID,
					"archive_dir": params.archiveDir,
				}).Info("Moving file to archive")

				if err = fileutil.MoveFile(ctx, res.TestID, params.inputDir, params.archiveDir); err != nil {
					log.WithError(ctx, err).Error("failed to move archive file")
				}

				log.Info(ctx, "File moved")
			}

		case err := <-pr.Errors():
			if err != nil {
				log.WithError(ctx, err).Error("processor error received")
			}

		case <-ctx.Done():
			log.Info(ctx, "cmd: context canceled")

			pr.Close()

			return
		}
	}
}

var availableExtensions = map[string]bool{
	"txt": true,
}
