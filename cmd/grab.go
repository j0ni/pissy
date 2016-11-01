package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// grabCmd represents the grab command
var grabCmd = &cobra.Command{
	Use:   "grab",
	Short: "Retrieve a single password.",
	Long:  `Retrieve a single password using a UUID.`,
	Run: func(cmd *cobra.Command, args []string) {

		// TODO: Work your own magic here
		fmt.Printf("grab called with args: %s\n", args)
	},
}

func init() {
	RootCmd.AddCommand(grabCmd)
}
