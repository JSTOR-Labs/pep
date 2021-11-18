package which

import "os/exec"

// Which is a package that discovers and caches the location of binaries

var whichCache map[string]string

func init() {
	whichCache = make(map[string]string)
}

func LookupExecutable(name string) string {
	if path, ok := whichCache[name]; ok {
		return path
	}
	path, err := exec.LookPath(name)
	if err != nil {
		return ""
	}
	whichCache[name] = path
	return path
}
