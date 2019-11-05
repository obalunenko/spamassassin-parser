package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/oleg-balunenko/spamassassin-parser/pkg/models"
	"github.com/oleg-balunenko/spamassassin-parser/pkg/processor"
	"github.com/oleg-balunenko/spamassassin-parser/pkg/utils"
)

var (
	reportFile = flag.String("report_file", "", "path to report file to process")
)

func main() {
	defer log.Println("Exit...")

	printVersion()

	flag.Parse()

	if *reportFile == "" {
		log.Fatal("report_file not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	cfg := processor.NewConfig()
	cfg.Receive.Errors = true
	pr := processor.NewProcessor(cfg)

	go pr.Process(ctx)

	file, err := os.Open(*reportFile)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to open file with report"))
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	go func() {
		pr.Input() <- models.NewProcessorInput(file, file.Name())
	}()

LOOP:
	for {
		select {
		case res := <-pr.Results():
			if res != nil {
				s, err := utils.PrettyPrint(res.Report, "", "\t")
				if err != nil {
					log.Fatal(errors.Wrap(err, "failed to print report"))
				}
				log.Printf("[TestID: %s] processed: \n %s \n",
					res.TestID, s)
			}

		case err := <-pr.Errors():
			if err != nil {
				log.Fatal(err)
			}
		case <-ctx.Done():
			log.Println("context deadline")
			pr.Close()
			break LOOP
		case <-stopChan:
			log.Println("ctrl+c received")
			pr.Close()
			break LOOP
		}
	}
}
