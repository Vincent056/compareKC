#!/bin/bash

set -e

function urldecode() { : "${*//+/ }"; echo -e "${_//%/\\x}"; }

display_description() {
    echo "Help to compare kubeletconfigs with rendered MachineConfig"
}

print_usage() {
    cmdname=$(basename $0)
    echo "Usage:"
    echo "    $cmdname -h                      Display this help message."
    echo "    $cmdname [kubeletConfig name] [machineConfig name]  Compare KubeletConfig with Rendered MachineConfig"
    exit 0
}

while getopts ":hdpcn:P:" opt; do
    case ${opt} in
        h ) # Display help
            display_description
            print_usage
            ;;
        \? ) 
            print_usage
            ;;
    esac
done

if [ -z "$1" ]; then
    echo "Error: KubeletConfig name not provided"
    print_usage
fi

if [ -z "$2" ]; then
    echo "Error: MachineConfig name not provided"
    print_usage
fi

echo "Remove existing tmp files"
rm -f ./kc.json
rm -f ./render.json


oc get kubeletconfig $1 -o json | jq '.spec.kubeletConfig' > kc.json
urldecode $(oc get mc $2 -o json | jq '.spec.config.storage.files[0].contents.source' | cut -d "," -f2- | rev | cut -c5- | rev) > render.json

grep -q '[^[:space:]]' kc.json || { echo "Error: KubeletConfig not found"; exit 1;}
grep -q '[^[:space:]]' render.json || { echo "Error: Rendered MachineConfig not found"; exit 1;}

echo "Comparing KubeletConfig with Rendered MachineConfig"
go run ./compareHelper.go



