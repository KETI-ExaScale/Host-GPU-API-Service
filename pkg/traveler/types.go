package traveler

var MaxLanes = 16 //nvlink를 확인하기 위해 도는 레인 수

var NodeGPUInfo_ *NodeGPUInfo

type NodeGPUInfo struct {
	IndexUUIDMap  map[string]int32
	TotalGPUCount int
	IsGPUNode     bool
	NvlinkInfo    []NVLink_
	NvlinkStatus  []NVLinkStatus
}

type NVLink_ struct {
	GPU1UUID  string
	GPU2UUID  string
	LaneCount int
}

type NVLinkStatus struct {
	UUID          string
	BusID         string
	Lanes         map[string]int
	P2PUUID       []string
	P2PDeviceType []int //0 GPU, 1 IBMNPU, 2 SWITCH, 255 = UNKNOWN
	P2PBusID      []string
}
