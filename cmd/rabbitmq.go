package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"platform-info/helpers"
	"platform-info/structs"
	"strings"
)



var rabbitmqCmd = &cobra.Command{
	Use:   "rabbitmq",
	Short: "Grab RabbitMQ cluster info",
	Long: `Grabs RabbitMQ cluster info`,
	Run: func(cmd *cobra.Command, args []string) {

		var serviceNames  = []string{"p-rabbitmq", "p.rabbitmq"}
		//fmt.Println(time.Now())
		viper.SetEnvPrefix("PI")
		viper.AutomaticEnv() // read in environment variables that match

		boshApi := helpers.GetFlagEnvironmentString(cmd, "bosh-api", "bosh_api", "Missing BOSH API", true)
		boshClient := helpers.GetFlagEnvironmentString(cmd, "bosh-client", "bosh_client", "Missing BOSH Client", true)
		boshClientSecret := helpers.GetFlagEnvironmentString(cmd, "bosh-client-secret", "bosh_client_secret","Missing BOSH Client Secret", true)
		boshBoshCACert := helpers.GetFlagEnvironmentString(cmd, "bosh-ca-cert", "bosh_ca_cert", "Missing BOSH CA Cert", true)

		vcenterUrl := helpers.GetFlagEnvironmentString(cmd, "vcenter-url", "vcenter_url", "Missing vCenter URL", true)
		vcenterUsername := helpers.GetFlagEnvironmentString(cmd, "vcenter-username", "vcenter_username", "Missing vCenter Username", true)
		vcenterPassword := helpers.GetFlagEnvironmentString(cmd, "vcenter-password", "vcenter_password", "Missing vCenter Password", true)

		environment := helpers.GetFlagEnvironmentString(cmd, "environment", "environment", "Missing Environment name", true)
		//wavefrontProxy := helpers.GetFlagEnvironmentString(cmd, "wavefront-proxy", "wavefront_proxy", "", false)
		//environment := ""
		//if wavefrontProxy != "" {
		//	environment = helpers.GetFlagEnvironmentString(cmd, "environment", "environment", "Missing Environment name", true)
		//}

		bci := helpers.BoshConnectionInfo{
			API:          boshApi,
			Client:       boshClient,
			ClientSecret: boshClientSecret,
			CACert:       boshBoshCACert,
		}

		vci := helpers.VCenterConnectionInfo{
			API:      vcenterUrl,
			Username: vcenterUsername,
			Password: vcenterPassword,
		}

		var instances []structs.ServiceInstance
		for _, name := range serviceNames {
			sis := helpers.GetServiceInstances(bci, vci, name, true)
			for _, si := range sis {
				instances = append(instances, si)
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
		for _, cluster := range instances {
			sb.WriteString(fmt.Sprintf("BOSH deployment: %s\n", cluster.Deployment))
			for _, vm := range cluster.VMs {
				sb.WriteString(fmt.Sprintf("VM instance: %s - vCPUS: %v\n", vm.Instance, vm.VCPUs))
				vcpus += vm.VCPUs
			}
			sb.WriteString("\n")
		}
		cores := vcpus / vCpuToCoreRatio

		sb.WriteString(fmt.Sprintf("RabbitMQ vCpus: %v\n", vcpus))
		sb.WriteString(fmt.Sprintf("RabbitMQ Cores: %v\n", cores))

		fmt.Println(sb.String())
		f.WriteString(sb.String())
		f.Sync()
		//if wavefrontProxy != "" {
		//
		//	proxyCfg := &wavefront.ProxyConfiguration{
		//		Host: wavefrontProxy,
		//		MetricsPort: 2878,
		//		DistributionPort: 2878,
		//		TracingPort: 30000,
		//		FlushIntervalSeconds: 10,
		//	}
		//
		//	sender, err := wavefront.NewProxySender(proxyCfg)
		//	if err != nil {
		//		helpers.Bail("Error setting up Wavefront connection", err)
		//	}
		//
		//	//helpers.SendK8sCPUDataToProxy(sender, vcpus, cores, environment)
		//	//helpers.SendServiceDataToProxy(sender, vmList, environment)
		//
		//	sender.Flush()
		//	sender.Close()
		//}
		//fmt.Println(time.Now())
	},
}

func init() {
	boshCmd.AddCommand(rabbitmqCmd)

	rabbitmqCmd.Flags().StringP("bosh-api", "", "","BOSH API [$PI_BOSH_API]")
	rabbitmqCmd.Flags().StringP("bosh-client", "", "",  "BOSH client [$PI_BOSH_CLIENT]")
	rabbitmqCmd.Flags().StringP("bosh-client-secret", "", "", "BOSH secret [$PI_BOSH_CLIENT_SECRET]")
	rabbitmqCmd.Flags().StringP("bosh-ca-cert", "", "", "BOSH CA Cert [$PI_BOSH_CA_CERT]")

	rabbitmqCmd.Flags().StringP("vcenter-url", "", "","vCenter URL [$PI_VCENTER_URL]")
	rabbitmqCmd.Flags().StringP("vcenter-username", "", "",  "vCenter admin user [$PI_VCENTER_USERNAME]")
	rabbitmqCmd.Flags().StringP("vcenter-password", "", "", "vCenter admin password [$PI_VCENTER_PASSWORD]")

	rabbitmqCmd.Flags().StringP("wavefront-proxy", "", "", "Wavefront Proxy [$PI_WAVEFRONT_PROXY]")
	rabbitmqCmd.Flags().StringP("environment", "", "", "Environment Name [$PI_ENVIRONMENT]")
}
