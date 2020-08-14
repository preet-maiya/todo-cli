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
	"time"

	"github.com/preet-maiya/todo/cmd/handlers"
	"github.com/spf13/cobra"
)

var createNote handlers.CreateNote
var listNotes handlers.ListNote
var status handlers.Status

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create notes quickly",
	Long: `Create notes with expected end date with your favourite editor or
	directly from the command line. To use editor, set EDITOR env variable to
	the executable of the editor. emacs or vim, I'm not judging!`,
	Run: createNote.Create,
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
	Run:  listNotes.List,
}

// updateStatusCmd represents the updateStatus command
var updateStatusCmd = &cobra.Command{
	Use:   "status [pattern]",
	Short: "Change status of a note",
	// TODO: Fill this
	Long: ``,
	Args: cobra.ExactArgs(1),
	Run:  status.Update,
}

func init() {
	rootCmd.AddCommand(notesCmd)
	notesCmd.AddCommand(createCmd)
	notesCmd.AddCommand(updateStatusCmd)

	beginning := time.Time{}.Format("2006-01-02")
	veryFarFuture := time.Now().AddDate(1000, 0, 0).Format("2006-01-02")
	notesCmd.Flags().StringVarP(&listNotes.EndDate, "end-date-before", "B", veryFarFuture, "End date for the note")
	notesCmd.Flags().StringVarP(&listNotes.StartDate, "end-date-after", "A", beginning, "End date for the note")
	notesCmd.Flags().StringVarP(&listNotes.CreatedEndDate, "created-before", "", veryFarFuture, "End date for the note")
	notesCmd.Flags().StringVarP(&listNotes.CreatedStartDate, "created-after", "", beginning, "End date for the note")
	notesCmd.Flags().BoolVarP(&listNotes.ShowCreated, "show-created", "c", false, "Show created at time for notes")
	notesCmd.Flags().StringVarP(&listNotes.OnDate, "on", "O", "", "On the date")
	notesCmd.Flags().BoolVarP(&listNotes.CaseInsensitive, "", "i", false, "Search pattern case insensitive")

	createCmd.Flags().StringVarP(&createNote.Content, "content", "m", "", "Content of the note")
	createCmd.Flags().StringVarP(&createNote.EndDate, "end-date", "e", "", "End date for the note")

	updateStatusCmd.Flags().BoolVarP(&status.CaseInsensitive, "", "i", false, "Search pattern case insensitive")
	updateStatusCmd.Flags().BoolVarP(&status.Created, "created", "c", false, "Change status to created")
	updateStatusCmd.Flags().BoolVarP(&status.Pending, "pending", "p", false, "Change status to pending")
	updateStatusCmd.Flags().BoolVarP(&status.Done, "done", "d", false, "Change status to done")
	updateStatusCmd.Flags().BoolVarP(&status.Expired, "expired", "e", false, "Change status to expired")
	updateStatusCmd.Flags().StringVarP(&status.EndDate, "end-date-before", "B", veryFarFuture, "End date for the note")
	updateStatusCmd.Flags().StringVarP(&status.StartDate, "end-date-after", "A", beginning, "End date for the note")
	updateStatusCmd.Flags().StringVarP(&status.CreatedEndDate, "created-before", "", veryFarFuture, "End date for the note")
	updateStatusCmd.Flags().StringVarP(&status.CreatedStartDate, "created-after", "", beginning, "End date for the note")
}
