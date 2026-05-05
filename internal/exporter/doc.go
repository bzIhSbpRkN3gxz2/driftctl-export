// Package exporter provides functionality to serialise parsed DriftReport
// data into structured output formats suitable for audit pipelines.
//
// Supported formats:
//
//   - JSON  — pretty-printed JSON encoding of the full report structure.
//   - CSV   — flat CSV with columns: type, id, status; one row per drifted
//             resource.
//
// Basic usage:
//
//	e := exporter.New(exporter.FormatJSON, os.Stdout)
//	if err := e.Export(report); err != nil {
//		log.Fatal(err)
//	}
//
// To write directly to a file use NewFileExporter which creates the
// destination file and returns the underlying *os.File so the caller can
// close it when done.
package exporter
