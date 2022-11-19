package version

import "runtime"

type Version struct {
	GoVersion  string
	GoOs       string
	GoArch     string
	GitVersion string
	GitCommit  string
	BuildDate  string
}

var ver = Version{
	GoVersion:  runtime.Version(),
	GoOs:       runtime.GOOS,
	GoArch:     runtime.GOARCH,
	GitVersion: "none",
	GitCommit:  "none",
	BuildDate:  "none",
}

func Get() Version {
	return ver
}
