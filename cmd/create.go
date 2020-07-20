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
var onDate string
var showCreated bool

var endDateCreate string
var content string
var caseInsensitive bool

var created, pending, done, expired bool
var listCreatedStartDateStr, listCreatedEndDateStr, listStartDateStr, listEndDateStr string
var updateCaseInsensitive bool

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

	if content == "" {
		log.Warning("Empty content given. Skipping...")
		return
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

	if onDate != "" {
		onDate = helpers.ParseDateOption(onDate)
		parsedStartDate = onDate
		endDate, err := time.Parse("2006-01-02", onDate)
		if err != nil {
			log.Errorf("Cannot parse onDate: %s, %v", onDate, err)
			return
		}
		parsedEndDate = endDate.AddDate(0, 0, 1).Format("2006-01-02")

	}

	notes, err := db.GetNotes(parsedCreatedStartDate, parsedCreatedEndDate, parsedStartDate, parsedEndDate, pattern, caseInsensitive)
	if err != nil {
		log.Errorf("Error fetching notes: %v", err)
		return
	}

	for i, note := range notes {
		// TODO: Refactor this block!

		if note.EndDate != "" {
			endDateParsed, err := time.Parse("2006-01-02", note.EndDate)
			if err != nil {
				log.Errorf("Cannot parse %s: %v", note.EndDate, err)
				return
			}

			if note.Status != configuration.Done && time.Now().After(endDateParsed) {
				note.Status = configuration.Expired
			}
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

// updateStatusCmd represents the updateStatus command
var updateStatusCmd = &cobra.Command{
	Use:   "status [pattern]",
	Short: "Change status of a note",
	// TODO: Fill this
	Long: ``,
	Args: cobra.ExactArgs(1),
	Run:  updateStatus,
}

func updateStatus(cobra *cobra.Command, args []string) {
	if nFlagsSet := helpers.NSet(created, pending, done, expired); nFlagsSet > 1 {
		log.Warning("More than one flag set")
		fmt.Println("More than one flag set. Please set only one")
		return
	}

	db := database.NewDB(config.DBFile)
	parsedCreatedStartDate := helpers.ParseDateOption(listCreatedStartDateStr)
	parsedCreatedEndDate := helpers.ParseDateOption(listCreatedEndDateStr)
	parsedStartDate := helpers.ParseDateOption(listStartDateStr)
	parsedEndDate := helpers.ParseDateOption(listEndDateStr)

	pattern := args[0]
	notes, err := db.GetNotes(parsedCreatedStartDate, parsedCreatedEndDate, parsedStartDate, parsedEndDate, pattern, updateCaseInsensitive)
	if err != nil {
		log.Errorf("Error fetching notes: %v", err)
		return
	}

	// TODO: Handle more than
	if len(notes) > 1 {
		fmt.Println("More than one notes found for the filter")
		return
	}

	status := configuration.Created

	if created {
		status = configuration.Created
	} else if pending {
		status = configuration.Pending
	} else if done {
		status = configuration.Done
	} else if expired {
		status = configuration.Expired
	}

	if err = db.UpdateStatus(notes[0].ID, status); err != nil {
		log.Errorf("Error updating status: %v", err)
		return
	}
}

func init() {
	rootCmd.AddCommand(notesCmd)
	notesCmd.AddCommand(createCmd)
	notesCmd.AddCommand(updateStatusCmd)

	beginning := time.Time{}.Format("2006-01-02")
	veryFarFuture := time.Now().AddDate(1000, 0, 0).Format("2006-01-02")
	notesCmd.Flags().StringVarP(&endDateStr, "end-date-before", "B", veryFarFuture, "End date for the note")
	notesCmd.Flags().StringVarP(&startDateStr, "end-date-after", "A", beginning, "End date for the note")
	notesCmd.Flags().StringVarP(&createdEndDateStr, "created-before", "", veryFarFuture, "End date for the note")
	notesCmd.Flags().StringVarP(&createdStartDateStr, "created-after", "", beginning, "End date for the note")
	notesCmd.Flags().BoolVarP(&showCreated, "show-created", "c", false, "Show created at time for notes")
	notesCmd.Flags().StringVarP(&onDate, "on", "O", "", "On the date")
	notesCmd.Flags().BoolVarP(&caseInsensitive, "", "i", false, "Search pattern case insensitive")

	createCmd.Flags().StringVarP(&content, "content", "m", "", "Content of the note")
	createCmd.Flags().StringVarP(&endDateCreate, "end-date", "e", "", "End date for the note")

	updateStatusCmd.Flags().BoolVarP(&updateCaseInsensitive, "", "i", false, "Search pattern case insensitive")
	updateStatusCmd.Flags().BoolVarP(&created, "created", "c", false, "Change status to created")
	updateStatusCmd.Flags().BoolVarP(&pending, "pending", "p", false, "Change status to pending")
	updateStatusCmd.Flags().BoolVarP(&done, "done", "d", false, "Change status to done")
	updateStatusCmd.Flags().BoolVarP(&expired, "expired", "e", false, "Change status to expired")
	updateStatusCmd.Flags().StringVarP(&listEndDateStr, "end-date-before", "B", veryFarFuture, "End date for the note")
	updateStatusCmd.Flags().StringVarP(&listStartDateStr, "end-date-after", "A", beginning, "End date for the note")
	updateStatusCmd.Flags().StringVarP(&listCreatedEndDateStr, "created-before", "", veryFarFuture, "End date for the note")
	updateStatusCmd.Flags().StringVarP(&listCreatedStartDateStr, "created-after", "", beginning, "End date for the note")
}
