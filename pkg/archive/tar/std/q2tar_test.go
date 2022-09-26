package tarq

import (
	"bytes"
	"testing"
)

func checkBuilder[T any](comp func(a, b T) (same bool)) func(got, expected T) func(*testing.T) {
	return func(got, expected T) func(*testing.T) {
		return func(t *testing.T) {
			var same bool = comp(got, expected)
			if !same {
				t.Errorf("Unexpected value got.\n")
				t.Errorf("Expected: %v\n", expected)
				t.Fatalf("Got:      %v\n", got)
			}
		}
	}
}

func check[T comparable](got, expected T) func(*testing.T) {
	return checkBuilder(func(a, b T) (same bool) { return a == b })(got, expected)
}

var checkBytes func(got, expected []byte) func(*testing.T) = checkBuilder(
	func(a, b []byte) (same bool) { return 0 == bytes.Compare(a, b) },
)
