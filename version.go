package hypersyncgo

import "runtime/debug"

// version returns the module version from build info (set by the go tool
// from the git tag, e.g. v0.1.1). Falls back to "unknown" during development.
func version() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "unknown"
}
