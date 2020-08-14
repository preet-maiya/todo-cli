package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/preet-maiya/todo/cmd/helpers"
	"github.com/preet-maiya/todo/configuration"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ListNote struct {
	CreatedStartDate string
	CreatedEndDate   string
	StartDate        string
	EndDate          string
	OnDate           string
	ShowCreated      bool
	CaseInsensitive  bool
}

func (ln *ListNote) List(cobra *cobra.Command, args []string) {
	parsedCreatedStartDate := helpers.ParseDateOption(ln.CreatedStartDate)
	parsedCreatedEndDate := helpers.ParseDateOption(ln.CreatedEndDate)
	parsedStartDate := helpers.ParseDateOption(ln.StartDate)
	parsedEndDate := helpers.ParseDateOption(ln.EndDate)

	pattern := ""
	if len(args) > 0 {
		pattern = args[0]
	}

	if ln.OnDate != "" {
		onDate := helpers.ParseDateOption(ln.OnDate)
		parsedStartDate = onDate
		endDate, err := time.Parse("2006-01-02", onDate)
		if err != nil {
			log.Errorf("Cannot parse onDate: %s, %v", onDate, err)
			return
		}
		parsedEndDate = endDate.AddDate(0, 0, 1).Format("2006-01-02")

	}

	notes, err := db.GetNotes(parsedCreatedStartDate, parsedCreatedEndDate, parsedStartDate, parsedEndDate, pattern, ln.CaseInsensitive)
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
		if ln.ShowCreated {
			fmt.Printf("Created At: %s\n", note.CreatedAt)
		}
		fmt.Printf("\n")
	}

}
