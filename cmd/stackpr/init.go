package stackpr

import (
	"fmt"
	"os"

	"test_cli/pkg/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize .stackpr.yaml configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("remote", "origin")
		viper.Set("syncMode", "rebase")
		if err := config.SaveConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Initialized .stackpr.yaml")

		// Set the base branch correctly in the config file
		fmt.Println("Setting base branch to main in config...")
		configCmd.Run(cmd, []string{"base", "main"})
	},
}
