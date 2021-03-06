package cmd

import (
  "fmt"
  "github.com/spf13/cobra"
  "platform-info/helpers"

  homedir "github.com/mitchellh/go-homedir"
  "github.com/spf13/viper"
)


const (
  vCpuToCoreRatio = 2
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:   "platform-info",
  Short: "Utility for gathering info about your Tanzu environment",
  Long: `Gathers information for TGF and TKGI usage`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    helpers.Bail("error starting app", err)
  }
}

func init() {
  cobra.OnInitialize(initConfig)
  rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.platform-info-config.yaml)")
}


// initConfig reads in config file and ENV variables if set.
func initConfig() {
  if cfgFile != "" {
    // Use config file from the flag.
    viper.SetConfigFile(cfgFile)
  } else {
    // Find home directory.
    home, err := homedir.Dir()
    if err != nil {
      helpers.Bail("error finding home dir", err)
    }

    // Search config in home directory with name ".platform-info-config" (without extension).
    viper.AddConfigPath(home)
    viper.SetConfigName(".platform-info-config")
  }

  // If a config file is found, read it in.
  if err := viper.ReadInConfig(); err == nil {
    fmt.Println("Using config file:", viper.ConfigFileUsed())
  }
}

