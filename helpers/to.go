package helpers

import (
	"fmt"
	"platform-info/structs"
)

func SendDataToProxy(clusters []structs.ClusterList) {
	fmt.Println("in send to proxy")
	for _, cluster := range clusters {
		fmt.Println(cluster)
	}
}
