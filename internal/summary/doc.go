// Package summary computes and renders aggregated statistics from a
// parsed driftctl report.
//
// Usage:
//
//	report, err := parser.ParseFile("report.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	stats := summary.Compute(report)
//	summary.Write(os.Stdout, stats)
//
// The Compute function derives counts for managed, unmanaged, deleted,
// and drifted resources as well as an overall coverage percentage.
// Write formats the Stats struct into a tab-aligned table suitable for
// terminal output or inclusion in audit logs.
package summary
