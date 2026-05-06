// Package reporter provides human-readable rendering of drift summaries.
//
// It supports multiple output formats:
//
//   - FormatTable: a padded tabular layout suitable for terminal output.
//   - FormatCompact: a single-line key=value representation useful for
//     log-based pipelines and CI output.
//
// Formats can be selected by name using ParseFormat, which returns an error
// for unrecognised format strings. This is useful when accepting a format
// value from a CLI flag or environment variable.
//
// Example usage:
//
//	fmt, err := reporter.ParseFormat(flagValue)
//	if err != nil {
//		log.Fatal(err)
//	}
//	r := reporter.New(fmt, os.Stdout)
//	if err := r.Render(s); err != nil {
//		log.Fatal(err)
//	}
package reporter
