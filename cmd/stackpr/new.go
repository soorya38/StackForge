package stackpr

import (
	"fmt"
	"os"

	"test_cli/pkg/config"
	"test_cli/pkg/git"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new <branch>",
	Short: "Create a new stacked branch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := git.GetRepo()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		branch := args[0]
		if git.BranchExists(repo, branch) {
			fmt.Printf("Branch %s already exists\n", branch)
		} else {
			// Create new branch
			if err := git.CreateBranch(repo, branch); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating branch: %v\n", err)
				os.Exit(1)
			}

			// Checkout the new branch
			if err := git.CheckoutBranch(repo, branch); err != nil {
				fmt.Fprintf(os.Stderr, "Error checking out %s: %v\n", branch, err)
				os.Exit(1)
			}
		}

		// Set base if not set
		base := viper.GetString("base")
		if base == "" {
			head, err := repo.Head()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting HEAD: %v\n", err)
				os.Exit(1)
			}
			base = head.Name().Short()
			viper.Set("base", base)
		}

		// Update config
		branches := viper.GetStringSlice("branches")
		if !config.Contains(branches, branch) {
			branches = append(branches, branch)
			viper.Set("branches", branches)
			if err := config.SaveConfig(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		}

		fmt.Printf("Created branch %s\n", branch)
	},
}
