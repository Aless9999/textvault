package assert

import (
	"testing"
)

func Equals[T comparable](t *testing.T, actual T, expected T) {
	t.Helper()
	if actual != expected {
		t.Errorf("Expected %v ,but actual %v", expected, actual)
	}
}
