package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type vmsResponse struct {
	Tables []struct {
		Rows []struct {
			Instance string `json:"instance"`
			IP string `json:"ips"`
			CID string `json:"vm_cid"`
		} `json:"Rows"`
	} `json:"Tables"`
}

type VMList []struct {
	Instance string `json:"instance"`
	IP string `json:"ips"`
	CID string `json:"vm_cid"`
}

func BOSHVMs (api string, username string, password string, ca_cert string, deployment string) VMList {
	buf := bytes.NewBuffer([]byte{})

	cmd := exec.Command("bosh", "--client="+ username, "--client-secret="+ password, "--ca-cert="+ ca_cert, "--environment="+api, "-d", deployment, "vms", "--json")
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("error running bosh vms", err)
		os.Exit(1)
	}

	response := vmsResponse{}

	if err := json.NewDecoder(buf).Decode(&response); err != nil {
		fmt.Println("Error decoding clusters json", err)
		os.Exit(1)
	}

	info := VMList{}

	for _, r := range response.Tables {
		info = r.Rows
	}

	return info
}