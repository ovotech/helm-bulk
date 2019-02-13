package cmd

import (
	"testing"
)

func TestTextFilename(t *testing.T) {
	actualString := textFilename()
	expectedString := "helm-releases.txt"
	if actualString != expectedString {
		t.Errorf("Release string was incorrect, got: %s, want: %s.",
			actualString, expectedString)
	}
}

func TestArchiveFileName(t *testing.T) {
	actualString := archiveFilename()
	expectedString := "helm-releases.tar.gz"
	if actualString != expectedString {
		t.Errorf("Release string was incorrect, got: %s, want: %s.",
			actualString, expectedString)
	}
}
