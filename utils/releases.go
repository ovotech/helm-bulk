package utils

import "k8s.io/helm/pkg/proto/hapi/release"

//ContainsRelease returns a bool indicating whether the provided Release is
//in the provided slice
func ContainsRelease(queryRelease *release.Release,
	targetReleases []*release.Release) (contains bool) {
	for _, release := range targetReleases {
		if queryRelease.GetName() == release.GetName() {
			contains = true
			break
		}
	}
	return
}
