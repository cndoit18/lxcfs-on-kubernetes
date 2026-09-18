package utils

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

const (
	defaultLxcfsParentDir = "lxcfs-on-k8s"
)

// DefaultProcFiles is the default set of lxcfs proc files mounted into pods.
var DefaultProcFiles = []string{
	"cpuinfo",
	"diskstats",
	"meminfo",
	"stat",
	"swaps",
	"uptime",
	"loadavg",
	"pressure",
	"slabinfo",
}

// ParseProcFiles parses a comma separated list of proc file names into
// normalized names relative to /proc. Both "cpuinfo" and "/proc/cpuinfo"
// forms are accepted and considered equivalent. An empty string produces an
// empty (non-nil) list, which disables all proc mounts.
func ParseProcFiles(s string) ([]string, error) {
	files := []string{}
	seen := make(map[string]struct{})
	for _, entry := range strings.Split(s, ",") {
		name := strings.TrimSpace(entry)
		name = strings.TrimPrefix(name, "/proc/")
		if name == "" {
			continue
		}
		if err := validateProcFile(name); err != nil {
			return nil, err
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		files = append(files, name)
	}
	return files, nil
}

func validateProcFile(name string) error {
	if strings.HasPrefix(name, "/") {
		return fmt.Errorf("invalid proc file %q: expected a name relative to /proc, e.g. \"cpuinfo\" or \"/proc/cpuinfo\"", name)
	}
	if name != path.Clean(name) || name == "." || name == ".." || strings.HasPrefix(name, "../") {
		return fmt.Errorf("invalid proc file %q", name)
	}
	return nil
}

func EnsureLxcfsParentDir(path string) error {
	lxcDir := filepath.Dir(strings.TrimRight(path, "/"))
	if !strings.HasSuffix(lxcDir, defaultLxcfsParentDir) {
		return fmt.Errorf("lxcfs path %s is not valid, it's parent directory should be '%s'", path, defaultLxcfsParentDir)
	}
	return nil
}
