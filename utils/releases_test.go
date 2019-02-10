package utils

import (
	"testing"

	"k8s.io/helm/pkg/proto/hapi/release"
)

func TestContainsRelease(t *testing.T) {
	queryRelease := release.Release{Name: "testRelease"}
	targetReleases := make([]*release.Release, 0)
	targetReleases = append(targetReleases, &queryRelease)
	containsRelease := ContainsRelease(&queryRelease, targetReleases)
	expected := true
	if containsRelease != expected {
		t.Errorf("Incorrect bool returned, got: %t, want: %t.", containsRelease, expected)
	}
}
