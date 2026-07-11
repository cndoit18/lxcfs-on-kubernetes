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
| affinity | object | `{}` | Affinity to add to the controller Pods |
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
| lxcfs.podAnnotations | object | `{}` | Additional annotations to add to the agent Pods |
| lxcfs.podMatchLabels | object | `{"platform.glm.ai/mount-lxcfs":"enabled"}` | For pods that match the labels, the Pods under it will mount lxcfs. |
| lxcfs.resources | object | `{"limits":{"cpu":"500m","memory":"300Mi"},"requests":{"cpu":"300m","memory":"200Mi"}}` | Expects input structure as per specification <https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core> |
| lxcfs.useDaemonset | bool | `true` | Installing lxcfs with daemonset |
| podAnnotations | object | `{}` | Additional annotations to add to the controller Pods |
| pullPolicy | string | `"IfNotPresent"` | The image pull policy. |
| replicas | int | `1` | Number of replicas for the controller |
| resources | object | `{"limits":{"cpu":"500m","memory":"300Mi"},"requests":{"cpu":"300m","memory":"200Mi"}}` | Expects input structure as per specification <https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core> |
| service.port | int | `443` | Expose port for WebHook controller |
| service.type | string | `"ClusterIP"` | Service type to use |

----------------------------------------------
Maintained manually. Update this table when `values.yaml` changes.
