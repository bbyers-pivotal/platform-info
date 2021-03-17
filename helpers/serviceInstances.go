package helpers

import (
	"fmt"
	"platform-info/structs"
)

type BoshConnectionInfo struct {
	API string
	Client string
	ClientSecret string
	CACert string
}

type VCenterConnectionInfo struct {
	API string
	Username string
	Password string
}

func GetServiceInstances(bci BoshConnectionInfo, vci VCenterConnectionInfo, serviceName string, includeParent bool) []structs.ServiceInstance {
	serviceInstances := ServiceDeployments(bci.API, bci.Client, bci.ClientSecret, bci.CACert, serviceName, includeParent)

	vmList := []structs.ServiceInstance{}

	for _, deploymentName := range serviceInstances {
		fmt.Println("Getting BOSH VMs for", deploymentName)
		vms := BOSHVMs(bci.API, bci.Client, bci.ClientSecret, bci.CACert, deploymentName)

		tempVms := []structs.VMInfo{}
		for _, vm := range vms {
			var vmInfo structs.VMInfo
			vmInfo.Instance = vm.Instance
			vmInfo.IP = vm.IP
			vmInfo.CID = vm.CID
			tempVms = append(tempVms, vmInfo)
		}
		vmList = append(vmList, structs.ServiceInstance{ Deployment: deploymentName, VMs: tempVms})
	}

	for i, cluster := range vmList {
		for j, vm := range cluster.VMs {
			fmt.Println("Getting VM info for", vm.IP)
			vmInfo := GetVMInfo(vci.Username, vci.Password, vci.API, vm.IP)
			vmList[i].VMs[j].VCPUs = vmInfo.CPUs
		}
	}

	return vmList
}
