package tracing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/blacklane/go-libs/logger"
)

var noopLogger = *logger.FromContext(context.Background())

// prettyJSONWriter takes a JSON string as []byte an ident it before writing to os.Stdout..
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
	buf.WriteByte('\n')

	n, err := buf.WriteTo(os.Stdout)
	return int(n), err
}
