// Package parser provides functionality to parse spamassassin result into json report.
package parser

import (
	"fmt"
	"strconv"

	"github.com/oleg-balunenko/spamassassin-parser/internal/processor/models"
)

func makeHeader(score, tag, description string) (models.Headers, error) {
	sc, err := strconv.ParseFloat(score, 64)
	if err != nil {
		return models.Headers{}, fmt.Errorf("failed to parse score: %w", err)
	}

	return models.Headers{
		Score:       sanitizeScore(sc),
		Tag:         tag,
		Description: description,
	}, nil
}
