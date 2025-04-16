package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expect T) {
	t.Helper()

	if actual != expect {
		t.Errorf("got %v; want %v", actual, expect)
	}
}

func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
	}
}
