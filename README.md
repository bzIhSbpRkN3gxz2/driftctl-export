# driftctl-export

> Utility to export Terraform drift reports into structured JSON/CSV for audit pipelines

---

## Installation

```bash
go install github.com/yourorg/driftctl-export@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/driftctl-export.git
cd driftctl-export
go build -o driftctl-export .
```

---

## Usage

Run `driftctl scan` and pipe the output directly into `driftctl-export`:

```bash
driftctl scan --output json://drift.json
driftctl-export --input drift.json --format csv --output report.csv
```

**Available flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--input` | Path to driftctl JSON output file | `drift.json` |
| `--format` | Output format: `json` or `csv` | `json` |
| `--output` | Path for the exported report file | `report.json` |
| `--filter` | Filter by resource type (e.g. `aws_s3_bucket`) | none |

**Example — export filtered CSV for audit:**

```bash
driftctl-export \
  --input drift.json \
  --format csv \
  --filter aws_iam_role \
  --output iam-drift-report.csv
```

---

## Requirements

- Go 1.21+
- [driftctl](https://github.com/snyk/driftctl) v0.39+

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

This project is licensed under the [MIT License](LICENSE).