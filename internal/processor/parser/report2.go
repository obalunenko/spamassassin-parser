// Package parser provides functionality to parse spamassassin result into json report.
package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/obalunenko/spamassassin-parser/internal/processor/models"
)

var reType2 = regexp.MustCompile(`(?m)(-?\d+.\d+)\s+(\w+)\s+(.*\n?)`)

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
				return emptyReport, fmt.Errorf(
					"failed to make header [line num: %d], [line: %s]: %w", lnum, line, err)
			}

			score += h.Score
			r.SpamAssassin.Headers = append(r.SpamAssassin.Headers, h)
		} else {
			indexShift := 1

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
