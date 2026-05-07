package pipeline_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/example/driftctl-export/internal/parser"
	"github.com/example/driftctl-export/internal/pipeline"
)

func writeTempReport(t *testing.T, dir string, name string) string {
	t.Helper()
	report := parser.DriftReport{
		Managed: []parser.Resource{
			{ID: "res-1", Type: "aws_s3_bucket", Source: "terraform"},
		},
		Unmanaged: []parser.Resource{
			{ID: "res-2", Type: "aws_instance", Source: "cloud"},
		},
		Deleted: []parser.Resource{},
	}
	data, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("marshal report: %v", err)
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write temp report: %v", err)
	}
	return path
}

func TestRun_Success(t *testing.T) {
	dir := t.TempDir()
	writeTempReport(t, dir, "report.json")
	outPath := filepath.Join(dir, "out.json")

	result, err := pipeline.Run(pipeline.Options{
		InputGlob:  filepath.Join(dir, "*.json"),
		OutputPath: outPath,
		Format:     "json",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.FilesProcessed != 1 {
		t.Errorf("expected 1 file processed, got %d", result.FilesProcessed)
	}
	if result.ResourcesTotal != 2 {
		t.Errorf("expected 2 resources, got %d", result.ResourcesTotal)
	}
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Error("expected output file to exist")
	}
}

func TestRun_NoMatchingFiles(t *testing.T) {
	dir := t.TempDir()
	_, err := pipeline.Run(pipeline.Options{
		InputGlob:  filepath.Join(dir, "*.json"),
		OutputPath: filepath.Join(dir, "out.json"),
		Format:     "json",
	})
	if err == nil {
		t.Fatal("expected error for no matching files")
	}
}

func TestRun_InvalidFormat(t *testing.T) {
	dir := t.TempDir()
	writeTempReport(t, dir, "report.json")

	_, err := pipeline.Run(pipeline.Options{
		InputGlob:  filepath.Join(dir, "*.json"),
		OutputPath: filepath.Join(dir, "out.xml"),
		Format:     "xml",
	})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}
