package helpers

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func GetFlagEnvironmentString(cmd *cobra.Command, flag string, env string, message string) string {
	value := cmd.Flag(flag).Value.String()
	if value == "" {
		value = viper.GetString(env)

		if value == "" {
			fmt.Println(message)
			os.Exit(1)
		}

		return value
	}
	return value
}