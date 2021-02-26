package helpers

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
)

type clusterList []struct {
	Name string `json:"name"`
	K8sVersion string `json:"k8s_version"`
	PKSVersion string `json:"pks_version"`
	UUID string `json:"uuid"`
}

func TKGILogin (api string, username string, password string) {
	buf := bytes.NewBuffer([]byte{})

	cmd := exec.Command("pks", "login", "-a", api, "-u", username, "-p", password, "-k")
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		Bail("error logging into TKGI foundation " + api, err)
	}
}

func TKGIClusters () clusterList {
	buf := bytes.NewBuffer([]byte{})

	cmd := exec.Command("pks", "clusters", "--json")
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		Bail("error running pks clusters", err)
	}

	list := clusterList{}

	if err := json.NewDecoder(buf).Decode(&list); err != nil {
		Bail("Error decoding clusters json", err)
	}
	return list
}

