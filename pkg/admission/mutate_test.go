/*
 *     Copyright 2021 Yinan Li <cndoit18@outlook.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package admission

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"

	"github.com/cndoit18/lxcfs-on-kubernetes/pkg/utils"
)

const testMutatePath = "/var/lib/lxcfs-on-k8s/lxcfs"

func newMutate(procFiles []string) *mutate {
	m := &mutate{}
	for _, opt := range []admissionOption{
		WithMutatePath(testMutatePath),
		WithMutateProcFiles(procFiles),
	} {
		opt(m)
	}
	return m
}

func volumeMountNames(mounts []corev1.VolumeMount) []string {
	names := make([]string, 0, len(mounts))
	for _, m := range mounts {
		names = append(names, m.Name)
	}
	return names
}

func volumeNames(volumes []corev1.Volume) []string {
	names := make([]string, 0, len(volumes))
	for _, v := range volumes {
		names = append(names, v.Name)
	}
	return names
}

func TestEnsureVolumeMount(t *testing.T) {
	tests := []struct {
		name      string
		procFiles []string
		existing  []corev1.VolumeMount
		wantNames []string
	}{
		{
			name:      "default proc files",
			procFiles: utils.DefaultProcFiles,
			wantNames: []string{
				"lxcfs-proc-cpuinfo",
				"lxcfs-proc-diskstats",
				"lxcfs-proc-loadavg",
				"lxcfs-proc-meminfo",
				"lxcfs-proc-pressure",
				"lxcfs-proc-slabinfo",
				"lxcfs-proc-stat",
				"lxcfs-proc-swaps",
				"lxcfs-proc-uptime",
				"lxcfs-root-parent-dir",
				"lxcfs-sys-devices-system-cpu",
			},
		},
		{
			name:      "custom proc files for old runc allowlist",
			procFiles: []string{"cpuinfo", "diskstats", "meminfo", "stat", "swaps", "uptime"},
			wantNames: []string{
				"lxcfs-proc-cpuinfo",
				"lxcfs-proc-diskstats",
				"lxcfs-proc-meminfo",
				"lxcfs-proc-stat",
				"lxcfs-proc-swaps",
				"lxcfs-proc-uptime",
				"lxcfs-root-parent-dir",
				"lxcfs-sys-devices-system-cpu",
			},
		},
		{
			name:      "empty proc files only keeps fixed mounts",
			procFiles: []string{},
			wantNames: []string{
				"lxcfs-root-parent-dir",
				"lxcfs-sys-devices-system-cpu",
			},
		},
		{
			name:      "existing lxcfs mounts are not duplicated",
			procFiles: []string{"cpuinfo", "meminfo"},
			existing: []corev1.VolumeMount{
				{Name: "lxcfs-proc-cpuinfo", MountPath: "/somewhere/else"},
			},
			wantNames: []string{
				"lxcfs-proc-cpuinfo",
				"lxcfs-proc-meminfo",
				"lxcfs-root-parent-dir",
				"lxcfs-sys-devices-system-cpu",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newMutate(tt.procFiles)
			got := m.ensureVolumeMount(tt.existing)
			if !reflect.DeepEqual(volumeMountNames(got), tt.wantNames) {
				t.Errorf("ensureVolumeMount() names = %v, want %v", volumeMountNames(got), tt.wantNames)
			}
		})
	}
}

func TestEnsureVolumeMountFields(t *testing.T) {
	m := newMutate([]string{"cpuinfo", "pressure"})
	mounts := m.ensureVolumeMount(nil)

	byName := make(map[string]corev1.VolumeMount, len(mounts))
	for _, m := range mounts {
		byName[m.Name] = m
	}

	wantMountPaths := map[string]string{
		"lxcfs-proc-cpuinfo":           "/proc/cpuinfo",
		"lxcfs-proc-pressure":          "/proc/pressure",
		"lxcfs-root-parent-dir":        "/var/lib/lxcfs-on-k8s",
		"lxcfs-sys-devices-system-cpu": "/sys/devices/system/cpu",
	}
	if len(byName) != len(wantMountPaths) {
		t.Fatalf("ensureVolumeMount() returned %d mounts, want %d", len(byName), len(wantMountPaths))
	}
	for name, wantPath := range wantMountPaths {
		mount, ok := byName[name]
		if !ok {
			t.Fatalf("ensureVolumeMount() missing volume mount %q", name)
		}
		if mount.MountPath != wantPath {
			t.Errorf("mount %q path = %q, want %q", name, mount.MountPath, wantPath)
		}
		if !mount.ReadOnly {
			t.Errorf("mount %q should be read-only", name)
		}
		wantPropagation := corev1.MountPropagationNone
		if name == rootParentDirVolume {
			wantPropagation = corev1.MountPropagationHostToContainer
		}
		if mount.MountPropagation == nil || *mount.MountPropagation != wantPropagation {
			t.Errorf("mount %q propagation = %v, want %v", name, mount.MountPropagation, wantPropagation)
		}
	}
}

func TestEnsureVolume(t *testing.T) {
	tests := []struct {
		name      string
		procFiles []string
		existing  []corev1.Volume
		wantNames []string
	}{
		{
			name:      "custom proc files",
			procFiles: []string{"cpuinfo", "meminfo"},
			wantNames: []string{
				"lxcfs-proc-cpuinfo",
				"lxcfs-proc-meminfo",
				"lxcfs-root-parent-dir",
				"lxcfs-sys-devices-system-cpu",
			},
		},
		{
			name:      "empty proc files only keeps fixed volumes",
			procFiles: []string{},
			wantNames: []string{
				"lxcfs-root-parent-dir",
				"lxcfs-sys-devices-system-cpu",
			},
		},
		{
			name:      "existing lxcfs volumes are not duplicated",
			procFiles: []string{"cpuinfo"},
			existing: []corev1.Volume{
				{Name: "lxcfs-proc-cpuinfo"},
			},
			wantNames: []string{
				"lxcfs-proc-cpuinfo",
				"lxcfs-root-parent-dir",
				"lxcfs-sys-devices-system-cpu",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newMutate(tt.procFiles)
			got := m.ensureVolume(tt.existing)
			if !reflect.DeepEqual(volumeNames(got), tt.wantNames) {
				t.Errorf("ensureVolume() names = %v, want %v", volumeNames(got), tt.wantNames)
			}
		})
	}
}

func TestEnsureVolumeHostPaths(t *testing.T) {
	m := newMutate([]string{"cpuinfo", "pressure"})
	volumes := m.ensureVolume(nil)

	byName := make(map[string]corev1.Volume, len(volumes))
	for _, v := range volumes {
		byName[v.Name] = v
	}

	wantHostPaths := map[string]string{
		"lxcfs-proc-cpuinfo":           testMutatePath + "/proc/cpuinfo",
		"lxcfs-proc-pressure":          testMutatePath + "/proc/pressure",
		"lxcfs-root-parent-dir":        "/var/lib/lxcfs-on-k8s",
		"lxcfs-sys-devices-system-cpu": testMutatePath + "/sys/devices/system/cpu",
	}
	if len(byName) != len(wantHostPaths) {
		t.Fatalf("ensureVolume() returned %d volumes, want %d", len(byName), len(wantHostPaths))
	}
	for name, wantPath := range wantHostPaths {
		volume, ok := byName[name]
		if !ok {
			t.Fatalf("ensureVolume() missing volume %q", name)
		}
		if volume.HostPath == nil {
			t.Fatalf("volume %q has no hostPath", name)
		}
		if volume.HostPath.Path != wantPath {
			t.Errorf("volume %q hostPath = %q, want %q", name, volume.HostPath.Path, wantPath)
		}
	}
}
