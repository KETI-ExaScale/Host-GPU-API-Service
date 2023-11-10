rm -rf /usr/lib/systemd/system/host-api-service.service
rm -rf /usr/bin/local/.keti
cp --force build/host-api-service.service /usr/lib/systemd/system/host-api-service.service
# mkdir /usr/bin/local
mkdir /usr/bin/local/.keti
cp --force build/bin/host-api-service /usr/bin/local/.keti/host-api-service