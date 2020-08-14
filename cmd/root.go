/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/preet-maiya/todo/cmd/handlers"
	"github.com/preet-maiya/todo/configuration"
	"github.com/spf13/viper"
)

var cfgFile string
var config configuration.Configuration

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "A simple TODO cli app",
	Long: `This app is for someone who likes the power and comfort of a command line
	Create a list of TODOS, reminders without leaving command line!
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Simple CLI Todo app")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.todo.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Set Defaults
	defaultDBFile := fmt.Sprintf("%s/%s", home, ".todo.db")
	viper.SetDefault("db_file", defaultDBFile)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".todo" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".todo")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Infof("Using config file:", viper.ConfigFileUsed())
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Errorf("Error unmarshaling config: %v", err)
	}
	handlers.InitConfig(config)
}
