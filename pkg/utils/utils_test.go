package utils

import (
	"reflect"
	"strings"
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
			path:    "/var/lib/lxcfs-on-k8s/test/",
			wantErr: false,
		},
		{
			name:    "valid path without trailing slash",
			path:    "/var/lib/lxcfs-on-k8s/test",
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
			path:    "/lxcfs-on-k8s/test",
			wantErr: false,
		},
		{
			name:    "path is just /lxcfs-on-k8s",
			path:    "/lxcfs-on-k8s",
			wantErr: true,
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
		{
			name:    "path with multiple trailing slashes",
			path:    "/var/lib/lxcfs-on-k8s/test///",
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

func TestParseProcFiles(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    []string
		wantErr bool
	}{
		{
			name: "empty disables proc mounts",
			in:   "",
			want: []string{},
		},
		{
			name: "single name",
			in:   "cpuinfo",
			want: []string{"cpuinfo"},
		},
		{
			name: "absolute proc path",
			in:   "/proc/cpuinfo",
			want: []string{"cpuinfo"},
		},
		{
			name: "mixed forms with spaces",
			in:   "cpuinfo, /proc/meminfo ,swaps",
			want: []string{"cpuinfo", "meminfo", "swaps"},
		},
		{
			name: "duplicates are deduplicated",
			in:   "cpuinfo,/proc/cpuinfo,meminfo",
			want: []string{"cpuinfo", "meminfo"},
		},
		{
			name: "trailing comma is tolerated",
			in:   "cpuinfo,",
			want: []string{"cpuinfo"},
		},
		{
			name: "directory entry",
			in:   "pressure",
			want: []string{"pressure"},
		},
		{
			name: "default list round trip",
			in:   strings.Join(DefaultProcFiles, ","),
			want: DefaultProcFiles,
		},
		{
			name:    "absolute path outside /proc",
			in:      "/etc/passwd",
			wantErr: true,
		},
		{
			name:    "leading parent directory",
			in:      "../etc/passwd",
			wantErr: true,
		},
		{
			name:    "nested parent directory traversal",
			in:      "proc/../../etc/passwd",
			wantErr: true,
		},
		{
			name:    "traversal through proc prefix",
			in:      "/proc/../etc/passwd",
			wantErr: true,
		},
		{
			name:    "dot",
			in:      ".",
			wantErr: true,
		},
		{
			name:    "dotdot",
			in:      "..",
			wantErr: true,
		},
		{
			name:    "non canonical path",
			in:      "cpuinfo//online",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseProcFiles(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseProcFiles(%q) error = %v, wantErr %v", tt.in, err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseProcFiles(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}
