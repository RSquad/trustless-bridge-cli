/*
Copyright © 2025 RSquad <hello@rsquad.io>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/rsquad/trustless-bridge-cli/internal/tonclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var tonClient *tonclient.TonClient

var rootCmd = &cobra.Command{
	Use:   "trustless-bridge-cli",
	Short: "A CLI tool for data preparation and retrieval for the Trustless Bridge",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		"config file (default is $HOME/.trustless-bridge-cli.yaml)",
	)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	viper.SetDefault("ton_config_url", "https://ton-blockchain.github.io/testnet-global.config.json")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".trustless-bridge-cli")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	tonClientURL := viper.GetString("ton_config_url")
	if tonClientURL == "" {
		panic("TonClient URL not set in config")
	}

	var err error
	tonClient, err = tonclient.NewTonClient(tonClientURL)
	if err != nil {
		panic(err)
	}
}
