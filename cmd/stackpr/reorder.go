package stackpr

import (
	"fmt"
	"os"

	"test_cli/pkg/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// reorderCmd represents the reorder command
var reorderCmd = &cobra.Command{
	Use:   "reorder <branch1> <branch2> ...",
	Short: "Reorder branches in the stack",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		currentBranches := viper.GetStringSlice("branches")
		if len(args) != len(currentBranches) {
			fmt.Fprintf(os.Stderr, "Error: Must provide all %d branches in new order, got %d\n", len(currentBranches), len(args))
			os.Exit(1)
		}

		// Validate that provided branches match current branches
		provided := make(map[string]bool)
		for _, branch := range args {
			provided[branch] = true
		}
		current := make(map[string]bool)
		for _, branch := range currentBranches {
			current[branch] = true
		}
		for branch := range provided {
			if !current[branch] {
				fmt.Fprintf(os.Stderr, "Error: Branch %s is not in the current stack\n", branch)
				os.Exit(1)
			}
		}
		for branch := range current {
			if !provided[branch] {
				fmt.Fprintf(os.Stderr, "Error: Branch %s missing from new order\n", branch)
				os.Exit(1)
			}
		}

		// Update branch order
		viper.Set("branches", args)
		if err := config.SaveConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Reordered branches:")
		for i, branch := range args {
			fmt.Printf("%d: %s\n", i+1, branch)
		}
	},
}
