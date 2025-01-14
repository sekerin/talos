#!/bin/bash

set -e

if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

function main {
  case "$1" in
    "up") up;;
    "down") down;;
    "workspace") workspace;;
    *)
      usage
      exit 2
      ;;
  esac
}

function usage {
  echo "USAGE: ${0##*/} <command>"
  echo "Commands:"
  echo -e "\up\t\spin up QEMU/KVM nodes on the talos0 bridge"
  echo -e "\down\t\tear down the QEMU/KVM nodes"
  echo -e "\workspace\t\run and enter a docker container ready for osctl and kubectl use"
}

NODES=(control-plane-1 control-plane-2 control-plane-3 worker-1)

VM_MEMORY=${VM_MEMORY:-2048}
VM_DISK=${VM_DISK:-10}

COMMON_VIRT_OPTS="--memory=${VM_MEMORY} --cpu=host --vcpus=1 --disk pool=default,size=${VM_DISK} --os-type=linux --os-variant=generic --noautoconsole --graphics none --events on_poweroff=preserve --rng /dev/urandom"

CONTROL_PLANE_1_NAME=control-plane-1
CONTROL_PLANE_1_MAC=52:54:00:a1:9c:ae

CONTROL_PLANE_2_NAME=control-plane-2
CONTROL_PLANE_2_MAC=52:54:00:b2:2f:86

CONTROL_PLANE_3_NAME=control-plane-3
CONTROL_PLANE_3_MAC=52:54:00:c3:61:77

WORKER_1_NAME=worker-1
WORKER_1_MAC=52:54:00:d7:99:c7

function up {
    cp $PWD/../../../build/initramfs.xz ./matchbox/assets/
    cp $PWD/../../../build/vmlinuz ./matchbox/assets/
    cd matchbox/assets
    $PWD/../../../../..//build/osctl-linux-amd64 config generate --install-image docker.io/autonomy/installer:latest integration-test https://kubernetes.talos.dev:6443
    yq w -i init.yaml machine.install.extraKernelArgs[+] 'console=ttyS0'
    yq w -i controlplane.yaml machine.install.extraKernelArgs[+] 'console=ttyS0'
    yq w -i join.yaml machine.install.extraKernelArgs[+] 'console=ttyS0'
    cd -
    virt-install --name $CONTROL_PLANE_1_NAME --network=bridge:talos0,model=e1000,mac=$CONTROL_PLANE_1_MAC $COMMON_VIRT_OPTS --boot=hd,network
    virt-install --name $CONTROL_PLANE_2_NAME --network=bridge:talos0,model=e1000,mac=$CONTROL_PLANE_2_MAC $COMMON_VIRT_OPTS --boot=hd,network
    virt-install --name $CONTROL_PLANE_3_NAME --network=bridge:talos0,model=e1000,mac=$CONTROL_PLANE_3_MAC $COMMON_VIRT_OPTS --boot=hd,network
    virt-install --name $WORKER_1_NAME        --network=bridge:talos0,model=e1000,mac=$WORKER_1_MAC        $COMMON_VIRT_OPTS --boot=hd,network
}

function down {
    for node in ${NODES[@]}; do
      virsh destroy $node
    done
    for node in ${NODES[@]}; do
      virsh undefine $node
    done
    virsh pool-refresh default
    for node in ${NODES[@]}; do
      virsh vol-delete --pool default $node.qcow2
    done
}

function workspace {
  docker run --rm -it -v $PWD:/workspace -v $PWD/../../../build/osctl-linux-amd64:/bin/osctl:ro --network talos --dns 172.28.1.1 -w /workspace/matchbox/assets -e TALOSCONFIG='/workspace/matchbox/assets/talosconfig' -e KUBECONFIG='/workspace/matchbox/assets/kubeconfig' k8s.gcr.io/hyperkube:v1.16.2 bash
}

main $@
