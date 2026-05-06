package version_test

import (
	"strings"
	"testing"

	"github.com/example/driftctl-export/internal/version"
)

func TestGet_DefaultValues(t *testing.T) {
	info := version.Get()

	if info.Version == "" {
		t.Error("expected Version to be non-empty")
	}
	if info.Commit == "" {
		t.Error("expected Commit to be non-empty")
	}
	if info.Date == "" {
		t.Error("expected Date to be non-empty")
	}
	if info.GoVer == "" {
		t.Error("expected GoVer to be non-empty")
	}
	if info.OS == "" {
		t.Error("expected OS to be non-empty")
	}
	if info.Arch == "" {
		t.Error("expected Arch to be non-empty")
	}
}

func TestGet_DefaultVersion(t *testing.T) {
	info := version.Get()
	// Default value set in the package.
	if info.Version != "dev" {
		t.Errorf("expected default version 'dev', got %q", info.Version)
	}
}

func TestInfo_String_ContainsVersion(t *testing.T) {
	info := version.Get()
	s := info.String()

	if !strings.Contains(s, info.Version) {
		t.Errorf("String() %q does not contain version %q", s, info.Version)
	}
	if !strings.Contains(s, info.Commit) {
		t.Errorf("String() %q does not contain commit %q", s, info.Commit)
	}
	if !strings.Contains(s, info.GoVer) {
		t.Errorf("String() %q does not contain Go version %q", s, info.GoVer)
	}
}

func TestInfo_String_Format(t *testing.T) {
	info := version.Info{
		Version: "1.2.3",
		Commit:  "abc1234",
		Date:    "2024-01-15T10:00:00Z",
		GoVer:   "go1.22.0",
		OS:      "linux",
		Arch:    "amd64",
	}

	s := info.String()
	expected := "1.2.3 (commit=abc1234, built=2024-01-15T10:00:00Z, linux/amd64, go1.22.0)"
	if s != expected {
		t.Errorf("expected %q, got %q", expected, s)
	}
}
