package handlers

import (
	"github.com/preet-maiya/todo/cmd/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type CreateNote struct {
	Content string
	EndDate string
}

func (cn *CreateNote) Create(*cobra.Command, []string) {
	if cn.Content == "" {
		contentBytes, err := helpers.CaptureInputFromEditor()
		if err != nil {
			log.Errorf("Error capturing input from $EDITOR: %v", err)
			return
		}
		cn.Content = string(contentBytes)
	}

	if cn.Content == "" {
		log.Warning("Empty content given. Skipping...")
		return
	}

	parsedEndDateCreate := helpers.ParseDateOption(cn.EndDate)
	if parsedEndDateCreate == "01-01-0001" {
		parsedEndDateCreate = ""
	}

	if err := db.AddNote(cn.Content, parsedEndDateCreate, 0); err != nil {
		log.Errorf("Error adding note: %v", err)
		return
	}

}
