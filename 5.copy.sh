#!/bin/bash
declare -A node_ip_mapping

node_info=$(kubectl get nodes -o custom-columns="NAME:.metadata.name,INTERNAL-IP:.status.addresses[0].address" --no-headers)

while IFS= read -r line; do
    node_name=$(echo "$line" | awk '{print $1}')
    ip_address=$(echo "$line" | awk '{print $2}')
    node_ip_mapping["$node_name"]=$ip_address
done <<< "$node_info"

for node in "${!node_ip_mapping[@]}"; do
    if [ "$node" == "c1-master" ]; then
        continue
    else
        scp -r ./build root@${node_ip_mapping[$node]}:/home/host-gpu-api-service
        scp "2.add-service.sh" root@${node_ip_mapping[$node]}:/home/host-gpu-api-service
        scp "3.start.sh" root@${node_ip_mapping[$node]}:/home/host-gpu-api-service
    fi
done
