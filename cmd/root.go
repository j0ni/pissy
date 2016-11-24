package cmd

import (
	"errors"
	"fmt"
	"os"

	"mig.ninja/mig/pgp/pinentry"

	"github.com/j0ni/pissy/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "pissy",
	Short: "A password manager",
	Long:  `Try it and see...`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(
		&cfgFile, "config", "", "config file (default is $HOME/.pissy.yaml)")
	RootCmd.PersistentFlags().StringP("path", "p", "", "path to keychain")

	RootCmd.PersistentFlags().StringP("name", "n", "", "Title of the entry")
	RootCmd.PersistentFlags().String("notes", "", "Additional notes")
	RootCmd.PersistentFlags().StringP("type", "t", "", "The entry category")
	RootCmd.PersistentFlags().StringP("uuid", "u", "", "UUID of the record to dump")
	RootCmd.PersistentFlags().BoolP(
		"update-secret", "s", false, "Update the secret (invokes UI)")

	viper.BindPFlags(RootCmd.PersistentFlags())
}

func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetDefault("path", fmt.Sprintf("%s/%s", os.Getenv("HOME"), "Dropbox/Apps/Pissy"))
	viper.BindPFlag("path", RootCmd.Flags().Lookup("path"))

	viper.SetConfigName(".pissy") // name of config file (without extension)
	viper.AddConfigPath("$HOME")  // adding home directory as first search path
	viper.SetEnvPrefix("pissy")   // add prefix for env vars
	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func acquirePassphrase() ([]byte, error) {
	request := &pinentry.Request{
		Desc:   "Passphrase dialog for Pissy",
		Prompt: "Enter passphrase",
	}
	passphrase, err := request.GetPIN()
	return []byte(passphrase), err
}

func parseCommonFieldsUpdate() (name, category, notes string) {
	name = viper.GetString("name")
	category = viper.GetString("type")
	notes = viper.GetString("notes")
	return
}

func parseCommonFieldsCreate() (name, category, notes string, err error) {
	name, category, notes = parseCommonFieldsUpdate()
	if len(name) == 0 {
		err = errors.New("name is a required argument")
		return
	}
	if len(category) == 0 {
		err = errors.New("type is a required argument")
	}
	return
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
