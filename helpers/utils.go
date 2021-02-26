package helpers

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func GetFlagEnvironmentString(cmd *cobra.Command, flag string, env string, message string, required bool) string {
	value := cmd.Flag(flag).Value.String()
	if value == "" {
		value = viper.GetString(env)

		if required {
			if value == "" {
				Bail(message, nil)
			}
		}
		return value
	}
	return value
}

func Bail(message string, err error) {
	if err == nil {
		fmt.Println(message)
	} else {
		fmt.Println(message, err)
	}
	os.Exit(1)
}