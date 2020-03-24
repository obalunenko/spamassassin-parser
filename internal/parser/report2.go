// Package parser provides functionality to parse spamassassin result into json report.
package parser

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/oleg-balunenko/spamassassin-parser/internal/models"
)

var (
	reType2 = regexp.MustCompile(`(?m)([-]?\d+.\d+)[\s]+([[:word:]]+)\s+(.*[\n]?)`)
)

type report2Parser struct{}

func (rp report2Parser) Parse(data io.Reader) (models.Report, error) {
	const (
		colFullMatch = iota
		colScore
		colTag
		colDescr
	)

	var (
		r     models.Report
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
			h, err := makeHeader(matches[colScore], matches[colTag], matches[colDescr])
			if err != nil {
				return emptyReport, errors.Wrapf(err,
					"failed to make header [line num: %d], [line: %s]", lnum, line)
			}

			score += h.Score
			r.SpamAssassin.Headers = append(r.SpamAssassin.Headers, h)
		} else {
			var indexShift = 1

			last := len(r.SpamAssassin.Headers) - indexShift
			if last >= 0 {
				line = strings.TrimSpace(line)
				r.SpamAssassin.Headers[last].Description += " " + line
			}
		}
	}

	r.SpamAssassin.Score = sanitizeScore(score)

	return r, nil
}
