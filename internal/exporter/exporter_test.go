package exporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/snyk/driftctl-export/internal/exporter"
	"github.com/snyk/driftctl-export/internal/parser"
)

func sampleReport() *parser.DriftReport {
	return &parser.DriftReport{
		Summary: parser.Summary{
			TotalResources:   3,
			DriftedResources: 2,
		},
		Differences: []parser.Difference{
			{Res: parser.Resource{Type: "aws_s3_bucket", ID: "my-bucket"}, Status: "changed"},
			{Res: parser.Resource{Type: "aws_iam_role", ID: "my-role"}, Status: "missing"},
		},
	}
}

func TestExportJSON(t *testing.T) {
	var buf bytes.Buffer
	e := exporter.New(exporter.FormatJSON, &buf)

	if err := e.Export(sampleReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got parser.DriftReport
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if got.Summary.TotalResources != 3 {
		t.Errorf("expected TotalResources=3, got %d", got.Summary.TotalResources)
	}
	if len(got.Differences) != 2 {
		t.Errorf("expected 2 differences, got %d", len(got.Differences))
	}
}

func TestExportCSV(t *testing.T) {
	var buf bytes.Buffer
	e := exporter.New(exporter.FormatCSV, &buf)

	if err := e.Export(sampleReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 { // header + 2 data rows
		t.Fatalf("expected 3 lines, got %d: %v", len(lines), lines)
	}
	if lines[0] != "type,id,status" {
		t.Errorf("unexpected header: %s", lines[0])
	}
	if !strings.Contains(lines[1], "aws_s3_bucket") {
		t.Errorf("expected aws_s3_bucket in row 1, got: %s", lines[1])
	}
}

func TestExportUnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	e := exporter.New("xml", &buf)
	if err := e.Export(sampleReport()); err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}
