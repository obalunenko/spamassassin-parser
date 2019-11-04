package parser

import (
	"bufio"
	"io"
	"regexp"
	"strconv"

	"github.com/pkg/errors"

	"github.com/oleg-balunenko/spamassassin-parser/pkg/models"
)

var (
	reType1 = regexp.MustCompile(`([*])[\s]+([-]?\d.\d)?[\s](([[:word:]]+)?[\s](.*))`)
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
			sc, err := strconv.ParseFloat(matches[colScore], 64)
			if err != nil {
				return emptyReport, errors.Wrapf(err,
					"failed to parse score [line num: %d], [line: %s], score[%s]",
					lnum, line, matches[colScore])
			}

			sc = sanitizeScore(sc)
			score = score + sc
			r.SpamAssassin.Headers = append(r.SpamAssassin.Headers, models.Headers{
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
