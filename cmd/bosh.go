package cmd

import (
	"github.com/spf13/cobra"
)

var boshCmd = &cobra.Command{
	Use:   "bosh",
	Short: "Grab BOSH related VM info",
	Long: `Grab BOSH related VM info`,
}

func init() {
	rootCmd.AddCommand(boshCmd)
}
