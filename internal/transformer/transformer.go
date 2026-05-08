package transformer

import (
	"strings"

	"github.com/snyk/driftctl-export/internal/models"
)

// Option holds configuration for the transformation pipeline.
type Option struct {
	NormalizeIDs   bool
	UppercaseTypes bool
	StripPrefix    string
}

// Apply transforms a DriftReport in-place according to the provided options.
// It normalises resource IDs and/or types to make downstream comparisons
// deterministic across different Terraform providers.
func Apply(report *models.DriftReport, opt Option) *models.DriftReport {
	if report == nil {
		return nil
	}

	for i := range report.Unmanaged {
		report.Unmanaged[i].ResourceID = transformID(report.Unmanaged[i].ResourceID, opt)
		report.Unmanaged[i].ResourceType = transformType(report.Unmanaged[i].ResourceType, opt)
	}

	for i := range report.Managed {
		report.Managed[i].ResourceID = transformID(report.Managed[i].ResourceID, opt)
		report.Managed[i].ResourceType = transformType(report.Managed[i].ResourceType, opt)
	}

	for i := range report.Drifted {
		report.Drifted[i].ResourceID = transformID(report.Drifted[i].ResourceID, opt)
		report.Drifted[i].ResourceType = transformType(report.Drifted[i].ResourceType, opt)
	}

	return report
}

func transformID(id string, opt Option) string {
	if opt.StripPrefix != "" {
		id = strings.TrimPrefix(id, opt.StripPrefix)
	}
	if opt.NormalizeIDs {
		id = strings.ToLower(strings.TrimSpace(id))
	}
	return id
}

func transformType(t string, opt Option) string {
	if opt.UppercaseTypes {
		return strings.ToUpper(strings.TrimSpace(t))
	}
	return strings.TrimSpace(t)
}
