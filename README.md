# lxcfs-on-kubernetes

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/lxcfs-on-kubernetes)](https://artifacthub.io/packages/search?repo=lxcfs-on-kubernetes)
[![CodeQL](https://github.com/cndoit18/lxcfs-on-kubernetes/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/cndoit18/lxcfs-on-kubernetes/actions/workflows/codeql-analysis.yml)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fcndoit18%2Flxcfs-on-kubernetes.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fcndoit18%2Flxcfs-on-kubernetes?ref=badge_large)


This project will automatically deploy [LXCFS](https://github.com/lxc/lxcfs) while mounted to the container

## Introduction

[LXCFS](https://github.com/lxc/lxcfs) is a small FUSE filesystem written with the intention of making Linux
containers feel more like a virtual machine. It started as a side-project of
`LXC` but is useable by any runtime.

[LXCFS](https://github.com/lxc/lxcfs) will take care that the information provided by crucial files in `procfs`
such as:

```
/proc/cpuinfo
/proc/diskstats
/proc/meminfo
/proc/stat
/proc/swaps

/proc/uptime
/proc/slabinfo
/sys/devices/system/cpu
/sys/devices/system/cpu/online
```

are container aware such that the values displayed (e.g. in `/proc/uptime`)
really reflect how long the container is running and not how long the host is
running.

## Prerequisites

1. `Kubernetes` cluster (v1.19+) is running. For local development purpose, check [Kind installation](https://kind.sigs.k8s.io/docs/user/quick-start/#installation).
1. `cert-manager` (v1.2+) is [installed](https://cert-manager.io/docs/installation/kubernetes/).
1. `helm` v3 is [installed](https://helm.sh/docs/intro/install/).

## Deploy

Run the helm command to install the lxcfs-on-kubernetes to your cluster:

```
helm repo add lxcfs-on-kubernetes https://cndoit18.github.io/lxcfs-on-kubernetes/
```

you can then do

```
helm upgrade --install lxcfs lxcfs-on-kubernetes/lxcfs-on-kubernetes -n lxcfs --create-namespace
```

For what settings you can override with `--set`, `--set-string`, `--set-file` or `--values`, you can refer to the [values.yaml](charts/lxcfs-on-kubernetes/README.md) file.

you can enable the namespace for injection.

```
kubectl label namespace default mount-lxcfs=enabled
```

> You can change it by setting [matchLabels](charts/lxcfs-on-kubernetes/README.md) during installation


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fcndoit18%2Flxcfs-on-kubernetes.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fcndoit18%2Flxcfs-on-kubernetes?ref=badge_large)