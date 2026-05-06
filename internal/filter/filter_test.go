package filter_test

import (
	"testing"

	"github.com/snyk/driftctl/pkg/resource"

	"github.com/example/driftctl-export/internal/filter"
)

// stubResource satisfies resource.Resource for testing purposes.
type stubResource struct {
	id      string
	resType string
	diff    map[string]interface{}
}

func (s stubResource) ResourceId() string                  { return s.id }
func (s stubResource) ResourceType() string                { return s.resType }
func (s stubResource) Diff() map[string]interface{}        { return s.diff }
func (s stubResource) Attributes() map[string]interface{}  { return nil }

func buildReport() *resource.ScanReport {
	return &resource.ScanReport{
		Managed: []resource.Resource{
			stubResource{id: "bucket-1", resType: "aws_s3_bucket", diff: map[string]interface{}{"acl": "public"}},
			stubResource{id: "vpc-1", resType: "aws_vpc", diff: nil},
		},
		Unmanaged: []resource.Resource{
			stubResource{id: "sg-1", resType: "aws_security_group"},
		},
		Deleted: []resource.Resource{
			stubResource{id: "igw-1", resType: "aws_internet_gateway"},
		},
	}
}

func TestApply_NoFilter(t *testing.T) {
	report := buildReport()
	results := filter.Apply(report, filter.Options{})
	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}
}

func TestApply_OnlyManaged(t *testing.T) {
	report := buildReport()
	results := filter.Apply(report, filter.Options{OnlyManaged: true})
	if len(results) != 2 {
		t.Fatalf("expected 2 managed results, got %d", len(results))
	}
	for _, r := range results {
		if r.Source != "managed" {
			t.Errorf("expected source=managed, got %s", r.Source)
		}
	}
}

func TestApply_OnlyDrifted(t *testing.T) {
	report := buildReport()
	results := filter.Apply(report, filter.Options{OnlyDrifted: true})
	// Only bucket-1 has a diff; unmanaged/deleted are excluded because changed==nil.
	if len(results) != 1 {
		t.Fatalf("expected 1 drifted result, got %d", len(results))
	}
	if results[0].ID != "bucket-1" {
		t.Errorf("unexpected resource id: %s", results[0].ID)
	}
}

func TestApply_FilterByType(t *testing.T) {
	report := buildReport()
	results := filter.Apply(report, filter.Options{
		ResourceTypes: []string{"aws_s3_bucket"},
	})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Type != "aws_s3_bucket" {
		t.Errorf("unexpected type: %s", results[0].Type)
	}
}

func TestApply_ChangedAttrs(t *testing.T) {
	report := buildReport()
	results := filter.Apply(report, filter.Options{OnlyManaged: true})
	var bucket filter.ResourceSummary
	for _, r := range results {
		if r.ID == "bucket-1" {
			bucket = r
		}
	}
	if len(bucket.ChangedAttrs) != 1 || bucket.ChangedAttrs[0] != "acl" {
		t.Errorf("expected changed attr 'acl', got %v", bucket.ChangedAttrs)
	}
}
