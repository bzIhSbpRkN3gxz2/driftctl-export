// Package pipeline provides an end-to-end orchestration layer for
// processing Terraform drift reports produced by driftctl.
//
// A pipeline run performs the following steps in order:
//
//  1. Load — reads one or more report files matching an input glob
//     pattern using the loader package.
//
//  2. Validate — ensures each loaded report conforms to the expected
//     schema using the validator package.
//
//  3. Filter — applies resource-level inclusion/exclusion rules via
//     the filter package.
//
//  4. Export — serialises the filtered reports to JSON or CSV at the
//     configured output path using the exporter package.
//
//  5. Summarise — computes aggregate statistics over the final set of
//     resources using the summary package.
//
// Typical usage:
//
//	result, err := pipeline.Run(pipeline.Options{
//		InputGlob:  "./reports/*.json",
//		OutputPath: "./out/drift.csv",
//		Format:     "csv",
//	})
package pipeline
