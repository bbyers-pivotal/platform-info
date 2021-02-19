package helpers

import (
	"fmt"
	"os"
	wavefront "github.com/wavefronthq/wavefront-sdk-go/senders"
	"platform-info/structs"
	"time"
)

func SendCPUDataToProxy(vcpus int, cores int, proxyAddress string, environment string) {
	proxyCfg := &wavefront.ProxyConfiguration{
		Host: proxyAddress,
		MetricsPort: 2878,
		DistributionPort: 2878,
		TracingPort: 30000,
		FlushIntervalSeconds: 10,
	}

	sender, err := wavefront.NewProxySender(proxyCfg)
	if err != nil {
		fmt.Println("Error setting up Wavefront connection", err)
		os.Exit(1)
	}

	sender.SendMetric("vmw.k8s.vcpus", float64(vcpus), time.Now().Unix(), "platform-info", map[string]string{"env": environment})
	sender.SendMetric("vmw.k8s.cores", float64(cores), time.Now().Unix(), "platform-info", map[string]string{"env": environment})
	sender.Flush()
	sender.Close()
}

func SendClusterDataToProxy(clusters []structs.ClusterList, proxyAddress string, environment string) {

}