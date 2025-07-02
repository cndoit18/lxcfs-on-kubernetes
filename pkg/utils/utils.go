package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

const (
	defaultLxcfsParentDir = "lxcfs-on-k8s"
)

func EnsureLxcfsParentDir(path string) error {
	lxcDir := filepath.Dir(strings.TrimRight(path, "/"))
	if !strings.HasSuffix(lxcDir, defaultLxcfsParentDir) {
		return fmt.Errorf("lxcfs path %s is not valid, it's parent directory should be 'lxc'", path)
	}
	return nil
}
