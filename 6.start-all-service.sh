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
        ./2.add-service.sh
        ./3.start.sh
    else
        ssh root@${node_ip_mapping[$node]} "cd /home/host-gpu-api-service && bash ./2.add-service.sh"
        ssh root@${node_ip_mapping[$node]} "cd /home/host-gpu-api-service && bash ./3.start.sh"
    fi
done