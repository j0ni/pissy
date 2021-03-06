package cmd

import (
	"github.com/j0ni/pissy/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the content of a specified record",
	Long:  `Update the content of a specified record`,
	RunE:  updateRecord,
}

func updateRecord(cmd *cobra.Command, args []string) (err error) {
	// dependencies
	path := viper.GetString("path")
	// command line
	name, category, notes := parseCommonFieldsUpdate()
	uuid := viper.GetString("uuid")
	// load the database
	database := db.New(path)
	err = database.Load()
	if err != nil {
		return
	}
	// find the record
	rec, err := database.Find(uuid)
	if err != nil {
		return
	}
	// maybe get the secret
	var secret, passphrase, keyMaterial []byte
	if viper.GetBool("update-secret") {
		secret, err = readSecret()
		if err != nil {
			return
		}
		// db passphrase, only needed for crypto
		passphrase, err = acquirePassphrase()
		if err != nil {
			return
		}
		// unlock the master key
		var key *db.EncryptionKey
		key, err = unlockKey(passphrase)
		if err != nil {
			return
		}
		keyMaterial = key.DecryptedKey
	}
	// update the record
	rec.UpdateFields(keyMaterial, &db.Record{
		Title:          name,
		TypeName:       category,
		Notes:          notes,
		DecryptedValue: secret,
	})
	return rec.Save(path)
}

func init() {
	RootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolP("update-secret", "s", false, "Update the secret (invokes UI)")
	viper.BindPFlags(updateCmd.Flags())
}
