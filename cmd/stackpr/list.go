package stackpr

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List branches in the stack",
	Run: func(cmd *cobra.Command, args []string) {
		branches := viper.GetStringSlice("branches")
		if len(branches) == 0 {
			fmt.Println("No stacked branches")
			return
		}
		fmt.Println("Stacked branches:")
		for i, branch := range branches {
			fmt.Printf("%d: %s\n", i+1, branch)
		}
	},
}
