package main

import "runtime"

var (
	// ReleaseNumber hold release number assigned manually.
	ReleaseNumber string

	// BuildNumber hold number assigned automatically externally by "go build".
	BuildNumber string

	// BuildTime holds build time assigned externally by "go build"
	BuildTime string

	BuildPlatform = runtime.GOOS

	// BuildGitHash holds git commit's hash, used for building this module.
	BuildGitHash string
)
