#!/usr/bin/env sh
set -x
if [[ -z "$1" ]]; then
  echo "Usage: $0 <lxcfs-mount-path>" >&2
  exit 1
fi
LXCFS="$1"
containers=$(crictl ps | grep -v pause | grep -v calico | grep -v cilium | awk '{print $1}' | grep -v CONTAINER)
for container in $containers; do
    # Get the container's mounts
    mounts=$(crictl inspect -o go-template --template='{{range .info.config.mounts}}{{.container_path}} -> {{.host_path}}{{println}}{{end}}' $container | grep "$LXCFS/")
    
    echo "Mounts for container $container:"
    echo "$mounts"
    
    # Check if the container has the LXCFS mount
    count=$(echo "$mounts" | grep  "$LXCFS/" | wc -l)
    if [ "$count" != "0" ]; then
        echo "remount $container"
        PID=$(crictl inspect --output go-template --template '{{- .info.pid -}}' $container)
        # mount /proc
        for file in meminfo cpuinfo loadavg stat diskstats swaps uptime; do
            echo nsenter --target $PID --mount -- mount -B "$LXCFS/proc/$file" "/proc/$file"
            nsenter --target $PID --mount -- mount -B "$LXCFS/proc/$file" "/proc/$file"
        done

        echo nsenter --target $PID --mount -- mount -B "$LXCFS/sys/devices/system/cpu" "/sys/devices/system/cpu"
        nsenter --target $PID --mount -- mount -B "$LXCFS/sys/devices/system/cpu" "/sys/devices/system/cpu"
    else
        echo "No LXCFS mount found for container $container"
    fi
done
exit 0