package handlers

import (
	"github.com/preet-maiya/todo/configuration"
	"github.com/preet-maiya/todo/database"
)

var db database.DB

func InitConfig(config configuration.Configuration) {
	db = database.NewDB(config.DBFile)
}
