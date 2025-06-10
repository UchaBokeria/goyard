package tests

import (
	"testing"

	"github.com/yourorg/goyard"
)

func TestVersion(t *testing.T) {
	if goyard.Version == "" {
		t.Error("Version should not be empty")
	}
}

// More tests will be added as the framework develops
