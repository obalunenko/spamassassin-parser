package parser

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

var EmptyReport Report

var errNotImplemented = errors.New("not implemented")

func ProcessReport(data io.Reader) (Report, error) {
	b, err := ioutil.ReadAll(data)
	if err != nil {
		return EmptyReport, errors.Wrap(err, "failed to read from reader")
	}

	rt, err := getReportType(bytes.NewReader(b))
	if err != nil {
		return Report{}, errors.Wrap(err, "failed to get report type")
	}

	return processReport(bytes.NewReader(b), rt)
}

func getReportType(data io.Reader) (reportType, error) {
	sc := bufio.NewScanner(data)
	for sc.Scan() {
		// line := sc.Text()
	}
	return reportTypeUnknown, errNotImplemented
}

func processReport(data io.Reader, rt reportType) (Report, error) {
	if !rt.Valid() {
		return Report{}, errors.New("invalid report type")
	}
	return EmptyReport, errNotImplemented
}
