package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/j0ni/pissy/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize an empty pissy database",
	Long:  `Initialize an empty pissy database`,
	RunE:  initDb,
}

func init() {
	RootCmd.AddCommand(initCmd)
}

func dbExists(path, configFile string) (found bool, err error) {
	if found, err = db.Exists(path); err != nil {
		return
	} else if found {
		encryptedKeyFile := fmt.Sprintf("%s/%s", path, db.EncryptionKeyFile)
		if found, err = db.Exists(encryptedKeyFile); err != nil || found {
			return
		}
	}

	found, err = db.Exists(configFile)
	return
}

func initDb(cmd *cobra.Command, args []string) error {
	path := viper.GetString("path")
	configFile := fmt.Sprintf("%s/%s", os.Getenv("HOME"), ".pissy.yaml")

	if ok, err := dbExists(path, configFile); ok {
		return errors.New("Looks like there's already a DB")
	} else if err != nil {
		return err
	}

	err := os.MkdirAll(path, 0700)
	if err != nil {
		return err
	}

	passphrase, err := acquirePassphrase()
	if err != nil {
		return err
	}

	key := db.NewKey()

	err = key.Lock(passphrase)
	if err != nil {
		return err
	}

	err = key.Save(path, db.EncryptionKeyFile)
	if err != nil {
		return err
	}

	configStr := []byte(fmt.Sprintf("path: \"%s\"\n", path))
	return ioutil.WriteFile(configFile, configStr, 0600)
}
