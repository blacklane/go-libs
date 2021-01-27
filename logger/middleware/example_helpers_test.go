package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

// prettyWriter takes a JSON string as []byte an ident it before writing to os.Stdout..
var prettyWriter = prettyJSONWriter{}

type prettyJSONWriter struct{}

func (pw prettyJSONWriter) Write(p []byte) (int, error) {
	buf := &bytes.Buffer{}

	data := &map[string]interface{}{}
	if err := json.Unmarshal(p, data); err != nil {
		return 0, fmt.Errorf("prettyJSONWriter: could not json.Unmarshal data: %w", err)
	}

	pp, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return 0, fmt.Errorf("prettyJSONWriter: could not json.MarshalIndent data: %w", err)
	}

	buf.Write(pp)

	n, err := buf.WriteTo(os.Stdout)
	return int(n), err
}
