package enricher_test

import (
	"testing"

	"github.com/snyk/driftctl-export/internal/enricher"
	"github.com/snyk/driftctl-export/internal/parser"
)

func buildReport() *parser.Report {
	return &parser.Report{
		UnmanagedResources: []parser.Resource{
			{ID: "bucket-1", Type: "aws_s3_bucket"},
			{ID: "sg-1", Type: "aws_security_group"},
		},
		ManagedResources: []parser.Resource{
			{ID: "vpc-1", Type: "aws_vpc"},
		},
	}
}

func TestApply_NilReport(t *testing.T) {
	result := enricher.Apply(nil, enricher.Options{})
	if result != nil {
		t.Fatal("expected nil for nil report input")
	}
}

func TestApply_SetsEnvironment(t *testing.T) {
	report := buildReport()
	enricher.Apply(report, enricher.Options{Environment: "production"})

	for _, r := range report.UnmanagedResources {
		if r.Environment != "production" {
			t.Errorf("unmanaged resource %s: expected environment 'production', got %q", r.ID, r.Environment)
		}
	}
	for _, r := range report.ManagedResources {
		if r.Environment != "production" {
			t.Errorf("managed resource %s: expected environment 'production', got %q", r.ID, r.Environment)
		}
	}
}

func TestApply_DefaultSeverity(t *testing.T) {
	report := buildReport()
	enricher.Apply(report, enricher.Options{DefaultSeverity: "low"})

	for _, r := range report.UnmanagedResources {
		if r.Severity != "low" {
			t.Errorf("resource %s: expected severity 'low', got %q", r.ID, r.Severity)
		}
	}
}

func TestApply_SeverityMapOverride(t *testing.T) {
	report := buildReport()
	enricher.Apply(report, enricher.Options{
		DefaultSeverity: "low",
		SeverityMap: map[string]string{
			"aws_security_group": "high",
		},
	})

	for _, r := range report.UnmanagedResources {
		switch r.Type {
		case "aws_security_group":
			if r.Severity != "high" {
				t.Errorf("expected 'high' for aws_security_group, got %q", r.Severity)
			}
		default:
			if r.Severity != "low" {
				t.Errorf("expected 'low' for %s, got %q", r.Type, r.Severity)
			}
		}
	}
}

func TestApply_FallbackDefaultSeverity(t *testing.T) {
	report := buildReport()
	// No DefaultSeverity set — should fall back to "medium".
	enricher.Apply(report, enricher.Options{})

	for _, r := range report.UnmanagedResources {
		if r.Severity != "medium" {
			t.Errorf("resource %s: expected fallback severity 'medium', got %q", r.ID, r.Severity)
		}
	}
}
