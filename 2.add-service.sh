rm -rf /usr/lib/systemd/system/host-api-service.service
rm -rf /usr/bin/local/.keti
cp --force build/host-gpu-api-service.service /usr/lib/systemd/system/host-gpu-api-service.service
# mkdir /usr/bin/local
mkdir /usr/bin/local/.keti
cp --force build/bin/host-gpu-api-service /usr/bin/local/.keti/host-gpu-api-service