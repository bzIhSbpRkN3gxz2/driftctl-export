package sorter_test

import (
	"testing"

	"github.com/snyk/driftctl/pkg/analyser"
	"github.com/snyk/driftctl/pkg/resource"

	"github.com/owner/driftctl-export/internal/sorter"
)

func buildReport() *analyser.Analysis {
	a := &analyser.Analysis{}
	a.AddManaged(
		&resource.Resource{Type: "aws_s3_bucket", Id: "bucket-b"},
		&resource.Resource{Type: "aws_instance", Id: "instance-a"},
	)
	a.AddUnmanaged(
		&resource.Resource{Type: "aws_vpc", Id: "vpc-z"},
		&resource.Resource{Type: "aws_subnet", Id: "subnet-m"},
	)
	a.AddDeleted(
		&resource.Resource{Type: "aws_iam_role", Id: "role-c"},
		&resource.Resource{Type: "aws_iam_policy", Id: "policy-a"},
	)
	return a
}

func TestApply_SortByTypeAsc(t *testing.T) {
	report := buildReport()
	err := sorter.Apply(report, sorter.Options{Field: sorter.ByType, Order: sorter.Asc})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	managed := report.Managed()
	if managed[0].Type > managed[1].Type {
		t.Errorf("expected ascending order by type, got %s before %s", managed[0].Type, managed[1].Type)
	}
}

func TestApply_SortByTypeDesc(t *testing.T) {
	report := buildReport()
	err := sorter.Apply(report, sorter.Options{Field: sorter.ByType, Order: sorter.Desc})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	managed := report.Managed()
	if managed[0].Type < managed[1].Type {
		t.Errorf("expected descending order by type, got %s before %s", managed[0].Type, managed[1].Type)
	}
}

func TestApply_SortByIDAsc(t *testing.T) {
	report := buildReport()
	err := sorter.Apply(report, sorter.Options{Field: sorter.ByID, Order: sorter.Asc})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	unmanaged := report.Unmanaged()
	if unmanaged[0].Id > unmanaged[1].Id {
		t.Errorf("expected ascending order by id, got %s before %s", unmanaged[0].Id, unmanaged[1].Id)
	}
}

func TestApply_UnsupportedField(t *testing.T) {
	report := buildReport()
	err := sorter.Apply(report, sorter.Options{Field: "unknown", Order: sorter.Asc})
	if err == nil {
		t.Fatal("expected error for unsupported field, got nil")
	}
}

func TestApply_NilReport(t *testing.T) {
	err := sorter.Apply(nil, sorter.Options{Field: sorter.ByType, Order: sorter.Asc})
	if err != nil {
		t.Fatalf("expected no error for nil report, got: %v", err)
	}
}
