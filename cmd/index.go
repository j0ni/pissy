package cmd

import (
	"fmt"

	"github.com/j0ni/pissy/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "List the contents of the keychain",
	Long:  `List the contents of the keychain`,
	RunE:  index,
}

func index(cmd *cobra.Command, args []string) error {
	database := db.New(viper.GetString("path"))
	err := database.Load()
	if err != nil {
		return err
	}
	for _, record := range database.Records {
		fmt.Println(record)
	}
	return nil
}

func init() {
	RootCmd.AddCommand(indexCmd)
}
