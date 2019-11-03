package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"

	"github.com/oleg-balunenko/spamassassin-parser/pkg/parser"
	"github.com/oleg-balunenko/spamassassin-parser/pkg/utils"
)

var (
	reportFile = flag.String("report_file", "", "path to report file to process")
)

func main() {
	printVersion()

	flag.Parse()

	if *reportFile == "" {
		log.Fatal("report_file not set")
	}

	file, err := os.Open(*reportFile)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to open file with report"))
	}

	report, err := parser.ProcessReport(file)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to process report"))
	}

	s, err := utils.PrettyPrint(report)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to print report"))
	}
	fmt.Println(s)
}
