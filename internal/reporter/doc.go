// Package reporter provides human-readable rendering of drift summaries.
//
// It supports multiple output formats:
//
//   - FormatTable: a padded tabular layout suitable for terminal output.
//   - FormatCompact: a single-line key=value representation useful for
//     log-based pipelines and CI output.
//
// Example usage:
//
//	r := reporter.New(reporter.FormatTable, os.Stdout)
//	if err := r.Render(s); err != nil {
//		log.Fatal(err)
//	}
package reporter
