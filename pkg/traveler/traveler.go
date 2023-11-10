package traveler

import (
	"fmt"
	"log"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"k8s.io/klog"
)

func NewNodeInfo() *NodeGPUInfo {
	return &NodeGPUInfo{
		IndexUUIDMap:  make(map[string]int32),
		TotalGPUCount: 0,
		IsGPUNode:     false,
		NvlinkInfo:    []NVLink_{},
		NvlinkStatus:  []NVLinkStatus{},
	}
}

func IsGPUNode() bool {
	nvmlReturn := nvml.Init()
	if nvmlReturn != nvml.SUCCESS {
		NodeGPUInfo_.IsGPUNode = false
		return false
	} else {
		NodeGPUInfo_.IsGPUNode = true
		return true
	}
}

func GetGPUs() {
	nvml.Init()
	defer func() {
		nvmlReturn := nvml.Shutdown()
		if nvmlReturn != nvml.SUCCESS {
			//log.Fatalf("Unable to shutdown NVML: %v", ret)
			klog.Infof("Unable to shutdown NVML: %v\n", nvmlReturn)
		}
	}()
	count, nvmlReturn := nvml.DeviceGetCount()
	if nvmlReturn != nvml.SUCCESS {
		//log.Fatalf("Unable to get device count: %v", ret)
		klog.Infof("Unable to get device count: %v", nvmlReturn)
		count = 0
	}

	NodeGPUInfo_.TotalGPUCount = count
}

func GetGPUUID() []string {
	count, err := nvml.DeviceGetCount()
	if err != nvml.SUCCESS {
		return nil
	}

	var UUIDs []string
	nvml.Init()
	for i := 0; i < count; i++ {
		device, ret := nvml.DeviceGetHandleByIndex(i)
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get device at index %d: %v", i, ret)
		}
		uuid, _ := device.GetUUID() //uuid
		NodeGPUInfo_.IndexUUIDMap[uuid] = int32(i)
	}

	return UUIDs
}

func GetNvlinkInfo() {
	count, err := nvml.DeviceGetCount()
	if err != nvml.SUCCESS {
		return
	}

	var nvlinkStatus []NVLinkStatus
	var nvlinkList []NVLink_

	for i := 0; i < count; i++ {
		nvlinkStatus = append(nvlinkStatus, NVLinkStatus{})
		device, _ := nvml.DeviceGetHandleByIndex(i)
		pciinfo, _ := device.GetPciInfo()
		tmparray := pciinfo.BusId

		var bytebus [32]byte
		for j := 0; j < 32; j++ {
			bytebus[j] = byte(tmparray[i])
		}

		nvlinkStatus[i].BusID = string(bytebus[:])
		nvlinkStatus[i].UUID, _ = device.GetUUID()
		nvlinkStatus[i].Lanes = make(map[string]int)

		for j := 0; j < MaxLanes; j++ { //Check nvlink by circling lanes as many as maxlane
			P2PPciInfo, err := device.GetNvLinkRemotePciInfo(j)
			if err != nvml.SUCCESS {
				break
			}
			tmparray := P2PPciInfo.BusId
			var bytebus [32]byte
			for k := 0; k < 32; k++ {
				bytebus[k] = byte(tmparray[k])
			}
			val, exists := nvlinkStatus[i].Lanes[string(bytebus[:])]
			if !exists {
				P2PDevice, err := nvml.DeviceGetHandleByPciBusId(string(bytebus[:]))
				if err != nvml.SUCCESS {
					fmt.Println("error can get device handle")
				} else {
					P2PIndex, _ := P2PDevice.GetIndex()
					if P2PIndex > j {
						types, _ := device.GetNvLinkRemoteDeviceType(j)
						nvlinkStatus[i].Lanes[string(bytebus[:])] = 1
						nvlinkStatus[i].P2PDeviceType = append(nvlinkStatus[i].P2PDeviceType, int(types))
						P2PUUID, _ := P2PDevice.GetUUID()
						nvlinkStatus[i].P2PUUID = append(nvlinkStatus[i].P2PUUID, P2PUUID)
						nvlinkStatus[i].P2PBusID = append(nvlinkStatus[i].P2PBusID, string(bytebus[:]))
					}
				}

			} else {
				nvlinkStatus[i].Lanes[string(bytebus[:])] = val + 1
			}
		}

		for j := 0; j < len(nvlinkStatus[i].P2PUUID); j++ {
			var nvlink NVLink_
			nvlink.GPU1UUID = nvlinkStatus[i].UUID
			nvlink.GPU2UUID = nvlinkStatus[i].P2PUUID[j]
			nvlink.LaneCount = nvlinkStatus[i].Lanes[nvlinkStatus[i].P2PBusID[j]]
			nvlinkList = append(nvlinkList, nvlink)
		}
	}

	NodeGPUInfo_.NvlinkInfo = nvlinkList
	NodeGPUInfo_.NvlinkStatus = nvlinkStatus
}

func DumpNodeGPUInfo() {
	fmt.Printf("1. Is GPU Node : %v\n", NodeGPUInfo_.IsGPUNode)
	fmt.Printf("2. Total GPU Count : %v\n", NodeGPUInfo_.TotalGPUCount)
	fmt.Printf("3. Node's GPU UUIDs\n")
	for uuid, index := range NodeGPUInfo_.IndexUUIDMap {
		fmt.Printf(" - GPU UUID : %s, Index: %d\n", uuid, index)
	}
	fmt.Printf("4. GPU's NVLink Info\n")
	for _, nvlink := range NodeGPUInfo_.NvlinkInfo {
		fmt.Printf(" - (GPU1 UUID: %s) (GPU2 UUID: %s) (LaneCout: %d)\n", nvlink.GPU1UUID, nvlink.GPU2UUID, nvlink.LaneCount)
	}
}
