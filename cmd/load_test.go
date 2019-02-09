package cmd

import (
	"bytes"
	"testing"

	"k8s.io/helm/pkg/proto/hapi/release"
)

func TestAddReleasesToBuffer(t *testing.T) {
	helmRelease := release.Release{Name: "dummy"}
	releases := make([]*release.Release, 0)
	releases = append(releases, &helmRelease)
	var buffer bytes.Buffer
	addReleasesToBuffer(releases, &buffer)
	expectedString := "    dummy\n"
	actualString := buffer.String()
	if actualString != expectedString {
		t.Errorf("Release string was incorrect, got: %s, want: %s.",
			actualString, expectedString)
	}
}

func TestAddHeaderToBuffer(t *testing.T) {
	var buffer bytes.Buffer
	addHeaderToBuffer("dummyHeader", &buffer)
	expectedString := "dummyHeader\n\n"
	actualString := buffer.String()
	if actualString != expectedString {
		t.Errorf("Release string was incorrect, got: %s, want: %s.",
			actualString, expectedString)
	}
}
