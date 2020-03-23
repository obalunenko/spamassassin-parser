// Package parser provides functionality to parse spamassassin result into json report.
package parser

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"math"
	"strings"

	"github.com/pkg/errors"

	"github.com/oleg-balunenko/spamassassin-parser/pkg/models"
)

var (
	emptyReport models.Report
)

// Parser is an interface that describes the basic functionality of a reports parser.
type Parser interface {
	// Parse takes a file as input and returns parsed from it report.
	Parse(data io.Reader) (models.Report, error)
}

func newParser(rt reportType) (Parser, error) {
	switch rt {
	case reportType1:
		return report1Parser{}, nil
	case reportType2:
		return report2Parser{}, nil
	}

	return nil, errors.Errorf("not implemented parser for report: %s", rt.String())
}

//go:generate stringer -type=reportType

// reportType
type reportType int

const (
	reportTypeUnknown reportType = iota
	reportType1
	reportType2

	reportTypeSentinel // should be always last.
)

func (i reportType) Valid() bool {
	return i > reportTypeUnknown && i < reportTypeSentinel
}

// ParseReport parses passed raw report to json representation.
func ParseReport(data io.Reader) (models.Report, error) {
	b, err := ioutil.ReadAll(data)
	if err != nil {
		return emptyReport, errors.Wrap(err, "failed to read from reader")
	}

	rt := getReportType(bytes.NewReader(b))

	if !rt.Valid() {
		return emptyReport, errors.New("invalid report type")
	}

	parser, err := newParser(rt)
	if err != nil {
		return emptyReport, errors.Wrap(err, "failed to get parser")
	}

	return parser.Parse(bytes.NewReader(b))
}

func getReportType(data io.Reader) reportType {
	sc := bufio.NewScanner(data)
	for sc.Scan() {
		line := sc.Text()
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "*"):
			return reportType1
		case strings.Contains(strings.ToLower(line), strings.ToLower("Spam detection software")):
			return reportType2
		}
	}

	return reportTypeUnknown
}

func sanitizeScore(f float64) float64 {
	if f == -0 {
		f = 0
	}

	var scaleHundred float64 = 100

	return math.Round(f*scaleHundred) / scaleHundred
}
