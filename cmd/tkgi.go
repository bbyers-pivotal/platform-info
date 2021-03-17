package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	wavefront "github.com/wavefronthq/wavefront-sdk-go/senders"
	"os"
	"platform-info/helpers"
	"platform-info/structs"
	"strings"
)

var tkgiCmd = &cobra.Command{
	Use:   "tkgi",
	Short: "Grab TKGI cluster info",
	Long: `Grabs TKGI cluster info`,
	Run: func(cmd *cobra.Command, args []string) {

		//fmt.Println(time.Now())
		viper.SetEnvPrefix("PI")
		viper.AutomaticEnv() // read in environment variables that match

		tkgiApi := helpers.GetFlagEnvironmentString(cmd, "tkgi-api", "tkgi_api", "Missing TKGI API", true)
		tkgiUsername := helpers.GetFlagEnvironmentString(cmd, "tkgi-username", "tkgi_username", "Missing TKGI Username", true)
		tkgiPassword := helpers.GetFlagEnvironmentString(cmd, "tkgi-password", "tkgi_password", "Missing TKGI Password", true)

		boshApi := helpers.GetFlagEnvironmentString(cmd, "bosh-api", "bosh_api", "Missing BOSH API", true)
		boshClient := helpers.GetFlagEnvironmentString(cmd, "bosh-client", "bosh_client", "Missing BOSH Client", true)
		boshClientSecret := helpers.GetFlagEnvironmentString(cmd, "bosh-client-secret", "bosh_client_secret","Missing BOSH Client Secret", true)
		boshBoshCACert := helpers.GetFlagEnvironmentString(cmd, "bosh-ca-cert", "bosh_ca_cert", "Missing BOSH CA Cert", true)

		vcenterUrl := helpers.GetFlagEnvironmentString(cmd, "vcenter-url", "vcenter_url", "Missing vCenter URL", true)
		vcenterUsername := helpers.GetFlagEnvironmentString(cmd, "vcenter-username", "vcenter_username", "Missing vCenter Username", true)
		vcenterPassword := helpers.GetFlagEnvironmentString(cmd, "vcenter-password", "vcenter_password", "Missing vCenter Password", true)

		wavefrontProxy := helpers.GetFlagEnvironmentString(cmd, "wavefront-proxy", "wavefront_proxy", "", false)
		environment := ""
		if wavefrontProxy != "" {
			environment = helpers.GetFlagEnvironmentString(cmd, "environment", "environment", "Missing Environment name", true)
		}



		fmt.Println("Logging into TKGI CLI")
		helpers.TKGILogin(tkgiApi, tkgiUsername, tkgiPassword)
		fmt.Println("Getting TKGI clusters")
		clusterList := helpers.TKGIClusters()

		vmList := []structs.ClusterList{}

		for _, c := range clusterList {
			deploymentName := "service-instance_"+c.UUID
			fmt.Println("Getting BOSH VMs for", deploymentName)
			vms := helpers.BOSHVMs(boshApi, boshClient, boshClientSecret, boshBoshCACert, deploymentName)

			tempVms := []structs.VMInfo{}
			for _, vm := range vms {
				var vmInfo structs.VMInfo
				vmInfo.Instance = vm.Instance
				vmInfo.IP = vm.IP
				vmInfo.CID = vm.CID
				tempVms = append(tempVms, vmInfo)
			}
			vmList = append(vmList, structs.ClusterList{ Deployment: deploymentName, PKSVersion: c.PKSVersion, K8sVersion: c.K8sVersion, Name: c.Name, VMs: tempVms})
		}

		for i, cluster := range vmList {
			for j, vm := range cluster.VMs {
				fmt.Println("Getting VM info for", vm.IP)
				vmInfo := helpers.GetVMInfo(vcenterUsername, vcenterPassword, vcenterUrl, vm.IP)
				vmList[i].VMs[j].VCPUs = vmInfo.CPUs
			}
		}

		f, err := os.Create("results")
		if err != nil {
			helpers.Bail("error creating output file", err)
		}

		defer f.Close()

		var sb strings.Builder
		sb.WriteString("Environment: " + environment + "\n\n")

		vcpus := 0
		for _, cluster := range vmList {
			sb.WriteString(fmt.Sprintf("BOSH deployment: %s\n", cluster.Deployment))
			for _, vm := range cluster.VMs {
				sb.WriteString(fmt.Sprintf("VM instance: %s - vCPUS: %v\n", vm.Instance, vm.VCPUs))
				vcpus += vm.VCPUs
			}
			sb.WriteString("\n")
		}

		cores := vcpus / vCpuToCoreRatio

		sb.WriteString(fmt.Sprintf("TKGI vCpus: %v\n", vcpus))
		sb.WriteString(fmt.Sprintf("TGKI Cores: %v\n", cores))

		fmt.Println(sb.String())
		f.WriteString(sb.String())
		f.Sync()

		if wavefrontProxy != "" {

			proxyCfg := &wavefront.ProxyConfiguration{
				Host: wavefrontProxy,
				MetricsPort: 2878,
				DistributionPort: 2878,
				TracingPort: 30000,
				FlushIntervalSeconds: 10,
			}

			sender, err := wavefront.NewProxySender(proxyCfg)
			if err != nil {
				helpers.Bail("Error setting up Wavefront connection", err)
			}

			helpers.SendK8sCPUDataToProxy(sender, vcpus, cores, environment)
			helpers.SendClusterDataToProxy(sender, vmList, environment)

			sender.Flush()
			sender.Close()
		}
	},
}

func init() {
	boshCmd.AddCommand(tkgiCmd)

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

	tkgiCmd.Flags().StringP("wavefront-proxy", "", "", "Wavefront Proxy [$PI_WAVEFRONT_PROXY]")
	tkgiCmd.Flags().StringP("environment", "", "", "Environment Name [$PI_ENVIRONMENT]")
}
