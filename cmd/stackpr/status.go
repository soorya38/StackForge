package stackpr

import (
	"fmt"
	"os"
	"strings"

	"test_cli/pkg/git"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of stacked branches",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := git.GetRepo()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		branches := viper.GetStringSlice("branches")
		for _, branch := range branches {
			if !git.BranchExists(repo, branch) {
				fmt.Fprintf(os.Stderr, "Warning: Branch %s does not exist\n", branch)
				continue
			}
			ref, err := repo.Reference(plumbing.NewBranchReferenceName(branch), true)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting ref for %s: %v\n", branch, err)
				continue
			}
			commit, err := repo.CommitObject(ref.Hash())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting commit for %s: %v\n", branch, err)
				continue
			}
			fmt.Printf("%s: %s (%s)\n", branch, strings.TrimSpace(commit.Message), commit.Hash.String()[:7])
		}
	},
}
