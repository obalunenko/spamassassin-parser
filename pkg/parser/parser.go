package parser

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var emptyReport Report

var errNotImplemented = errors.New("not implemented")
var (
	reType1 = regexp.MustCompile(`([*])[\s]+([-]?\d.\d)?[\s](([[:word:]]+)?[\s](.*))`)
	reType2 = regexp.MustCompile(`(?m)([-]?\d.\d)[\s]+([[:word:]]+)\s+(.*[\n]?)`)
)

// ProcessReport parses passed raw report to json representation.
func ProcessReport(data io.Reader) (Report, error) {
	b, err := ioutil.ReadAll(data)
	if err != nil {
		return emptyReport, errors.Wrap(err, "failed to read from reader")
	}

	rt := getReportType(bytes.NewReader(b))

	return processReport(bytes.NewReader(b), rt)
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

func processReport(data io.Reader, rt reportType) (Report, error) {
	if !rt.Valid() {
		return Report{}, errors.New("invalid report type")
	}
	switch rt {
	case reportType1:
		return processReportType1(data)

	case reportType2:
		return processReportType2(data)
	}

	return emptyReport, errNotImplemented
}

func processReportType1(data io.Reader) (Report, error) {
	const (
		colFullMatch = iota
		colAsterisk
		colScore
		colTagWithDescr
		colTag
		colDescr
	)
	var (
		r     Report
		score float64
		lnum  int
	)
	sc := bufio.NewScanner(data)
	for sc.Scan() {
		lnum++
		line := sc.Text()

		matches := reType1.FindStringSubmatch(line)
		if len(matches) == 0 {
			return emptyReport, errors.Errorf("failed to find matches for regex [line num: %d], [line: %s]",
				lnum, line)
		}

		if matches[colScore] != "" {
			sc, err := strconv.ParseFloat(matches[colScore], 64)
			if err != nil {
				return emptyReport, errors.Wrapf(err,
					"failed to parse score [line num: %d], [line: %s], score[%s]",
					lnum, line, matches[colScore])
			}

			sc = sanitizeScore(sc)
			score = score + sc
			r.SpamAssassin.Headers = append(r.SpamAssassin.Headers, Headers{
				Score:       sc,
				Tag:         matches[colTag],
				Description: matches[colDescr],
			})

		} else {
			last := len(r.SpamAssassin.Headers) - 1
			if last >= 0 {
				r.SpamAssassin.Headers[last].Description += " " + matches[colDescr]
			}
		}

	}

	r.SpamAssassin.Score = sanitizeScore(score)

	return r, nil
}

func sanitizeScore(sc float64) float64 {
	if sc == -0 {
		sc = 0
	}
	return math.RoundToEven(sc*100) / 100
}

func processReportType2(data io.Reader) (Report, error) {

	const (
		colFullMatch = iota
		colScore
		colTag
		colDescr
	)
	var (
		r     Report
		score float64
		lnum  int
		start bool
	)
	sc := bufio.NewScanner(data)
	for sc.Scan() {
		lnum++
		line := sc.Text()
		if !start {
			if strings.Contains(line, "----") {
				start = true
			}
			continue
		}

		matches := reType2.FindStringSubmatch(line)
		if len(matches) != 0 {
			sc, err := strconv.ParseFloat(matches[colScore], 64)
			if err != nil {
				return emptyReport, errors.Wrapf(err,
					"failed to parse score [line num: %d], [line: %s], score[%s]",
					lnum, line, matches[colScore])
			}

			sc = sanitizeScore(sc)
			score = score + sc
			r.SpamAssassin.Headers = append(r.SpamAssassin.Headers, Headers{
				Score:       sc,
				Tag:         matches[colTag],
				Description: matches[colDescr],
			})

		} else {
			last := len(r.SpamAssassin.Headers) - 1
			if last >= 0 {
				line = strings.TrimSpace(line)
				r.SpamAssassin.Headers[last].Description += " " + line
			}
		}
	}

	r.SpamAssassin.Score = sanitizeScore(score)

	return r, nil
}
