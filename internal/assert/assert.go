package assert

import (
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expect T) {
	t.Helper()

	if actual != expect {
		t.Errorf("got %v; want %v", actual, expect)
	}
}
