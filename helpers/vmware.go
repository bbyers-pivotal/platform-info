package helpers

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
)

type vmInfoResponse struct {
	VirtualMachines []struct {
		Name string `json:"Name"`
		Config struct {
			Hardware struct {
				CPU int `json:"NumCPU"`
			} `json:"Hardware"`
		} `json:"Config"`
	} `json:"VirtualMachines"`
}

type VMInfo struct {
	Name string `json:"name"`
	CPUs int    `json:"cpu"`
}

func GetVMInfo(username, password, url, ip string) VMInfo {
	buf := bytes.NewBuffer([]byte{})

	cmd := exec.Command("govc", "vm.info", "-vm.ip", ip, "-json")
	cmd.Env = append(cmd.Env, "GOVC_USERNAME="+username, "GOVC_PASSWORD="+password, "GOVC_URL="+url, "GOVC_INSECURE=true")
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		Bail("error running govc", err)
	}

	list := vmInfoResponse{}

	if err := json.NewDecoder(buf).Decode(&list); err != nil {
		Bail("Error decoding vm info json", err)
	}

	if len(list.VirtualMachines) > 1 {
		Bail("Found more than one VM with IP " + ip, nil)
	}

	vm := list.VirtualMachines[0]
	vmInfo := VMInfo{Name: vm.Name, CPUs: vm.Config.Hardware.CPU}

	return vmInfo
}