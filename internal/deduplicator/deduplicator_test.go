package deduplicator_test

import (
	"testing"

	"github.com/snyk/driftctl-export/internal/deduplicator"
	"github.com/snyk/driftctl-export/internal/parser"
)

func buildReport() *parser.Report {
	return &parser.Report{
		UnmanagedResources: []parser.Resource{
			{Type: "aws_s3_bucket", ResourceID: "bucket-1"},
			{Type: "aws_s3_bucket", ResourceID: "bucket-1"}, // duplicate
			{Type: "aws_s3_bucket", ResourceID: "bucket-2"},
		},
		ManagedResources: []parser.ManagedResource{
			{Type: "aws_iam_role", ResourceID: "role-1"},
			{Type: "aws_iam_role", ResourceID: "role-1"}, // duplicate
			{Type: "aws_iam_role", ResourceID: "role-2"},
		},
		DeletedResources: []parser.Resource{
			{Type: "aws_instance", ResourceID: "i-abc"},
			{Type: "aws_instance", ResourceID: "i-abc"}, // duplicate
		},
	}
}

func TestApply_RemovesDuplicates(t *testing.T) {
	report := buildReport()
	result := deduplicator.Apply(report)

	if len(result.UnmanagedResources) != 2 {
		t.Errorf("expected 2 unmanaged resources, got %d", len(result.UnmanagedResources))
	}
	if len(result.ManagedResources) != 2 {
		t.Errorf("expected 2 managed resources, got %d", len(result.ManagedResources))
	}
	if len(result.DeletedResources) != 1 {
		t.Errorf("expected 1 deleted resource, got %d", len(result.DeletedResources))
	}
}

func TestApply_NilReport(t *testing.T) {
	result := deduplicator.Apply(nil)
	if result != nil {
		t.Error("expected nil result for nil input")
	}
}

func TestApply_NoDuplicates(t *testing.T) {
	report := &parser.Report{
		UnmanagedResources: []parser.Resource{
			{Type: "aws_s3_bucket", ResourceID: "bucket-1"},
			{Type: "aws_s3_bucket", ResourceID: "bucket-2"},
		},
	}
	result := deduplicator.Apply(report)
	if len(result.UnmanagedResources) != 2 {
		t.Errorf("expected 2 resources, got %d", len(result.UnmanagedResources))
	}
}

func TestApply_EmptyReport(t *testing.T) {
	report := &parser.Report{}
	result := deduplicator.Apply(report)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result.UnmanagedResources) != 0 {
		t.Errorf("expected 0 unmanaged resources, got %d", len(result.UnmanagedResources))
	}
}

func TestApply_PreservesOrder(t *testing.T) {
	report := &parser.Report{
		UnmanagedResources: []parser.Resource{
			{Type: "aws_s3_bucket", ResourceID: "bucket-b"},
			{Type: "aws_s3_bucket", ResourceID: "bucket-a"},
			{Type: "aws_s3_bucket", ResourceID: "bucket-b"}, // duplicate
		},
	}
	result := deduplicator.Apply(report)
	if result.UnmanagedResources[0].ResourceID != "bucket-b" {
		t.Errorf("expected first resource to be bucket-b, got %s", result.UnmanagedResources[0].ResourceID)
	}
	if result.UnmanagedResources[1].ResourceID != "bucket-a" {
		t.Errorf("expected second resource to be bucket-a, got %s", result.UnmanagedResources[1].ResourceID)
	}
}
