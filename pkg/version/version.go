// Package version is used to inject build version and time
package version

import "fmt"

// injected values by ldflags
// nolint:gochecknoglobals
var (
	version   string // $ git rev-parse HEAD | cut -c -10
	buildtime string // $ date +%Y%m%d%H%M%S
)

// FullVersion returns the full version string with format {version}-{buildtime}
// ex. abababab-20200202000000
func FullVersion() string {
	return fmt.Sprintf("%s-%s", Version(), BuildTime())
}

// Version returns the version.
func Version() string {
	if len(version) == 0 {
		return "UNKNOWN_VERSION"
	}
	return version
}

// BuildTime returns the buildtime.
func BuildTime() string {
	if len(buildtime) == 0 {
		return "UNKNOWN_BUILD_TIME"
	}
	return buildtime
}
