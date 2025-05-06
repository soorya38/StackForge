package stackpr

import (
	"fmt"
	"os"
	"os/exec"

	"test_cli/pkg/config"
	"test_cli/pkg/git"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge <branch>",
	Short: "Merge a stacked branch into its parent",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := git.GetRepo()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		branch := args[0]
		branches := viper.GetStringSlice("branches")
		if !config.Contains(branches, branch) {
			fmt.Fprintf(os.Stderr, "Error: Branch %s is not in the stack\n", branch)
			os.Exit(1)
		}

		base := viper.GetString("base")
		if base == "" {
			fmt.Fprintf(os.Stderr, "Error: Base branch not set in .stackpr.yaml\n")
			os.Exit(1)
		}

		// Find parent
		parent := base
		for i, b := range branches {
			if b == branch && i > 0 {
				parent = branches[i-1]
			}
		}

		if !git.BranchExists(repo, parent) {
			fmt.Fprintf(os.Stderr, "Error: Parent branch %s does not exist\n", parent)
			os.Exit(1)
		}
		if !git.BranchExists(repo, branch) {
			fmt.Fprintf(os.Stderr, "Error: Branch %s does not exist\n", branch)
			os.Exit(1)
		}

		wt, err := repo.Worktree()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting worktree: %v\n", err)
			os.Exit(1)
		}

		// Checkout parent
		if err := git.CheckoutBranch(repo, parent); err != nil {
			fmt.Fprintf(os.Stderr, "Error checking out %s: %v\n", parent, err)
			os.Exit(1)
		}

		// Try remote pull first
		err = wt.Pull(&gogit.PullOptions{
			RemoteName:    viper.GetString("remote"),
			ReferenceName: plumbing.NewBranchReferenceName(branch),
		})
		if err != nil && err != gogit.NoErrAlreadyUpToDate {
			// Fall back to local merge using git command
			_, err := repo.Reference(plumbing.NewBranchReferenceName(branch), true)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting ref for %s: %v\n", branch, err)
				os.Exit(1)
			}
			cmd := exec.Command("git", "merge", branch)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error merging %s into %s: %v\n", branch, parent, err)
				os.Exit(1)
			}
		}

		fmt.Printf("Merged %s into %s\n", branch, parent)
	},
}
