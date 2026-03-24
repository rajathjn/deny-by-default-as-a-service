package favicon

import (
	"testing"
)

func TestGetFavicon(t *testing.T) {
	favicon, err := GetFavicon()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(favicon) == 0 {
		t.Error("Expected favicon data, got empty")
	}
}