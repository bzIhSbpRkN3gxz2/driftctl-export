package validator_test

import (
	"testing"

	"github.com/snyk/driftctl-export/internal/parser"
	"github.com/snyk/driftctl-export/internal/validator"
)

func buildValidReport() *parser.DriftReport {
	return &parser.DriftReport{
		Differences: []parser.Difference{
			{
				Res: parser.Resource{
					ResourceType: "aws_s3_bucket",
					ResourceID:   "my-bucket",
				},
			},
		},
		UnmanagedResources: []parser.Resource{
			{
				ResourceType: "aws_iam_role",
				ResourceID:   "my-role",
			},
		},
	}
}

func TestValidate_ValidReport(t *testing.T) {
	report := buildValidReport()
	if err := validator.Validate(report); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_NilReport(t *testing.T) {
	err := validator.Validate(nil)
	if err == nil {
		t.Fatal("expected error for nil report, got nil")
	}
}

func TestValidate_MissingResourceType(t *testing.T) {
	report := buildValidReport()
	report.Differences[0].Res.ResourceType = ""

	err := validator.Validate(report)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(ve.Issues))
	}
}

func TestValidate_MissingUnmanagedResourceID(t *testing.T) {
	report := buildValidReport()
	report.UnmanagedResources[0].ResourceID = ""

	err := validator.Validate(report)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Issues) == 0 {
		t.Fatal("expected at least one issue")
	}
}

func TestValidate_EmptyReport(t *testing.T) {
	report := &parser.DriftReport{}
	if err := validator.Validate(report); err != nil {
		t.Fatalf("expected no error for empty report, got: %v", err)
	}
}
