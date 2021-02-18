package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"platform-info/helpers"
)

type VMInfo struct {
	Instance string `json:"instance"`
	IP string `json:"ips"`
	CID string `json:"vm_cid"`
	vCPUs int `json:"vcpus"`
}
type ClusterList struct {
	deployment string
	vms []VMInfo
}

const (
	vCpuToCoreRatio = 2
)

var tkgiCmd = &cobra.Command{
	Use:   "tkgi",
	Short: "Grab TKGI cluster info",
	Long: `Grabs TKGI cluster info`,
	Run: func(cmd *cobra.Command, args []string) {

		//fmt.Println(time.Now())
		viper.SetEnvPrefix("PI")
		viper.AutomaticEnv() // read in environment variables that match

		tkgiApi := helpers.GetFlagEnvironmentString(cmd, "tkgi-api", "tkgi_api", "Missing TKGI API")
		tkgiUsername := helpers.GetFlagEnvironmentString(cmd, "tkgi-username", "tkgi_username", "Missing TKGI Username")
		tkgiPassword := helpers.GetFlagEnvironmentString(cmd, "tkgi-password", "tkgi_password", "Missing TKGI Password")

		boshApi := helpers.GetFlagEnvironmentString(cmd, "bosh-api", "bosh_api", "Missing BOSH API")
		boshClient := helpers.GetFlagEnvironmentString(cmd, "bosh-client", "bosh_client", "Missing BOSH Client")
		boshClientSecret := helpers.GetFlagEnvironmentString(cmd, "bosh-client-secret", "bosh_client_secret","Missing BOSH Client Secret")
		boshBoshCACert := helpers.GetFlagEnvironmentString(cmd, "bosh-ca-cert", "bosh_ca_cert", "Missing BOSH CA Cert")

		vcenterUrl := helpers.GetFlagEnvironmentString(cmd, "vcenter-url", "vcenter_url", "Missing vCenter URL")
		vcenterUsername := helpers.GetFlagEnvironmentString(cmd, "vcenter-username", "vcenter_username", "Missing vCenter Username")
		vcenterPassword := helpers.GetFlagEnvironmentString(cmd, "vcenter-password", "vcenter_password", "Missing vCenter Password")

		fmt.Println("Logging into TKGI CLI")
		helpers.TKGILogin(tkgiApi, tkgiUsername, tkgiPassword)
		fmt.Println("Getting TKGI clusters")
		clusterList := helpers.TKGIClusters()

		vmList := []ClusterList{}

		for _, c := range clusterList {
			deploymentName := "service-instance_"+c.UUID
			fmt.Println("Getting BOSH VMs for", deploymentName)
			vms := helpers.BOSHVMs(boshApi, boshClient, boshClientSecret, boshBoshCACert, deploymentName)

			tempVms := []VMInfo{}
			for _, vm := range vms {
				var vmInfo VMInfo
				vmInfo.Instance = vm.Instance
				vmInfo.IP = vm.IP
				vmInfo.CID = vm.CID
				tempVms = append(tempVms, vmInfo)
			}
			vmList = append(vmList, ClusterList{deploymentName,tempVms})
		}

		for i, cluster := range vmList {
			for j, vm := range cluster.vms {
				fmt.Println("Getting VM info for", vm.IP)
				vmInfo := helpers.GetVMInfo(vcenterUsername, vcenterPassword, vcenterUrl, vm.IP)
				vmList[i].vms[j].vCPUs = vmInfo.CPUs
			}
		}

		vcpus := 0
		for _, cluster := range vmList {
			fmt.Println("BOSH deployment", cluster.deployment)
			for _, vm := range cluster.vms {
				fmt.Println("VM instance:", vm.Instance, "vCPUs:", vm.vCPUs)
				vcpus += vm.vCPUs
			}
		}

		fmt.Println("TKGI vCpus:", vcpus)
		fmt.Println("TKGI Cores:", vcpus / vCpuToCoreRatio)

		//fmt.Println(time.Now())
	},
}

func init() {
	rootCmd.AddCommand(tkgiCmd)

	tkgiCmd.Flags().StringP("tkgi-api", "", "","TKGI (PKS) API [$PI_TKGI_API]")
	tkgiCmd.Flags().StringP("tkgi-username", "", "",  "TKGI admin user [$PI_TKGI_USERNAME]")
	tkgiCmd.Flags().StringP("tkgi-password", "", "", "TKGI admin password [$PI_TKGI_PASSWORD]")

	tkgiCmd.Flags().StringP("bosh-api", "", "","BOSH API [$PI_BOSH_API]")
	tkgiCmd.Flags().StringP("bosh-client", "", "",  "BOSH client [$PI_BOSH_CLIENT]")
	tkgiCmd.Flags().StringP("bosh-client-secret", "", "", "BOSH secret [$PI_BOSH_CLIENT_SECRET]")
	tkgiCmd.Flags().StringP("bosh-ca-cert", "", "", "BOSH CA Cert [$PI_BOSH_CA_CERT]")

	tkgiCmd.Flags().StringP("vcenter-url", "", "","vCenter URL [$PI_VCENTER_URL]")
	tkgiCmd.Flags().StringP("vcenter-username", "", "",  "vCenter admin user [$PI_VCENTER_USERNAME]")
	tkgiCmd.Flags().StringP("vcenter-password", "", "", "vCenter admin password [$PI_VCENTER_PASSWORD]")
}
