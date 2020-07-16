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
	"time"

	"github.com/preet-maiya/todo/cmd/helpers"
	"github.com/preet-maiya/todo/database"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var createdEndDateStr string
var createdStartDateStr string
var endDateStr string
var startDateStr string
var endDateCreate string
var content string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: create,
}

func create(*cobra.Command, []string) {
	db := database.NewDB(config.DBFile)
	if content == "" {
		contentBytes, err := helpers.CaptureInputFromEditor()
		if err != nil {
			log.Errorf("Error capturing input from $EDITOR: %v", err)
			return
		}
		content = string(contentBytes)
	}

	if err := db.AddNote(content, endDateCreate, 0); err != nil {
		log.Errorf("Error adding note: %v", err)
		return
	}
}

// createCmd represents the create command
var notesCmd = &cobra.Command{
	Use:   "notes",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: notes,
}

func notes(*cobra.Command, []string) {
	db := database.NewDB(config.DBFile)
	// 	createdStartDate := time.Parse("2006-01-02", createdStartDateStr)
	// 	createdEndDate := time.Parse("2006-01-02", createdEndDateStr)
	// 	startDate := time.Parse("2006-01-02", startDateStr)
	// 	endDate := time.Parse("2006-01-02", endDateStr)
	notes, err := db.GetNotes(createdStartDateStr, createdEndDateStr, startDateStr, endDateStr)
	if err != nil {
		log.Errorf("Error fetching notes: %v", err)
		return
	}

	for i, note := range notes {
		fmt.Printf("%d: %s\nCreated At: %s\nEnd Date: %s\n\n", i, note.Content, note.CreatedAt, note.EndDate)
	}

}

func init() {
	rootCmd.AddCommand(notesCmd)
	notesCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	beginning := time.Time{}.Format("2006-01-02")
	veryFarFuture := time.Now().AddDate(1000, 0, 0).Format("2006-01-02")
	notesCmd.Flags().StringVarP(&endDateStr, "end-date-before", "", veryFarFuture, "End date for the note")
	notesCmd.Flags().StringVarP(&startDateStr, "end-date-after", "", beginning, "End date for the note")
	notesCmd.Flags().StringVarP(&createdEndDateStr, "created-before", "", veryFarFuture, "End date for the note")
	notesCmd.Flags().StringVarP(&createdStartDateStr, "created-after", "", beginning, "End date for the note")
	createCmd.Flags().StringVarP(&content, "content", "m", "", "Content of the note")
	createCmd.Flags().StringVarP(&endDateCreate, "end-date", "", "", "End date for the note")
}