package sample_test

import (
	"testing"

	"github.com/kynrai/go-template/pkg/sample"
)

func TestSample(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name     string
		expected string
	}{
		{
			name:     "test 1",
			expected: "sample",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := sample.Sample()
			if result != tc.expected {
				t.Fatalf("[Sample][%s]: got: %v, expected: %v", tc.name, result, tc.expected)
			}
		})
	}
}
