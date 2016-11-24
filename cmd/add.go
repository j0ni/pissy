package cmd

import (
	"errors"

	"github.com/howeyc/gopass"
	"github.com/j0ni/pissy/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new record to the database",
	Long:  `Add an record to the database. The secret is read from stdin, or using pinentry`,
	RunE:  addRecord,
}

func readSecret() (secret []byte, err error) {
	secret, err = gopass.GetPasswd()
	if err != nil {
		return
	}
	if len(secret) == 0 {
		secret, err = acquirePassphrase()
		if err != nil {
			return
		}
	}
	if len(secret) == 0 {
		err = errors.New("secret cannot be empty")
	}
	return
}

func createRecord(name, category, notes string, secret, key []byte) error {
	path := viper.GetString("path")
	record := db.NewRecord()
	record.Title = name
	record.TypeName = category
	record.Notes = notes
	record.DecryptedValue = secret
	err := record.Encrypt(key)
	if err != nil {
		return err
	}
	return record.Save(path)
}

func addRecord(cmd *cobra.Command, args []string) (err error) {
	// command line
	name, category, notes, err := parseCommonFieldsCreate()
	if err != nil {
		return
	}
	// secret from stdin or pinentry
	secret, err := readSecret()
	if err != nil {
		return
	}
	// db passphrase
	passphrase, err := acquirePassphrase()
	if err != nil {
		return
	}
	// unlock the db master key
	key, err := unlockKey(passphrase)
	if err != nil {
		return
	}
	// make and save the new record
	return createRecord(name, category, notes, secret, key.DecryptedKey)
}

func init() {
	RootCmd.AddCommand(addCmd)
}
