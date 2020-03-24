// Package parser provides functionality to parse spamassassin result into json report.
package parser

import (
	"bufio"
	"io"
	"regexp"

	"github.com/pkg/errors"

	"github.com/oleg-balunenko/spamassassin-parser/internal/models"
)

var (
	reType1 = regexp.MustCompile(`([*])[\s]+([-]?\d+.\d+)?[\s](([[:word:]]+)?[\s](.*))`)
)

type report1Parser struct{}

func (rp report1Parser) Parse(data io.Reader) (models.Report, error) {
	const (
		colFullMatch = iota
		colAsterisk
		colScore
		colTagWithDescr
		colTag
		colDescr
	)

	var (
		r     models.Report
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
				r.SpamAssassin.Headers[last].Description += " " + matches[colDescr]
			}
		}
	}

	r.SpamAssassin.Score = sanitizeScore(score)

	return r, nil
}
