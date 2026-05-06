package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/snyk/driftctl-export/internal/reporter"
	"github.com/snyk/driftctl-export/internal/summary"
)

func buildSummary() summary.Summary {
	return summary.Summary{
		Managed:     10,
		Unmanaged:   3,
		Missing:     2,
		Differences: 5,
		Total:       15,
	}
}

func TestRenderTable(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(reporter.FormatTable, &buf)

	if err := r.Render(buildSummary()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"METRIC", "COUNT", "Managed", "10", "Unmanaged", "3", "Missing", "2", "Differences", "5", "Total", "15"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderCompact(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(reporter.FormatCompact, &buf)

	if err := r.Render(buildSummary()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "managed=10 unmanaged=3 missing=2 differences=5 total=15"
	if !strings.Contains(buf.String(), want) {
		t.Errorf("expected %q, got %q", want, buf.String())
	}
}

func TestRenderUnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New("xml", &buf)

	if err := r.Render(buildSummary()); err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}

func TestRenderEmptySummary(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(reporter.FormatCompact, &buf)

	if err := r.Render(summary.Summary{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "managed=0 unmanaged=0 missing=0 differences=0 total=0"
	if !strings.Contains(buf.String(), want) {
		t.Errorf("expected %q, got %q", want, buf.String())
	}
}
