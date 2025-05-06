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

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync child branches with parent",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := git.GetRepo()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		syncMode := viper.GetString("syncMode")
		if !config.Contains([]string{"rebase", "merge", "reset"}, syncMode) {
			fmt.Fprintf(os.Stderr, "Error: Invalid syncMode '%s'. Must be 'rebase', 'merge', or 'reset'\n", syncMode)
			os.Exit(1)
		}

		branches := viper.GetStringSlice("branches")
		base := viper.GetString("base")
		if base == "" {
			fmt.Fprintf(os.Stderr, "Error: Base branch not set in .stackpr.yaml\n")
			os.Exit(1)
		}

		for i, branch := range branches {
			if !git.BranchExists(repo, branch) {
				fmt.Fprintf(os.Stderr, "Error: Branch %s does not exist\n", branch)
				os.Exit(1)
			}

			parent := base
			if i > 0 {
				parent = branches[i-1]
			}
			if !git.BranchExists(repo, parent) {
				fmt.Fprintf(os.Stderr, "Error: Parent branch %s does not exist\n", parent)
				os.Exit(1)
			}

			wt, err := repo.Worktree()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting worktree: %v\n", err)
				os.Exit(1)
			}

			if err := git.CheckoutBranch(repo, branch); err != nil {
				fmt.Fprintf(os.Stderr, "Error checking out %s: %v\n", branch, err)
				os.Exit(1)
			}

			switch syncMode {
			case "rebase":
				_, err := repo.Reference(plumbing.NewBranchReferenceName(parent), true)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error getting parent ref for %s: %v\n", parent, err)
					os.Exit(1)
				}
				// Use git command directly since go-git doesn't support rebase
				cmd := exec.Command("git", "rebase", parent)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					fmt.Fprintf(os.Stderr, "Error rebasing %s onto %s: %v\n", branch, parent, err)
					os.Exit(1)
				}
			case "merge":
				// Try remote pull first
				err := wt.Pull(&gogit.PullOptions{
					RemoteName:    viper.GetString("remote"),
					ReferenceName: plumbing.NewBranchReferenceName(parent),
				})
				if err != nil && err != gogit.NoErrAlreadyUpToDate {
					// Fall back to local merge using git command
					_, err := repo.Reference(plumbing.NewBranchReferenceName(branch), true)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error getting ref for %s: %v\n", branch, err)
						os.Exit(1)
					}
					cmd := exec.Command("git", "merge", parent)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						fmt.Fprintf(os.Stderr, "Error merging %s into %s: %v\n", parent, branch, err)
						os.Exit(1)
					}
				}
			case "reset":
				parentRef, err := repo.Reference(plumbing.NewBranchReferenceName(parent), true)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error getting parent ref for %s: %v\n", parent, err)
					os.Exit(1)
				}
				if err := wt.Reset(&gogit.ResetOptions{Mode: gogit.HardReset, Commit: parentRef.Hash()}); err != nil {
					fmt.Fprintf(os.Stderr, "Error resetting %s to %s: %v\n", branch, parent, err)
					os.Exit(1)
				}
			}
			fmt.Printf("Synced %s with %s\n", branch, parent)
		}
	},
}
