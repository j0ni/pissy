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

func parseCommandLine() (name, category, notes string, err error) {
	name = viper.GetString("name")
	category = viper.GetString("type")
	if len(name) == 0 {
		err = errors.New("name is a required argument")
		return
	}
	if len(category) == 0 {
		err = errors.New("type is a required argument")
		return
	}
	notes = viper.GetString("notes")
	return
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

func unlockKey(passphrase []byte) (key *db.EncryptionKey, err error) {
	path := viper.GetString("path")
	key = &db.EncryptionKey{}
	err = key.Load(path, db.EncryptionKeyFile)
	if err != nil {
		return
	}
	err = key.Unlock(passphrase)
	return
}

func addRecord(cmd *cobra.Command, args []string) (err error) {
	// command line
	name, category, notes, err := parseCommandLine()
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

	addCmd.Flags().StringP("name", "n", "", "Title of the entry")
	addCmd.Flags().String("notes", "", "Additional notes")
	addCmd.Flags().StringP("type", "t", "", "The entry category")

	viper.BindPFlags(addCmd.Flags())
}
