// Package parser provides functionality to parse spamassassin result into json report.
package parser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"strings"

	"github.com/obalunenko/spamassassin-parser/internal/processor/models"
)

var emptyReport models.Report

// Parser is an interface that describes the basic functionality of a reports parser.
type Parser interface {
	// Parse takes a file as input and returns parsed from it report.
	Parse(data io.Reader) (models.Report, error)
}

func newParser(rt reportType) (Parser, error) {
	switch rt {
	case reportTypeSentinel, reportTypeUnknown:
		return nil, errors.New("invalid report type")
	case reportType1:
		return report1Parser{}, nil
	case reportType2:
		return report2Parser{}, nil
	default:
		return nil, fmt.Errorf("not implemented parser for report: %s", rt.String())
	}
}

//go:generate stringer -type=reportType

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
		return emptyReport, fmt.Errorf("failed to read from reader: %w", err)
	}

	rt := getReportType(bytes.NewReader(b))

	if !rt.Valid() {
		return emptyReport, fmt.Errorf("invalid report type: %w", err)
	}

	parser, err := newParser(rt)
	if err != nil {
		return emptyReport, fmt.Errorf("failed to get parser: %w", err)
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
