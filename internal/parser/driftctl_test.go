package parser_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/example/driftctl-export/internal/parser"
)

func writeTempReport(t *testing.T, report parser.DriftReport) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "driftreport-*.json")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if err := json.NewEncoder(f).Encode(report); err != nil {
		t.Fatalf("encoding report: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseFile_ValidReport(t *testing.T) {
	input := parser.DriftReport{}
	input.Summary.TotalResources = 3
	input.Summary.DriftedCount = 1
	input.Drifted = []parser.DriftResource{
		{
			ResourceID:   "aws_s3_bucket.my-bucket",
			ResourceType: "aws_s3_bucket",
			Source:       "aws",
			Differences: []parser.Difference{
				{Attribute: "acl", Wanted: "private", Found: "public-read"},
			},
		},
	}

	path := writeTempReport(t, input)
	report, err := parser.ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Summary.TotalResources != 3 {
		t.Errorf("expected TotalResources=3, got %d", report.Summary.TotalResources)
	}
	if len(report.Drifted) != 1 {
		t.Fatalf("expected 1 drifted resource, got %d", len(report.Drifted))
	}
	if report.Drifted[0].ResourceID != "aws_s3_bucket.my-bucket" {
		t.Errorf("unexpected resource ID: %s", report.Drifted[0].ResourceID)
	}
}

func TestParseFile_MissingFile(t *testing.T) {
	_, err := parser.ParseFile("/nonexistent/path/report.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestParseFile_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "bad-*.json")
	f.WriteString("not valid json{{{")
	f.Close()

	_, err := parser.ParseFile(f.Name())
	if err == nil {
		t.Fatal("expected JSON decode error, got nil")
	}
}
