package utils

import (
	"testing"
)

func TestGetNegativeReason(t *testing.T) {
	reason := GetNegativeReason()
	if reason == "" {
		t.Error("GetNegativeReason returned an empty string")
	}
}

func TestGetPositiveReason(t *testing.T) {
	reason := GetPositiveReason()
	if reason == "" {
		t.Error("GetPositiveReason returned an empty string")
	}
}

func TestGetNegativeReasonRandomness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 50; i++ {
		seen[GetNegativeReason()] = true
	}
	if len(seen) < 2 {
		t.Error("GetNegativeReason does not appear to be random — got same result 50 times")
	}
}

func TestGetPositiveReasonRandomness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 50; i++ {
		seen[GetPositiveReason()] = true
	}
	if len(seen) < 2 {
		t.Error("GetPositiveReason does not appear to be random — got same result 50 times")
	}
}

func TestReasonsLoaded(t *testing.T) {
	if lenNoReasons == 0 {
		t.Error("No negative reasons loaded from reasons.json")
	}
	if lenYesReasons == 0 {
		t.Error("No positive reasons loaded from reasons.json")
	}
}
