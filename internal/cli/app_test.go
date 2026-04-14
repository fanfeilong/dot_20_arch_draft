package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestHelp(t *testing.T) {
	var out bytes.Buffer

	if err := runWithIO([]string{"help"}, &out); err != nil {
		t.Fatalf("runWithIO returned error: %v", err)
	}

	if !strings.Contains(out.String(), "d2a init <target-dir>") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}
