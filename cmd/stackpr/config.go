package stackpr

import (
	"fmt"
	"os"

	"test_cli/pkg/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config <key> [value]",
	Short: "View or modify configuration",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		if len(args) == 1 {
			// View config
			fmt.Printf("%s: %v\n", key, viper.Get(key))
			return
		}

		// Set config
		value := args[1]
		if key == "syncMode" && !config.Contains([]string{"rebase", "merge", "reset"}, value) {
			fmt.Fprintf(os.Stderr, "Error: Invalid syncMode '%s'. Must be 'rebase', 'merge', or 'reset'\n", value)
			os.Exit(1)
		}
		viper.Set(key, value)
		if err := config.SaveConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Set %s to %s\n", key, value)
	},
}
