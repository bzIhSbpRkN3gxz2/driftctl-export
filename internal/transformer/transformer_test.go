package transformer_test

import (
	"testing"

	"github.com/snyk/driftctl-export/internal/models"
	"github.com/snyk/driftctl-export/internal/transformer"
)

func buildReport() *models.DriftReport {
	return &models.DriftReport{
		Unmanaged: []models.Resource{
			{ResourceID: "  AWS::S3::Bucket/my-bucket  ", ResourceType: "aws_s3_bucket"},
		},
		Managed: []models.Resource{
			{ResourceID: "PREFIX/res-1", ResourceType: "aws_iam_role"},
		},
		Drifted: []models.Resource{
			{ResourceID: "PREFIX/res-2", ResourceType: "aws_lambda_function"},
		},
	}
}

func TestApply_NilReport(t *testing.T) {
	result := transformer.Apply(nil, transformer.Option{})
	if result != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestApply_NormalizeIDs(t *testing.T) {
	report := buildReport()
	transformer.Apply(report, transformer.Option{NormalizeIDs: true})

	want := "aws::s3::bucket/my-bucket"
	got := report.Unmanaged[0].ResourceID
	if got != want {
		t.Errorf("NormalizeIDs: got %q, want %q", got, want)
	}
}

func TestApply_UppercaseTypes(t *testing.T) {
	report := buildReport()
	transformer.Apply(report, transformer.Option{UppercaseTypes: true})

	want := "AWS_S3_BUCKET"
	got := report.Unmanaged[0].ResourceType
	if got != want {
		t.Errorf("UppercaseTypes: got %q, want %q", got, want)
	}
}

func TestApply_StripPrefix(t *testing.T) {
	report := buildReport()
	transformer.Apply(report, transformer.Option{StripPrefix: "PREFIX/"})

	wantManaged := "res-1"
	if got := report.Managed[0].ResourceID; got != wantManaged {
		t.Errorf("StripPrefix managed: got %q, want %q", got, wantManaged)
	}

	wantDrifted := "res-2"
	if got := report.Drifted[0].ResourceID; got != wantDrifted {
		t.Errorf("StripPrefix drifted: got %q, want %q", got, wantDrifted)
	}
}

func TestApply_CombinedOptions(t *testing.T) {
	report := buildReport()
	transformer.Apply(report, transformer.Option{
		NormalizeIDs:   true,
		UppercaseTypes: true,
		StripPrefix:    "PREFIX/",
	})

	if got := report.Managed[0].ResourceID; got != "res-1" {
		t.Errorf("combined managed id: got %q", got)
	}
	if got := report.Managed[0].ResourceType; got != "AWS_IAM_ROLE" {
		t.Errorf("combined managed type: got %q", got)
	}
}

func TestApply_EmptyReport(t *testing.T) {
	report := &models.DriftReport{}
	result := transformer.Apply(report, transformer.Option{NormalizeIDs: true, UppercaseTypes: true})
	if result == nil {
		t.Fatal("expected non-nil result for empty report")
	}
}
