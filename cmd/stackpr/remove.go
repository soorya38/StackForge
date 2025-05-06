package stackpr

import (
	"fmt"
	"os"

	"test_cli/pkg/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove <branch>",
	Short: "Remove a branch from the stack",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branch := args[0]
		branches := viper.GetStringSlice("branches")
		if !config.Contains(branches, branch) {
			fmt.Fprintf(os.Stderr, "Error: Branch %s is not in the stack\n", branch)
			os.Exit(1)
		}

		newBranches := []string{}
		for _, b := range branches {
			if b != branch {
				newBranches = append(newBranches, b)
			}
		}
		viper.Set("branches", newBranches)
		if err := config.SaveConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Removed %s from stack\n", branch)
	},
}
