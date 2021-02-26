package structs

type VMInfo struct {
	Instance string `json:"instance"`
	IP string `json:"ips"`
	CID string `json:"vm_cid"`
	VCPUs int `json:"vcpus"`
}
type ClusterList struct {
	Deployment string
	PKSVersion string
	K8sVersion string
	Name string
	VMs []VMInfo
}

type ServiceInstance struct {
	Deployment string
	VMs []VMInfo
}