# Compare KubeletConfig with Rendered Machine Config
You can use this script to check if the KubeletConfig has been rendered into MachineConfig

# Prerequisite

Install `go`, you can following the guide here: https://go.dev/doc/install

You also need to install OpenShift CLI (`oc`) and login onto the cluster.

# Usage

```bash
git clone https://github.com/Vincent056/compareKC.git
cd compareKC
./compareKC.sh [kubeletConfig name] [machineConfig name]
```

Examples:

```bash
./compareKC.sh compliance-operator-kubelet-worker 99-master-generated-kubelet
```
