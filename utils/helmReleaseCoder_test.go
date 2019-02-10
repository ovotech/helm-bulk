package utils

import (
	"testing"

	"k8s.io/helm/pkg/proto/hapi/release"
)

func TestEncodeReleases(t *testing.T) {
	release := release.Release{Name: "testRelease"}
	actualString, err := EncodeRelease(&release)
	if err != nil {
		t.Error("Error encoding Helm Release", err)
	}
	expectedString := "H4sIAAAAAAAC/+LiLkktLglKzUlNLE4FBAAA//9q7y4QDQAAAA=="
	if actualString != expectedString {
		t.Errorf("Release string was incorrect, got: %s, want: %s.",
			actualString, expectedString)
	}
}

func TestDecodeReleases(t *testing.T) {
	decodedRelease, err := DecodeRelease("H4sIAAAAAAAC/+LiLkktLglKzUlNLE4FBAAA//9q7y4QDQAAAA==")
	if err != nil {
		t.Error("Error decoding Helm Release", err)
	}
	actualString := decodedRelease.GetName()
	expectedSTring := "testRelease"
	if actualString != expectedSTring {
		t.Errorf("Decoded release name was incorrect, got: %s, want: %s.",
			actualString, expectedSTring)
	}
}
