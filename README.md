# lxcfs-on-kubernetes

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

1. `Kubernetes` cluster (v1.19+) is running (see [Applicative Kubernetes versions](../../README.md#applicative-kubernetes-versions)
   for more information). For local development purpose, check [Kind installation](./kind-installation.md).
1. `cert-manager` (v1.2+) is [installed](https://cert-manager.io/docs/installation/kubernetes/).
1. `helm` v3 is [installed](https://helm.sh/docs/intro/install/).

## Deploy

Run the helm command to install the lxcfs-on-kubernetes to your cluster:

```
helm repo add lxcfs-on-kubernetes https://cndoit18.github.io/lxcfs-on-kubernetes/
```

you can then do

```
helm upgrade --install lxcfs lxcfs-on-kubernetes/lxcfs-on-kubernetes
```

For what settings you can override with `--set`, `--set-string`, `--set-file` or `--values`, you can refer to the [values.yaml](charts/lxcfs-on-kubernetes/values.yaml) file.