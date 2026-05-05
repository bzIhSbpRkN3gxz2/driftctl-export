// Package parser provides functionality for reading and deserialising
// driftctl JSON output files into Go structs that can be consumed by
// downstream exporters (JSON, CSV, etc.).
//
// Typical usage:
//
//	report, err := parser.ParseFile("drift-report.json")
//	if err != nil {
//		log.Fatalf("failed to parse report: %v", err)
//	}
//	fmt.Printf("Drifted resources: %d\n", report.Summary.DriftedCount)
//
// The package intentionally does no output formatting; see the exporter
// sub-packages for JSON and CSV rendering.
package parser
