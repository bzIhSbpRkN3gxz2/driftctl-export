package reporter

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/snyk/driftctl-export/internal/summary"
)

// Format represents the output format for the drift report.
type Format string

const (
	FormatTable Format = "table"
	FormatCompact Format = "compact"
)

// Reporter renders a drift summary to an output stream.
type Reporter struct {
	format Format
	writer io.Writer
}

// New creates a new Reporter with the given format and writer.
func New(format Format, w io.Writer) *Reporter {
	return &Reporter{format: format, writer: w}
}

// Render writes the summary to the configured writer in the configured format.
func (r *Reporter) Render(s summary.Summary) error {
	switch r.format {
	case FormatTable:
		return r.renderTable(s)
	case FormatCompact:
		return r.renderCompact(s)
	default:
		return fmt.Errorf("unsupported reporter format: %s", r.format)
	}
}

func (r *Reporter) renderTable(s summary.Summary) error {
	w := tabwriter.NewWriter(r.writer, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "METRIC\tCOUNT")
	fmt.Fprintln(w, "------\t-----")
	fmt.Fprintf(w, "Managed\t%d\n", s.Managed)
	fmt.Fprintf(w, "Unmanaged\t%d\n", s.Unmanaged)
	fmt.Fprintf(w, "Missing\t%d\n", s.Missing)
	fmt.Fprintf(w, "Differences\t%d\n", s.Differences)
	fmt.Fprintf(w, "Total\t%d\n", s.Total)
	return w.Flush()
}

func (r *Reporter) renderCompact(s summary.Summary) error {
	_, err := fmt.Fprintf(
		r.writer,
		"managed=%d unmanaged=%d missing=%d differences=%d total=%d\n",
		s.Managed, s.Unmanaged, s.Missing, s.Differences, s.Total,
	)
	return err
}
