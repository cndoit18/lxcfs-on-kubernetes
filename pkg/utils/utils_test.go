package utils

import (
	"testing"
)

func TestEnsureLxcfsParentDir(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "valid path with trailing slash",
			path:    "/var/lib/lxc/test/",
			wantErr: false,
		},
		{
			name:    "valid path without trailing slash",
			path:    "/var/lib/lxc/test",
			wantErr: false,
		},
		{
			name:    "invalid parent directory",
			path:    "/var/lib/lxcfs/test",
			wantErr: true,
		},
		{
			name:    "parent directory is not lxc",
			path:    "/foo/bar/test",
			wantErr: true,
		},
		{
			name:    "parent directory is lxc at root",
			path:    "/lxc/test",
			wantErr: false,
		},
		{
			name:    "path is just /lxc",
			path:    "/lxc",
			wantErr: true,
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
		{
			name:    "path with multiple trailing slashes",
			path:    "/var/lib/lxc/test///",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EnsureLxcfsParentDir(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnsureLxcfsParentDir(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
		})
	}
}
