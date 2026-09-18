package utils

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	defaultLxcfsParentDir = "lxcfs-on-k8s"
	// DefaultLxcfsPath is the default host path of the lxcfs mount root.
	DefaultLxcfsPath = "/var/lib/lxcfs-on-k8s/lxcfs"
	// procVolumePrefix is the prefix of the pod volume names used to mount lxcfs proc files.
	procVolumePrefix = "lxcfs-proc-"
)

// dns1123LabelRegexp matches the pod volume names kubernetes accepts:
// RFC 1123 labels, i.e. lowercase alphanumeric plus '-' (see ProcVolumeName).
var dns1123LabelRegexp = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

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

// ProcVolumeName returns the pod volume name used to mount a lxcfs proc
// file, e.g. "cpuinfo" becomes "lxcfs-proc-cpuinfo".
func ProcVolumeName(name string) string {
	return procVolumePrefix + strings.ReplaceAll(name, "/", "-")
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
	if err := ValidateProcFiles(files); err != nil {
		return nil, err
	}
	return files, nil
}

// ValidateProcFiles validates a normalized list of proc file names (as
// produced by ParseProcFiles): every entry must yield a valid pod volume
// name and two distinct entries must not collide on the same volume name.
func ValidateProcFiles(names []string) error {
	owners := make(map[string]string, len(names))
	for _, name := range names {
		if err := validateProcFile(name); err != nil {
			return err
		}
		volume := ProcVolumeName(name)
		if other, ok := owners[volume]; ok && other != name {
			return fmt.Errorf("proc files %q and %q collide on pod volume name %q", other, name, volume)
		}
		owners[volume] = name
	}
	return nil
}

func validateProcFile(name string) error {
	if strings.HasPrefix(name, "/") {
		return fmt.Errorf("invalid proc file %q: expected a name relative to /proc, e.g. \"cpuinfo\" or \"/proc/cpuinfo\"", name)
	}
	if name != path.Clean(name) || name == "." || name == ".." || strings.HasPrefix(name, "../") {
		return fmt.Errorf("invalid proc file %q", name)
	}
	if strings.SplitN(name, "/", 2)[0] == "proc" {
		if rest := strings.TrimPrefix(name, "proc/"); rest != name {
			return fmt.Errorf("invalid proc file %q: names are relative to /proc, did you mean %q", name, rest)
		}
		return fmt.Errorf("invalid proc file %q: names are relative to /proc", name)
	}
	for _, segment := range strings.Split(name, "/") {
		if strings.HasPrefix(segment, "-") || strings.HasSuffix(segment, "-") {
			return fmt.Errorf("invalid proc file %q: path segments must not start or end with '-'", name)
		}
	}
	volume := ProcVolumeName(name)
	if len(volume) > 63 || !dns1123LabelRegexp.MatchString(volume) {
		return fmt.Errorf("proc file %q yields pod volume name %q: volume names must be lowercase RFC 1123 labels of at most 63 characters", name, volume)
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
