// Package parser provides functionality to parse spamassassin result into json report.
package parser

import (
	"fmt"
	"strconv"

	"github.com/obalunenko/spamassassin-parser/internal/processor/models"
)

func makeHeader(score, tag, description string) (models.Headers, error) {
	const bitsize = 64

	sc, err := strconv.ParseFloat(score, bitsize)
	if err != nil {
		return models.Headers{}, fmt.Errorf("failed to parse score: %w", err)
	}

	return models.Headers{
		Score:       sanitizeScore(sc),
		Tag:         tag,
		Description: description,
	}, nil
}
