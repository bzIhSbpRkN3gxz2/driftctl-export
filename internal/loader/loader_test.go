package loader_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/example/driftctl-export/internal/loader"
)

func writeTempReport(t *testing.T, dir, name string, data interface{}) string {
	t.Helper()
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("failed to marshal report: %v", err)
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatalf("failed to write temp report: %v", err)
	}
	return path
}

func validReportData() map[string]interface{} {
	return map[string]interface{}{
		"summary": map[string]interface{}{
			"total_resources": 3,
			"total_changed":   1,
			"total_unmanaged": 1,
			"total_missing":   1,
			"total_managed":   1,
		},
		"managed":   []interface{}{},
		"unmanaged": []interface{}{},
		"missing":   []interface{}{},
		"differences": []interface{}{},
	}
}

func TestLoadFile_Valid(t *testing.T) {
	dir := t.TempDir()
	path := writeTempReport(t, dir, "report.json", validReportData())

	report, err := loader.LoadFile(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if report == nil {
		t.Fatal("expected non-nil report")
	}
}

func TestLoadFile_Missing(t *testing.T) {
	_, err := loader.LoadFile("/nonexistent/path/report.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadGlob_MultipleFiles(t *testing.T) {
	dir := t.TempDir()
	writeTempReport(t, dir, "report1.json", validReportData())
	writeTempReport(t, dir, "report2.json", validReportData())

	result, err := loader.LoadGlob(filepath.Join(dir, "*.json"))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(result.Reports) != 2 {
		t.Errorf("expected 2 reports, got %d", len(result.Reports))
	}
	if len(result.Errors) != 0 {
		t.Errorf("expected no errors, got %d", len(result.Errors))
	}
}

func TestLoadGlob_NoMatches(t *testing.T) {
	dir := t.TempDir()
	_, err := loader.LoadGlob(filepath.Join(dir, "*.json"))
	if err == nil {
		t.Fatal("expected error for no matches, got nil")
	}
}

func TestLoadGlob_PartialErrors(t *testing.T) {
	dir := t.TempDir()
	writeTempReport(t, dir, "valid.json", validReportData())

	badPath := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(badPath, []byte("not-json"), 0o644); err != nil {
		t.Fatalf("failed to write bad file: %v", err)
	}

	result, err := loader.LoadGlob(filepath.Join(dir, "*.json"))
	if err != nil {
		t.Fatalf("expected no top-level error, got: %v", err)
	}
	if len(result.Reports) != 1 {
		t.Errorf("expected 1 valid report, got %d", len(result.Reports))
	}
	if len(result.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(result.Errors))
	}
}
