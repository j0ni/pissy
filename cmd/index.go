package cmd

import (
	"fmt"

	"github.com/j0ni/pissy/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "List the contents of the keychain",
	Long:  `List the contents of the keychain`,
	Run: func(cmd *cobra.Command, args []string) {
		database := db.New(viper.Get("path").(string))
		err := database.Load()
		if err != nil {
			panic(err)
		}
		for _, record := range database.Records {
			fmt.Println(record)
		}
	},
}

func init() {
	RootCmd.AddCommand(indexCmd)
}
