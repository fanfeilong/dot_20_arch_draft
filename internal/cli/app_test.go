package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestHelp(t *testing.T) {
	var out bytes.Buffer

	if err := runWithIO([]string{"help"}, &out, "test"); err != nil {
		t.Fatalf("runWithIO returned error: %v", err)
	}

	if !strings.Contains(out.String(), "d2a init <target-dir>") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestVersion(t *testing.T) {
	var out bytes.Buffer

	if err := runWithIO([]string{"version"}, &out, "v0.0.1"); err != nil {
		t.Fatalf("runWithIO returned error: %v", err)
	}

	if got := out.String(); got != "v0.0.1\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}
