package utils

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
)

// PrettyPrint appends to passed struct indents and returns a human readable form of struct.
func PrettyPrint(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal")
	}

	var out bytes.Buffer
	if err := json.Indent(&out, b, "", "\t"); err != nil {
		return "", errors.Wrap(err, "failed to indent")
	}
	if _, err := out.WriteString("\n"); err != nil {
		return "", errors.Wrap(err, "failed to write string")
	}

	return out.String(), nil
}
