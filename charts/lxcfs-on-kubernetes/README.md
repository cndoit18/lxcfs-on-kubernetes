# lxcfs-on-kubernetes

![Version: 0.2.1](https://img.shields.io/badge/Version-0.2.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.2.1](https://img.shields.io/badge/AppVersion-0.2.1-informational?style=flat-square)

[LXCFS](https://github.com/lxc/lxcfs) is a small FUSE filesystem designed to make Linux containers feel more like virtual machines.

## Maintainers

| Name      | Email                | Url                              |
|-----------|----------------------|----------------------------------|
| cndoit18  | cndoit18@outlook.com | https://github.com/cndoit18      |

## Source Code

* <https://github.com/cndoit18/lxcfs-on-kubernetes>
* <https://github.com/lxc/lxcfs>

## Requirements

Kubernetes: `>= 1.16.0-0`

## Breaking Changes in v0.2.1

### LXCFS Mount Path Validation

Starting from version v0.2.1, the `lxcfs.mountPath` configuration now requires that the parent directory of the mount path must end with **`lxcfs-on-k8s`** (e.g., `/var/lib/lxcfs-on-k8s/lxcfs`).

**Note**: The parent directory can have any name as long as it ends with 'lxcfs-on-k8s' (e.g., `/custom/path/lxcfs-on-k8s/lxcfs`).

Previously, the parent directory could be **`lxc`** (e.g., `/var/lib/lxc/lxcfs`). This change affects:

1. **Existing deployments** that use a custom mount path with parent directory `lxc`
2. **Manual LXCFS installations** on nodes separate from the chart deployment

#### What You Need to Do

- If you have existing deployments using the old path `/var/lib/lxc/lxcfs`, you must update your configuration to use `/var/lib/lxcfs-on-k8s/lxcfs`
- If you manually installed LXCFS on nodes, ensure the mount path follows the new validation rule

#### Example Configuration

```yaml
# Before v0.2.1 (valid)
lxcfs:
  mountPath: /var/lib/lxc/lxcfs

# From v0.2.1 onward (valid)
lxcfs:
  mountPath: /var/lib/lxcfs-on-k8s/lxcfs
```

#### Additional Notes

- This validation is enforced by the `EnsureLxcfsParentDir` function in the controller manager
- The controller will fail to start if the path does not meet this requirement
- Ensure all nodes in your cluster use the same mount path configuration

## Values

| Key                      | Type   | Default                                             | Description |
|--------------------------|--------|-----------------------------------------------------|-------------|
| affinity                 | object | `{}`                                                | Affinity to add to the controller Pods |
| image.agent              | string | `"ghcr.io/cndoit18/lxcfs-agent:v0.2.1"`            | lxcfs-on-kubernetes agent image |
| image.manager            | string | `"ghcr.io/cndoit18/lxcfs-manager:v0.2.1"`          | lxcfs-on-kubernetes controller image |
| imagePullSecrets         | list   | `[]`                                                | Reference to one or more secrets to be used when pulling images. See [Kubernetes documentation](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/) |
| leaderElection.enable    | bool   | `false`                                             | Whether to enable leader election |
| leaderElection.id        | string | `"lxcfs-on-kubernetes-leader-election"`            | The id used to store the ConfigMap for leader election |
| leaderElection.namespace | string | `"kube-system"`                                     | The namespace used to store the ConfigMap for leader election |
| logLevel                 | int    | `4`                                                 | Set the verbosity of controller. Range of 0-6 with 6 being the most verbose. Info level is 4. |
| lxcfs.args               | list   | `["-l","--enable-cfs","--enable-pidfd"]`           | Adjusting the boot parameters of lxcfs |
| lxcfs.configMaps         | object | `{"crictlConfig":{"endpoint":"/run/containerd/containerd.sock"}}` | ConfigMaps to be created for lxcfs |
| lxcfs.matchLabels        | object | `{"mount-lxcfs":"enabled"}`                        | For namespaces that match the labels, the Pods under it will mount lxcfs |
| lxcfs.mountPath          | string | `"/var/lib/lxcfs-on-k8s/lxcfs"`                    | Specify the mount path of lxcfs on the host |
| lxcfs.podAnnotations     | object | `{}`                                                | Additional annotations to add to the agent Pods |
| lxcfs.resources          | object | `{"limits":{"cpu":"500m","memory":"300Mi"},"requests":{"cpu":"300m","memory":"200M"}}` | Expects input structure as per specification. See [Kubernetes API](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| lxcfs.useDaemonset       | bool   | `true`                                              | Installing lxcfs with daemonset |
| podAnnotations           | object | `{}`                                                | Additional annotations to add to the controller Pods |
| pullPolicy               | string | `"IfNotPresent"`                                    | The image pull policy |
| replicas                 | int    | `1`                                                 | Number of replicas for the controller |
| resources                | object | `{"limits":{"cpu":"500m","memory":"300Mi"},"requests":{"cpu":"300m","memory":"200Mi"}}` | Expects input structure as per specification. See [Kubernetes API](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| service.port             | int    | `443`                                               | Expose port for WebHook controller |
| service.type             | string | `"ClusterIP"`                                       | Service type to use |

---

*Autogenerated from chart metadata using [helm-docs v1.4.0](https://github.com/norwoodj/helm-docs/releases/v1.4.0)*
