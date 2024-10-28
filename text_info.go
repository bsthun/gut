package gut

import (
	"runtime/debug"
)

var Version string

var Commit = func() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		var hash string
		var modified string
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				hash = setting.Value[:7]
			}
			if setting.Key == "vcs.modified" {
				if setting.Value == "false" {
					modified = ".c" // Clean build
				} else {
					modified = ".d" // Dirty build
				}
			}
		}
		if hash == "" || modified == "" {
			return ""
		}
		return hash + modified + "." + Build
	}
	return ""
}()
