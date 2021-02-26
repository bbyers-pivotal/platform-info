package helpers

import (
	wavefront "github.com/wavefronthq/wavefront-sdk-go/senders"
	"platform-info/structs"
	"strings"
	"time"
)

func SendK8sCPUDataToProxy(sender wavefront.Sender, vcpus int, cores int, environment string) {
	sender.SendMetric("vmw.k8s.vcpus", float64(vcpus), time.Now().Unix(), "platform-info", map[string]string{"env": environment})
	sender.SendMetric("vmw.k8s.cores", float64(cores), time.Now().Unix(), "platform-info", map[string]string{"env": environment})
}

func SendClusterDataToProxy(sender wavefront.Sender, clusters []structs.ClusterList, environment string) {
	sender.SendMetric("vmw.k8s.clusters", float64(len(clusters)), time.Now().Unix(), "platform-info", map[string]string{"env": environment})

	for _, cluster := range clusters {
		sender.SendMetric("vmw.k8s.cluster.vms", float64(len(cluster.VMs)), time.Now().Unix(), "platform-info", map[string]string{"env": environment, "deployment": cluster.Deployment})
		masterVMs := 0
		workerVMs := 0
		for _, vm := range cluster.VMs {
			if strings.Contains(vm.Instance,"master") {
				masterVMs += 1
			}

			if strings.Contains(vm.Instance,"worker") {
				workerVMs += 1
			}
		}
		sender.SendMetric("vmw.k8s.cluster.masters", float64(masterVMs), time.Now().Unix(), "platform-info", map[string]string{"env": environment, "deployment": cluster.Deployment, "name": cluster.Name, "k8sVersion": cluster.K8sVersion, "pksVersion": cluster.PKSVersion})
		sender.SendMetric("vmw.k8s.cluster.workers", float64(workerVMs), time.Now().Unix(), "platform-info", map[string]string{"env": environment, "deployment": cluster.Deployment, "name": cluster.Name, "k8sVersion": cluster.K8sVersion, "pksVersion": cluster.PKSVersion})
	}
}

func SendServiceDataToProxy(sender wavefront.Sender, clusters []structs.ServiceInstance, environment string) {
	sender.SendMetric("vmw.k8s.clusters", float64(len(clusters)), time.Now().Unix(), "platform-info", map[string]string{"env": environment})

	for _, cluster := range clusters {
		sender.SendMetric("vmw.k8s.cluster.vms", float64(len(cluster.VMs)), time.Now().Unix(), "platform-info", map[string]string{"env": environment, "deployment": cluster.Deployment})
		masterVMs := 0
		workerVMs := 0
		for _, vm := range cluster.VMs {
			if strings.Contains(vm.Instance,"master") {
				masterVMs += 1
			}

			if strings.Contains(vm.Instance,"worker") {
				workerVMs += 1
			}
		}
		sender.SendMetric("vmw.k8s.cluster.masters", float64(masterVMs), time.Now().Unix(), "platform-info", map[string]string{"env": environment, "deployment": cluster.Deployment})
		sender.SendMetric("vmw.k8s.cluster.workers", float64(workerVMs), time.Now().Unix(), "platform-info", map[string]string{"env": environment, "deployment": cluster.Deployment})
	}
}