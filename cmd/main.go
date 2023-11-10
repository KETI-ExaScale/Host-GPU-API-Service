package main

import (
	"host-api-service/pkg/api"
	"host-api-service/pkg/traveler"
)

func main() {
	traveler.NodeGPUInfo_ = traveler.NewNodeInfo()
	api.StartServer()
}
