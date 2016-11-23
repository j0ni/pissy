package cmd

import (
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

func initDb(cmd *cobra.Command, args []string) error {
	path := viper.GetString("path")

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
	configFile := fmt.Sprintf("%s/%s", os.Getenv("HOME"), ".pissy.yaml")
	return ioutil.WriteFile(configFile, configStr, 0600)
}
