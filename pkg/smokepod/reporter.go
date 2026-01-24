package smokepod

import (
	"encoding/json"
	"io"
)

// Reporter outputs test results.
type Reporter struct {
	writer io.Writer
	pretty bool
}

// NewReporter creates a new reporter that writes to w.
func NewReporter(w io.Writer) *Reporter {
	return &Reporter{writer: w}
}

// SetPretty enables or disables pretty-printed JSON output.
func (r *Reporter) SetPretty(p bool) {
	r.pretty = p
}

// Report writes the result as JSON to the writer.
func (r *Reporter) Report(result *Result) error {
	var data []byte
	var err error

	if r.pretty {
		data, err = json.MarshalIndent(result, "", "  ")
	} else {
		data, err = json.Marshal(result)
	}
	if err != nil {
		return err
	}

	// Add trailing newline for pretty output
	if r.pretty {
		data = append(data, '\n')
	}

	_, err = r.writer.Write(data)
	return err
}
