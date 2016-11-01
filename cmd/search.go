package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for a password",
	Long:  `Perform a regular expression search for a password`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("search called")
	},
}

func init() {
	RootCmd.AddCommand(searchCmd)
}
