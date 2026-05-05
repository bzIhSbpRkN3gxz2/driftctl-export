package exporter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/snyk/driftctl-export/internal/parser"
)

// Format represents the output format for the export.
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// Exporter writes drift report data to the given writer in the requested format.
type Exporter struct {
	format Format
	writer io.Writer
}

// New creates a new Exporter targeting the given format and writer.
func New(format Format, w io.Writer) *Exporter {
	return &Exporter{format: format, writer: w}
}

// NewFileExporter creates an Exporter that writes to the specified file path.
func NewFileExporter(format Format, path string) (*Exporter, *os.File, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, nil, fmt.Errorf("creating output file: %w", err)
	}
	return New(format, f), f, nil
}

// Export serialises the report according to the configured format.
func (e *Exporter) Export(report *parser.DriftReport) error {
	switch e.format {
	case FormatJSON:
		return e.exportJSON(report)
	case FormatCSV:
		return e.exportCSV(report)
	default:
		return fmt.Errorf("unsupported format: %s", e.format)
	}
}

func (e *Exporter) exportJSON(report *parser.DriftReport) error {
	enc := json.NewEncoder(e.writer)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}
	return nil
}

func (e *Exporter) exportCSV(report *parser.DriftReport) error {
	w := csv.NewWriter(e.writer)
	defer w.Flush()

	header := []string{"type", "id", "status"}
	if err := w.Write(header); err != nil {
		return fmt.Errorf("writing CSV header: %w", err)
	}

	for _, r := range report.Differences {
		row := []string{r.Res.Type, r.Res.ID, string(r.Status)}
		if err := w.Write(row); err != nil {
			return fmt.Errorf("writing CSV row: %w", err)
		}
	}

	if err := w.Error(); err != nil {
		return fmt.Errorf("flushing CSV writer: %w", err)
	}
	return nil
}
