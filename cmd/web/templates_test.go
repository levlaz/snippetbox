package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	// create slice of anon structs containing the test case name,
	// input, and expected output
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "UTC",
			input:    time.Date(2024, time.May, 15, 10, 0, 0, 0, time.UTC),
			expected: "15 May 2024 at 10:00",
		},
		{
			name:     "Empty",
			input:    time.Time{},
			expected: "",
		},
		{
			name:     "CET",
			input:    time.Date(2024, time.May, 15, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			expected: "15 May 2024 at 09:00",
		},
	}

	// loop over test cases
	for _, tt := range tests {
		// run subtest for each test case
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.input)
			if hd != tt.expected {
				t.Errorf("got %q; want %q", hd, tt.expected)
			}
		})
	}
}
