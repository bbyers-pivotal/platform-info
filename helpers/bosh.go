package helpers

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
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

type deploymentsResponse struct {
	Tables []struct {
		Rows []struct {
			Name string `json:"name"`
			Service string `json:"team_s"`
		} `json:"Rows"`
	} `json:"Tables"`
}

type VMList []struct {
	Instance string `json:"instance"`
	IP string `json:"ips"`
	CID string `json:"vm_cid"`
}

type DeploymentList []struct {
	Name string `json:"name"`
}

func BOSHVMs (api string, username string, password string, ca_cert string, deployment string) VMList {
	buf := bytes.NewBuffer([]byte{})

	cmd := exec.Command("bosh", "--client="+ username, "--client-secret="+ password, "--ca-cert="+ ca_cert, "--environment="+api, "-d", deployment, "vms", "--json")
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		Bail("error running bosh vms", err)
	}

	response := vmsResponse{}

	if err := json.NewDecoder(buf).Decode(&response); err != nil {
		Bail("Error decoding clusters json", err)
	}

	info := VMList{}

	for _, r := range response.Tables {
		info = r.Rows
	}

	return info
}

func ServiceDeployments (api string, username string, password string, ca_cert string, serviceName string) []string {
	buf := bytes.NewBuffer([]byte{})

	cmd := exec.Command("bosh", "--client="+ username, "--client-secret="+ password, "--ca-cert="+ ca_cert, "--environment="+api, "deployments", "--json")
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		Bail("error running bosh deployments", err)
	}

	response := deploymentsResponse{}

	if err := json.NewDecoder(buf).Decode(&response); err != nil {
		Bail("Error decoding deployments json", err)
	}

	deployments := []string{}

	for _, r := range response.Tables {
		for _, deployment := range r.Rows {
			if strings.Contains(deployment.Service, serviceName) {
				deployments = append(deployments, deployment.Name)
			}
		}
	}

	return deployments
}