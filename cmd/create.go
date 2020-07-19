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
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/preet-maiya/todo/cmd/helpers"
	"github.com/preet-maiya/todo/configuration"
	"github.com/preet-maiya/todo/database"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var createdEndDateStr string
var createdStartDateStr string
var endDateStr string
var startDateStr string
var showCreated bool
var endDateCreate string
var content string
var caseInsensitive bool

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create notes quickly",
	Long: `Create notes with expected end date with your favourite editor or
	directly from the command line. To use editor, set EDITOR env variable to
	the executable of the editor. emacs or vim, I'm not judging!`,
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

	parsedEndDateCreate := helpers.ParseDateOption(endDateCreate)

	if err := db.AddNote(content, parsedEndDateCreate, 0); err != nil {
		log.Errorf("Error adding note: %v", err)
		return
	}
}

// createCmd represents the create command
var notesCmd = &cobra.Command{
	Use:   "notes [pattern]",
	Short: "List all notes with details",
	Long: `List all the notes created from the beginning of the time!
	Control what to show with flags, like whether to show the created date
	It is not shown by default but passing flag -c will show it.
	An optional pattern to search can also be given`,
	Args: cobra.MaximumNArgs(1),
	Run:  notes,
}

func notes(cobra *cobra.Command, args []string) {
	db := database.NewDB(config.DBFile)
	parsedCreatedStartDate := helpers.ParseDateOption(createdStartDateStr)
	parsedCreatedEndDate := helpers.ParseDateOption(createdEndDateStr)
	parsedStartDate := helpers.ParseDateOption(startDateStr)
	parsedEndDate := helpers.ParseDateOption(endDateStr)

	pattern := ""
	if len(args) > 0 {
		pattern = args[0]
	}

	notes, err := db.GetNotes(parsedCreatedStartDate, parsedCreatedEndDate, parsedStartDate, parsedEndDate, pattern, caseInsensitive)
	if err != nil {
		log.Errorf("Error fetching notes: %v", err)
		return
	}

	for i, note := range notes {
		// TODO: Refactor this block!

		endDateParsed, err := time.Parse("2006-01-02", note.EndDate)
		if err != nil {
			log.Errorf("Cannot parse %s: %v", note.EndDate, err)
			return
		}

		if note.Status != configuration.Done && time.Now().After(endDateParsed) {
			note.Status = configuration.Expired
		}

		fmt.Printf("%d: ", i)

		color.Cyan("%s\n", note.Content)
		statusColor := color.New(color.FgWhite).SprintFunc()
		switch note.Status {
		case configuration.Created:
			statusColor = color.New(color.FgWhite).SprintFunc()
		case configuration.Pending:
			statusColor = color.New(color.FgYellow).SprintFunc()
		case configuration.Done:
			statusColor = color.New(color.FgGreen).SprintFunc()
		case configuration.Expired:
			statusColor = color.New(color.FgRed).SprintFunc()
		default:
			statusColor = color.New(color.FgWhite).SprintFunc()
		}
		fmt.Printf("Status: %s\n", statusColor(strings.ToTitle(note.Status)))
		fmt.Printf("End Date: %s\n", note.EndDate)
		if showCreated {
			fmt.Printf("Created At: %s\n", note.CreatedAt)
		}
		fmt.Printf("\n")
	}

}

func init() {
	rootCmd.AddCommand(notesCmd)
	notesCmd.AddCommand(createCmd)

	beginning := time.Time{}.Format("2006-01-02")
	veryFarFuture := time.Now().AddDate(1000, 0, 0).Format("2006-01-02")
	notesCmd.Flags().StringVarP(&endDateStr, "end-date-before", "B", veryFarFuture, "End date for the note")
	notesCmd.Flags().StringVarP(&startDateStr, "end-date-after", "A", beginning, "End date for the note")
	notesCmd.Flags().StringVarP(&createdEndDateStr, "created-before", "", veryFarFuture, "End date for the note")
	notesCmd.Flags().StringVarP(&createdStartDateStr, "created-after", "", beginning, "End date for the note")
	notesCmd.Flags().BoolVarP(&showCreated, "show-created", "c", false, "Show created at time for notes")
	notesCmd.Flags().BoolVarP(&caseInsensitive, "", "i", false, "Search pattern case insensitive")
	createCmd.Flags().StringVarP(&content, "content", "m", "", "Content of the note")
	createCmd.Flags().StringVarP(&endDateCreate, "end-date", "e", "", "End date for the note")
}
