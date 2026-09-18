#!/usr/bin/env sh
set -x
if [ -z "$1" ]; then
  echo "Usage: $0 <lxcfs-mount-path>" >&2
  exit 1
fi
LXCFS="$1"
containers=$(crictl ps | grep -v pause | grep -v calico | grep -v cilium | awk '{print $1}' | grep -v CONTAINER)
for container in $containers; do
    # Get the container's lxcfs mounts, one "<container_path> -> <host_path>" per line.
    mounts=$(crictl inspect -o go-template --template='{{range .info.config.mounts}}{{.container_path}} -> {{.host_path}}{{println}}{{end}}' "$container" | grep "$LXCFS/")

    echo "Mounts for container $container:"
    echo "$mounts"

    if [ -n "$mounts" ]; then
        echo "remount $container"
        PID=$(crictl inspect --output go-template --template '{{- .info.pid -}}' "$container")
        # Re-bind exactly the lxcfs mounts the container already has, so the
        # healed set always follows the webhook's --lxcfs-proc-files configuration.
        echo "$mounts" | while read -r container_path arrow host_path; do
            if [ "$arrow" != "->" ]; then
                continue
            fi
            echo nsenter --target "$PID" --mount -- mount -o bind "$host_path" "$container_path"
            nsenter --target "$PID" --mount -- mount -o bind "$host_path" "$container_path"
        done
    else
        echo "No LXCFS mount found for container $container"
    fi
done
exit 0
