package summary_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/driftctl-export/internal/parser"
	"github.com/your-org/driftctl-export/internal/summary"
)

func buildReport() *parser.Report {
	return &parser.Report{
		ManagedResources: []parser.Resource{
			{ID: "res-1", Type: "aws_s3_bucket", Differences: []parser.Difference{
				{Attribute: "tags", Before: "a", After: "b"},
			}},
			{ID: "res-2", Type: "aws_instance", Differences: nil},
		},
		UnmanagedResources: []parser.Resource{
			{ID: "res-3", Type: "aws_vpc"},
		},
		DeletedResources: []parser.Resource{
			{ID: "res-4", Type: "aws_subnet"},
			{ID: "res-5", Type: "aws_subnet"},
		},
	}
}

func TestCompute(t *testing.T) {
	report := buildReport()
	s := summary.Compute(report)

	if s.TotalResources != 5 {
		t.Errorf("expected TotalResources=5, got %d", s.TotalResources)
	}
	if s.ManagedCount != 2 {
		t.Errorf("expected ManagedCount=2, got %d", s.ManagedCount)
	}
	if s.UnmanagedCount != 1 {
		t.Errorf("expected UnmanagedCount=1, got %d", s.UnmanagedCount)
	}
	if s.DeletedCount != 2 {
		t.Errorf("expected DeletedCount=2, got %d", s.DeletedCount)
	}
	if s.DriftedCount != 1 {
		t.Errorf("expected DriftedCount=1, got %d", s.DriftedCount)
	}
	if s.CoveragePercent != 40.0 {
		t.Errorf("expected CoveragePercent=40.0, got %.1f", s.CoveragePercent)
	}
}

func TestComputeEmptyReport(t *testing.T) {
	s := summary.Compute(&parser.Report{})
	if s.TotalResources != 0 || s.CoveragePercent != 0 {
		t.Error("expected zero stats for empty report")
	}
}

func TestWrite(t *testing.T) {
	report := buildReport()
	s := summary.Compute(report)

	var buf bytes.Buffer
	if err := summary.Write(&buf, s); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	output := buf.String()
	for _, expected := range []string{"Total Resources", "Managed", "Unmanaged", "Deleted", "Drifted", "Coverage", "40.0%"} {
		if !strings.Contains(output, expected) {
			t.Errorf("expected output to contain %q, got:\n%s", expected, output)
		}
	}
}
