package handlers

import (
	"fmt"

	"github.com/preet-maiya/todo/cmd/helpers"
	"github.com/preet-maiya/todo/configuration"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Status struct {
	Pending          bool
	Done             bool
	Expired          bool
	CreatedStartDate string
	CreatedEndDate   string
	StartDate        string
	EndDate          string
	Created          bool
	CaseInsensitive  bool
}

func (st *Status) Update(cobra *cobra.Command, args []string) {
	if nFlagsSet := helpers.NSet(st.Created, st.Pending, st.Done, st.Expired); nFlagsSet > 1 {
		log.Warning("More than one flag set")
		fmt.Println("More than one flag set. Please set only one")
		return
	}

	parsedCreatedStartDate := helpers.ParseDateOption(st.CreatedStartDate)
	parsedCreatedEndDate := helpers.ParseDateOption(st.CreatedEndDate)
	parsedStartDate := helpers.ParseDateOption(st.StartDate)
	parsedEndDate := helpers.ParseDateOption(st.EndDate)

	pattern := args[0]
	notes, err := db.GetNotes(parsedCreatedStartDate, parsedCreatedEndDate, parsedStartDate, parsedEndDate, pattern, st.CaseInsensitive)
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

	if st.Created {
		status = configuration.Created
	} else if st.Pending {
		status = configuration.Pending
	} else if st.Done {
		status = configuration.Done
	} else if st.Expired {
		status = configuration.Expired
	}

	if err = db.UpdateStatus(notes[0].ID, status); err != nil {
		log.Errorf("Error updating status: %v", err)
		return
	}
}
