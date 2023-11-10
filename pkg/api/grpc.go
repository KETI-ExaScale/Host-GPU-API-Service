package api

import (
	"context"
	"fmt"
	"host-api-service/pkg/traveler"
	"net"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	grpc "google.golang.org/grpc"
	"k8s.io/klog"
)

type TravelServer struct {
	UnsafeTravelerServer
}

func (*TravelServer) Node(ctx context.Context, req *NodeRequest) (*NodeResponse, error) {
	nodeName := os.Getenv("HOSTNAME")
	clusterName := string(req.GetClusterName())
	gpuCount := int32(traveler.NodeGPUInfo_.TotalGPUCount)

	responseStruct := &NodeResponse{
		ClusterName: clusterName,
		NodeName:    nodeName,
		GPU:         gpuCount,
	}

	return responseStruct, nil
}

func (*TravelServer) Delete(ctx context.Context, req *DockerRequest) (*DockerResponse, error) {
	containerID := string(req.GetDockerid())
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		klog.Errorln(err)
	}

	// 컨테이너 삭제
	err = deleteContainer(cli, containerID)
	if err != nil {
		klog.Errorf("Error deleting container: %v\n", err)
	} else {
		klog.Errorln("Container deleted successfully.")
	}

	return &DockerResponse{}, nil
}

func (*TravelServer) NodeGPUInfo(ctx context.Context, req *Request) (*NodeGPUResponse, error) {
	fmt.Println("# Node GPU Info Called!!")
	gpuCount := traveler.NodeGPUInfo_.TotalGPUCount
	indexUUIDMap := traveler.NodeGPUInfo_.IndexUUIDMap
	nvlink := traveler.NodeGPUInfo_.NvlinkInfo

	var nvlinkInfo []*NVLink

	for _, n := range nvlink {
		n := &NVLink{
			Gpu1Uuid:  n.GPU1UUID,
			Gpu2Uuid:  n.GPU2UUID,
			Lanecount: int32(n.LaneCount),
		}
		nvlinkInfo = append(nvlinkInfo, n)
	}

	responseStruct := &NodeGPUResponse{
		TotalGpuCount: int32(gpuCount),
		IndexUuidMap:  indexUUIDMap,
		NvlinkInfo:    nvlinkInfo,
	}

	return responseStruct, nil
}

func StartServer() {
	lis, err := net.Listen("tcp", ":9321")
	if err != nil {
		klog.Fatalf("failed to listen: %v", err)
	}
	nodeServer := grpc.NewServer()
	RegisterTravelerServer(nodeServer, &TravelServer{})
	fmt.Println("HostAPI Server Started...")

	initHostAPIService()
	traveler.DumpNodeGPUInfo()

	if err := nodeServer.Serve(lis); err != nil {
		klog.Fatalf("failed to serve: %v", err)
	}
}

func initHostAPIService() {
	flag := traveler.IsGPUNode()
	if flag {
		traveler.GetGPUs()
		traveler.GetGPUUID()
		traveler.GetNvlinkInfo()
	}
}

func deleteContainer(cli *client.Client, containerID string) error {
	err := stopContainer(cli, containerID)
	if err != nil {
		return err
	}
	ctx := context.Background()
	options := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         true,
	}

	err = cli.ContainerRemove(ctx, containerID, options)
	if err != nil {
		return err
	}

	return nil
}

func stopContainer(cli *client.Client, containerID string) error {
	ctx := context.Background()

	err := cli.ContainerStop(ctx, containerID, container.StopOptions{})
	if err != nil {
		return err
	}

	return nil
}
