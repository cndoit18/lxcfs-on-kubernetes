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
	"context"
	"encoding/json"
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
	logr "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	k8sadmission "sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	log                      = logr.Log.WithName("mutate")
	_   k8sadmission.Handler = &mutate{}
)

type admissionOption func(m *mutate)

func WithMutatePath(mutatePath string) admissionOption {
	return func(m *mutate) {
		if !strings.HasSuffix(mutatePath, "/") {
			mutatePath = mutatePath + "/"
		}
		m.mutatePath = mutatePath
	}
}
func WithMutateDecoder(decoder k8sadmission.Decoder) admissionOption {
	return func(m *mutate) {
		m.decoder = decoder
	}
}

// nolint: lll
// +kubebuilder:webhook:path=/mount-lxcfs,mutating=true,failurePolicy=ignore,groups="",resources=pods,verbs=create,versions=v1,name=club.cndoit18.lxcfs-on-kubernetes,sideEffects=NoneOnDryRun,admissionReviewVersions=v1
// +kubebuilder:rbac:groups=core,resources=pods,verbs=*

func AddToManager(mgr manager.Manager, opts ...admissionOption) error {
	m := &mutate{}
	for _, opt := range opts {
		opt(m)
	}

	mgr.GetWebhookServer().Register("/mount-lxcfs", &k8sadmission.Webhook{
		Handler: m,
	})

	return nil
}

type mutate struct {
	decoder    k8sadmission.Decoder
	mutatePath string
}

func (m *mutate) Handle(ctx context.Context, req k8sadmission.Request) k8sadmission.Response {
	pod := &corev1.Pod{}
	err := m.decoder.DecodeRaw(req.Object, pod)
	if err != nil {
		return k8sadmission.Errored(http.StatusInternalServerError, err)
	}

	// TODO: Add filter conditions.
	pod.Spec = m.ensurePodSpec(pod.Spec)

	marshalled, err := json.Marshal(pod)
	if err != nil {
		return k8sadmission.Errored(http.StatusInternalServerError, err)
	}

	// Create the patch
	return k8sadmission.PatchResponseFromRaw(req.Object.Raw, marshalled)
}

func (m *mutate) ensurePodSpec(spec corev1.PodSpec) corev1.PodSpec {
	spec.Volumes = m.ensureVolume(spec.Volumes)
	spec.Containers = m.ensureContainer(spec.Containers)

	return spec
}

func (m *mutate) ensureContainer(cs []corev1.Container) []corev1.Container {
	containers := make([]corev1.Container, 0, len(cs))
	containers = append(containers, cs...)
	for i := range containers {
		containers[i].VolumeMounts = m.ensureVolumeMount(containers[i].VolumeMounts)
	}

	return containers
}

func (m *mutate) ensureVolumeMount(volumeMounts []corev1.VolumeMount) []corev1.VolumeMount {
	mounts := map[string]string{
		"lxcfs-proc-cpuinfo":           "/proc/cpuinfo",
		"lxcfs-proc-diskstats":         "/proc/diskstats",
		"lxcfs-proc-meminfo":           "/proc/meminfo",
		"lxcfs-proc-stat":              "/proc/stat",
		"lxcfs-proc-swaps":             "/proc/swaps",
		"lxcfs-proc-uptime":            "/proc/uptime",
		"lxcfs-proc-loadavg":           "/proc/loadavg",
		"lxcfs-sys-devices-system-cpu": "/sys/devices/system/cpu",
		"lxcfs-root-parent-dir":        filepath.Dir(strings.TrimRight(m.mutatePath, "/")),
	}
	for _, v := range volumeMounts {
		if _, ok := mounts[v.Name]; !ok {
			continue
		}
		delete(mounts, v.Name)
	}
	mountPropagationMode := func(mountName string) *corev1.MountPropagationMode {
		config := map[string]*corev1.MountPropagationMode{
			"lxcfs-root-parent-dir": ptr.To(corev1.MountPropagationHostToContainer),
		}
		if _, ok := config[mountName]; ok {
			return config[mountName]
		}
		return ptr.To(corev1.MountPropagationNone)
	}

	result := make([]corev1.VolumeMount, 0, len(volumeMounts)+len(mounts))
	result = append(result, volumeMounts...)
	for k, v := range mounts {
		result = append(result,
			corev1.VolumeMount{
				Name:             k,
				MountPath:        v,
				ReadOnly:         true,
				MountPropagation: mountPropagationMode(k),
			})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

func (m *mutate) ensureVolume(vs []corev1.Volume) []corev1.Volume {
	volumes := map[string]corev1.VolumeSource{
		"lxcfs-proc-cpuinfo": {
			HostPath: &corev1.HostPathVolumeSource{
				Path: m.mutatePath + "proc/cpuinfo",
			},
		},
		"lxcfs-proc-diskstats": {
			HostPath: &corev1.HostPathVolumeSource{
				Path: m.mutatePath + "proc/diskstats",
			},
		},
		"lxcfs-proc-meminfo": {
			HostPath: &corev1.HostPathVolumeSource{
				Path: m.mutatePath + "proc/meminfo",
			},
		},
		"lxcfs-proc-stat": {
			HostPath: &corev1.HostPathVolumeSource{
				Path: m.mutatePath + "proc/stat",
			},
		},
		"lxcfs-proc-swaps": {
			HostPath: &corev1.HostPathVolumeSource{
				Path: m.mutatePath + "proc/swaps",
			},
		},
		"lxcfs-proc-uptime": {
			HostPath: &corev1.HostPathVolumeSource{
				Path: m.mutatePath + "proc/uptime",
			},
		},
		"lxcfs-proc-loadavg": {
			HostPath: &corev1.HostPathVolumeSource{
				Path: m.mutatePath + "proc/loadavg",
			},
		},
		"lxcfs-sys-devices-system-cpu": {
			HostPath: &corev1.HostPathVolumeSource{
				Path: m.mutatePath + "sys/devices/system/cpu",
			},
		},
		"lxcfs-root-parent-dir": {
			HostPath: &corev1.HostPathVolumeSource{
				Path: filepath.Dir(strings.TrimRight(m.mutatePath, "/")),
			},
		},
	}

	for _, v := range vs {
		if _, ok := volumes[v.Name]; !ok {
			continue
		}
		delete(volumes, v.Name)
	}

	result := make([]corev1.Volume, 0, len(vs)+len(volumes))
	result = append(result, vs...)
	for k, v := range volumes {
		result = append(result, corev1.Volume{
			Name:         k,
			VolumeSource: v,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}
