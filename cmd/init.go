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
	"github.com/preet-maiya/todo/database"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize everything needed for the app to run. Only to be run once",
	Long: `This command is to initialize the app. It is to be run only once.
	Running it the second time may mess up the app :\`,
	Run: initDB,
}

func initDB(*cobra.Command, []string) {
	db := database.NewDB(config.DBFile)
	err := db.InitDB()
	if err != nil {
		log.Errorf("Error initializing DB: %v", err)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
}
