# lxcfs-on-kubernetes

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![CodeQL](https://github.com/cndoit18/lxcfs-on-kubernetes/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/cndoit18/lxcfs-on-kubernetes/actions/workflows/codeql-analysis.yml)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fcndoit18%2Flxcfs-on-kubernetes.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fcndoit18%2Flxcfs-on-kubernetes?ref=badge_shield)


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

you can enable the namespace for injection.

```
kubectl label namespace default mount-lxcfs=enabled
```

## Release

Releases are produced exclusively by the [`Release`](.github/workflows/release.yaml)
workflow. It builds the Helm chart package, pushes the `lxcfs-manager` / `lxcfs-agent`
images, creates the GitHub release with the chart `.tgz` attached, and updates the
Helm repository index (`gh-pages`).

Two equivalent ways to release:

1. **Recommended:** run the `Release` workflow (Actions → Release → Run workflow) with
   a `v`-prefixed semver version, e.g. `v0.2.9`.
2. Push an annotated tag (`git tag --annotate v0.2.9 && git push origin v0.2.9`). The
   same workflow triggers on `v*` tag pushes.

Never create a GitHub release by hand: doing so skips attaching the chart asset and
updating the Helm repository index (this broke the `v0.2.8` release). If a release is
incomplete, re-run the `Release` workflow for that version — it is idempotent.

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fcndoit18%2Flxcfs-on-kubernetes.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fcndoit18%2Flxcfs-on-kubernetes?ref=badge_large)
