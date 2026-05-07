// Package validator provides pre-export validation for driftctl DriftReport
// structures.
//
// It checks that required fields such as ResourceType and ResourceID are
// present on each resource entry before the report is handed off to the
// exporter or summary pipeline.  Any issues are collected and returned as a
// single *ValidationError so callers receive a complete picture of all
// problems in one pass.
//
// Typical usage:
//
//	report, err := parser.ParseFile("drift.json")
//	if err != nil { ... }
//
//	if err := validator.Validate(report); err != nil {
//		log.Fatalf("invalid report: %v", err)
//	}
package validator
