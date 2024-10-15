package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	// init new time.Time object and pass into humanDate()
	tm := time.Date(2024, time.May, 15, 10, 0, 0, 0, time.UTC)
	hd := humanDate(tm)

	// check that output string has expected format
	expected := "15 May 2024 at 10:00"
	if hd != expected {
		t.Errorf("got %q; want %q", hd, expected)
	}

	// // uncomment to force failure by passing in known invalid time
	// notAMatch := time.Date(2024, time.May, 15, 10, 0, 0, 0, time.UTC)
	// hd = humanDate(notAMatch)
	// expected = "15 May 2025 at 10:00"
	// if hd != expected {
	// 	t.Errorf("got %q; want %q", hd, expected)
	// }
}
