// Package parser provides functionality to parse spamassassin result into json report.
package parser

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/oleg-balunenko/spamassassin-parser/internal/models"
)

func makeHeader(score, tag, description string) (models.Headers, error) {
	sc, err := strconv.ParseFloat(score, 64)
	if err != nil {
		return models.Headers{}, errors.Wrapf(err,
			"failed to parse score")
	}

	return models.Headers{
		Score:       sanitizeScore(sc),
		Tag:         tag,
		Description: description,
	}, nil
}
