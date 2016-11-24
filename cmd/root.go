package cmd

import (
	"fmt"
	"os"

	"mig.ninja/mig/pgp/pinentry"

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
