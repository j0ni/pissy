package cmd

import (
	"fmt"

	"github.com/j0ni/pissy/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for a password",
	Long:  `Perform a regular expression search for a password`,
	RunE:  search,
}

func search(cmd *cobra.Command, args []string) error {
	// load the database
	database := db.New(viper.GetString("path"))
	err := database.Load()
	if err != nil {
		return err
	}
	reString := viper.GetString("regexp")
	database, err = database.Filter(reString)
	if err != nil {
		return err
	}
	for _, rec := range database.Records {
		fmt.Println(rec)
	}
	return nil
}

func init() {
	RootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringP("regexp", "r", "", "Regular expression to use in search")
	viper.BindPFlags(searchCmd.Flags())
}
