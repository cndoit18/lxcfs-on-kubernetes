# lxcfs-on-kubernetes

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.1.0](https://img.shields.io/badge/AppVersion-0.1.0-informational?style=flat-square)

`LXCFS` is a small FUSE filesystem written with the intention of making Linux containers feel more like a virtual machine.

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| cndoit18 | cndoit18@outlook.com | https://github.com/cndoit18 |

## Source Code

* <https://github.com/cndoit18/lxcfs-on-kubernetes>
* <https://github.com/lxc/lxcfs>

## Requirements

Kubernetes: `>= 1.16.0-0`

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Affinity for the controller Pods (podAntiAffinity, nodeAffinity, etc.). May be set together with nodeSelector: the scheduler requires all of them to match (AND). |
| image.agent | string | `"ghcr.io/cndoit18/lxcfs-agent:latest"` | lxcfs-on-kubernetes agent image |
| image.manager | string | `"ghcr.io/cndoit18/lxcfs-manager:latest"` | lxcfs-on-kubernetes controller image |
| imagePullSecrets | list | `[]` | Reference to one or more secrets to be used when pulling images <https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/> For example: `[   {"name":"image-pull-secret"} ]` |
| leaderElection.enable | bool | `false` | Whether to enabled leaderElection |
| leaderElection.id | string | `"lxcfs-on-kubernetes-leader-election"` | The id used to store the ConfigMap for leader election |
| leaderElection.namespace | string | `"kube-system"` | The namespace used to store the ConfigMap for leader election |
| logLevel | int | `4` | Set the verbosity of controller. Range of 0 - 6 with 6 being the most verbose. Info level is 4. |
| lxcfs.args | list | `["-l","--enable-cfs","--enable-pidfd"]` | Adjusting the boot parameters of lxcfs |
| lxcfs.configMaps.crictlConfig.endpoint | string | `"/run/containerd/containerd.sock"` | The endpoint of the CRI runtime |
| lxcfs.mountPath | string | `"/var/lib/lxcfs-on-k8s/lxcfs"` | Specify the mount path of lxcfs on the host |
| lxcfs.namespaceMatchLabels | object | `{}` | For namespaces that match the labels, only pods in matching namespaces will be considered. Empty by default (namespaceSelector disabled). |
| lxcfs.nodeSelector | object | `{}` | Node selector for the lxcfs DaemonSet |
| lxcfs.podAnnotations | object | `{}` | Additional annotations to add to the agent Pods |
| lxcfs.podMatchLabels | object | `{"platform.glm.ai/mount-lxcfs":"enabled"}` | For pods that match the labels, the Pods under it will mount lxcfs. |
| lxcfs.podSecurityContext | object | `{}` | Pod-level security context for the lxcfs DaemonSet |
| lxcfs.procFiles | list | `["cpuinfo","diskstats","meminfo","stat","swaps","uptime","loadavg","pressure","slabinfo"]` | LXCFS proc files to mount into pods by the webhook. Entries are names relative to `/proc` (e.g. `cpuinfo`); an absolute `/proc/...` prefix is also accepted. Older docker/runc releases only allow a subset of `/proc` files to be bind-mounted (`cpuinfo`, `diskstats`, `meminfo`, `stat`, `swaps`, `uptime`); tune this list accordingly if pods crash with `cannot be mounted because it is located inside "/proc"`. Set to an empty list (`[]`) to disable proc mounts entirely. |
| lxcfs.resources | object | `{"limits":{"cpu":"500m","memory":"300Mi"},"requests":{"cpu":"300m","memory":"200Mi"}}` | Expects input structure as per specification <https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core> |
| lxcfs.tolerations | list | `[]` | Tolerations for the lxcfs DaemonSet (e.g. to run on tainted nodes) |
| lxcfs.updateStrategy | object | `{"type":"RollingUpdate","rollingUpdate":{"maxUnavailable":1}}` | DaemonSet update strategy |
| lxcfs.useDaemonset | bool | `true` | Installing lxcfs with daemonset |
| nodeSelector | object | `{}` | Node selector for the controller Deployment |
| podAnnotations | object | `{}` | Additional annotations to add to the controller Pods |
| podSecurityContext | object | `{}` | Pod-level security context for the controller |
| pullPolicy | string | `"IfNotPresent"` | The image pull policy. |
| replicas | int | `1` | Number of replicas for the controller |
| resources | object | `{"limits":{"cpu":"500m","memory":"300Mi"},"requests":{"cpu":"300m","memory":"200Mi"}}` | Expects input structure as per specification <https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core> |
| revisionHistoryLimit | int | `10` | Number of old ReplicaSets to retain for the Deployment |
| service.port | int | `443` | Expose port for WebHook controller |
| service.type | string | `"ClusterIP"` | Service type to use |
| tolerations | list | `[]` | Tolerations for the controller Deployment |

----------------------------------------------
Maintained manually. Update this table when `values.yaml` changes.
