package cmd

import (
	"errors"
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/j0ni/pissy/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// grabCmd represents the grab command
var grabCmd = &cobra.Command{
	Use:   "grab",
	Short: "Retrieve a single password.",
	Long:  `Retrieve a single password using a UUID.`,
	RunE:  findRecord,
}

func findRecord(cmd *cobra.Command, args []string) error {
	// dependencies
	path := viper.GetString("path")
	uuid := viper.GetString("uuid")
	if len(uuid) == 0 {
		return errors.New("UUID cannot be empty")
	}
	// load the database
	database := db.New(path)
	err := database.Load()
	if err != nil {
		return err
	}
	// find the record
	rec, err := database.Find(uuid)
	if err != nil {
		return err
	}
	// db passphrase
	passphrase, err := acquirePassphrase()
	if err != nil {
		return err
	}
	// unlock the master key
	key, err := unlockKey(passphrase)
	if err != nil {
		return err
	}
	// decrypt the secret
	err = rec.Decrypt(key.DecryptedKey)
	if err != nil {
		return err
	}
	// output the secret
	if viper.GetBool("clipboard") {
		err = clipboard.WriteAll(string(rec.DecryptedValue))
		if err != nil {
			return err
		}
	} else {
		fmt.Println(string(rec.DecryptedValue))
	}
	return nil
}

func init() {
	RootCmd.AddCommand(grabCmd)

	grabCmd.Flags().StringP("uuid", "u", "", "UUID of the record to dump")
	grabCmd.Flags().BoolP("clipboard", "y", false, "Send secret to clipboard")

	viper.BindPFlags(grabCmd.Flags())
}
