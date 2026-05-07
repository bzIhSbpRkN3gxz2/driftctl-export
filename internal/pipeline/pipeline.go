// Package pipeline orchestrates the full drift report processing
// workflow: load, validate, filter, export, and summarize.
package pipeline

import (
	"fmt"

	"github.com/example/driftctl-export/internal/exporter"
	"github.com/example/driftctl-export/internal/filter"
	"github.com/example/driftctl-export/internal/loader"
	"github.com/example/driftctl-export/internal/summary"
	"github.com/example/driftctl-export/internal/validator"
)

// Options holds configuration for a pipeline run.
type Options struct {
	// InputGlob is a file glob pattern for input report files.
	InputGlob string
	// OutputPath is the destination file for the exported report.
	OutputPath string
	// Format is the export format: "json" or "csv".
	Format string
	// FilterOptions controls which resources are included.
	FilterOptions filter.Options
	// ReporterFormat controls summary output: "table" or "compact".
	ReporterFormat string
}

// Result holds the outcome of a pipeline run.
type Result struct {
	FilesProcessed int
	ResourcesTotal int
	Summary        summary.Summary
}

// Run executes the full pipeline with the given options.
func Run(opts Options) (*Result, error) {
	reports, err := loader.LoadGlob(opts.InputGlob)
	if err != nil {
		return nil, fmt.Errorf("pipeline: load: %w", err)
	}
	if len(reports) == 0 {
		return nil, fmt.Errorf("pipeline: no reports matched glob %q", opts.InputGlob)
	}

	for i, r := range reports {
		if err := validator.Validate(r); err != nil {
			return nil, fmt.Errorf("pipeline: validate report[%d]: %w", i, err)
		}
	}

	filtered := filter.Apply(reports, opts.FilterOptions)

	exp, err := exporter.NewFileExporter(opts.OutputPath, opts.Format)
	if err != nil {
		return nil, fmt.Errorf("pipeline: exporter: %w", err)
	}
	if err := exp.Export(filtered); err != nil {
		return nil, fmt.Errorf("pipeline: export: %w", err)
	}

	sum := summary.Compute(filtered)

	total := 0
	for _, r := range filtered {
		total += len(r.Managed) + len(r.Unmanaged) + len(r.Deleted)
	}

	return &Result{
		FilesProcessed: len(reports),
		ResourcesTotal: total,
		Summary:        sum,
	}, nil
}
