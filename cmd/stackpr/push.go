package stackpr

import (
	"fmt"
	"os"

	"test_cli/pkg/git"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	forcePush bool
	debugPush bool
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push all stacked branches",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := git.GetRepo()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		remoteName := viper.GetString("remote")
		branches := viper.GetStringSlice("branches")

		if debugPush {
			fmt.Printf("Debug: Pushing to remote %s\n", remoteName)
			fmt.Printf("Debug: Branches to push: %v\n", branches)
		}

		for _, branch := range branches {
			if !git.BranchExists(repo, branch) {
				fmt.Fprintf(os.Stderr, "Error: Branch %s does not exist\n", branch)
				os.Exit(1)
			}
			pushOptions := &gogit.PushOptions{
				RemoteName: remoteName,
				RefSpecs:   []config.RefSpec{config.RefSpec(fmt.Sprintf("refs/heads/%s:refs/heads/%s", branch, branch))},
				Force:      forcePush,
			}
			if debugPush {
				fmt.Printf("Debug: Pushing branch %s with options: %+v\n", branch, pushOptions)
			}
			if err := repo.Push(pushOptions); err != nil {
				if debugPush {
					fmt.Printf("Debug: Push error: %v\n", err)
				}
				fmt.Fprintf(os.Stderr, "Error pushing %s: %v\n", branch, err)
				os.Exit(1)
			}
			fmt.Printf("Pushed %s\n", branch)
		}
	},
}

func init() {
	pushCmd.Flags().BoolVar(&forcePush, "force", false, "Force push branches")
	pushCmd.Flags().BoolVar(&debugPush, "debug", false, "Enable debug output for push command")
}
