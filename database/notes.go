package database

import (
	"database/sql"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type Note struct {
	ID        int
	Content   string
	CreatedAt string
	EndDate   string
	Status    string
	Group     string
}

func (db *DB) AddNote(content string, endDate string, groupID int) error {
	tx, err := db.conn.Begin()
	if err != nil {
		log.Errorf("Error initializing transaction: %v", err)
		return err
	}

	createNoteStmt := `
		INSERT INTO notes(created_at, content, end_date, group_id)
		VALUES(?, ?, ?, ?)
	`
	stmt, err := tx.Prepare(createNoteStmt)
	if err != nil {
		log.Errorf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	now := time.Now().Format("2006-01-02 15:04:05")

	groupIDSQL := sql.NullInt64{Int64: int64(groupID)}
	if groupID > 0 {
		groupIDSQL.Valid = true
	}

	endDateSQL := sql.NullString{String: endDate}
	if endDate != "" {
		endDateSQL.Valid = true
	}

	if _, err = stmt.Exec(now, content, endDateSQL, nil); err != nil {
		log.Errorf("Error inserting note: %v", err)
		return err
	}
	tx.Commit()
	return nil
}

func (db *DB) GetNotes(createdStartDate, createdEndDate, endStartDate, endEndDate string, pattern string, caseInsensitive bool) ([]Note, error) {
	// TODO: Fix case sensitive flag
	// caseSensitiveStmt := `PRAGMA case_sensitive_like = $1;`
	// _, err := db.conn.Query(caseSensitiveStmt, 1)
	// if err != nil {
	// 	log.Errorf("Error setting case sensitive setting: %v", err)
	// 	return []Note{}, err
	// }

	pattern = fmt.Sprintf("%%%s%%", pattern)

	getNotesStmt := `
		SELECT notes.id, notes.content, notes.created_at, notes.end_date, notes.status, groups.group_name from notes
		left join groups on notes.group_id=groups.id
		where (%s) and (%s);
	`

	timeFilters := `created_at>=? and created_at<?
					and end_date>=? and end_date<? `

	if createdStartDate == "0001-01-01" {
		timeFilters = fmt.Sprintf("%s or created_at is null", timeFilters)
	}
	if endStartDate == "0001-01-01" {
		timeFilters = fmt.Sprintf("%s or end_date is null", timeFilters)
	}

	patternFilter := "content like ?"

	getNotesStmt = fmt.Sprintf(getNotesStmt, timeFilters, patternFilter)

	println(db)
	rows, err := db.conn.Query(getNotesStmt, createdStartDate, createdEndDate, endStartDate, endEndDate, pattern)

	if err != nil {
		log.Errorf("Error getting notes: %v", err)
		return []Note{}, err
	}
	defer rows.Close()

	notes := []Note{}
	for rows.Next() {
		var id int
		var content string
		var createdAt string
		var endDateSQL sql.NullString
		var status string
		var groupNameSQL sql.NullString
		if err = rows.Scan(&id, &content, &createdAt, &endDateSQL, &status, &groupNameSQL); err != nil {
			log.Errorf("Error reading value: %v", err)
			return []Note{}, err
		}

		groupName := ""
		if groupNameSQL.Valid {
			groupName = groupNameSQL.String
		}

		endDate := ""
		if endDateSQL.Valid {
			endDate = endDateSQL.String
		}
		notes = append(notes, Note{
			ID:        id,
			Content:   content,
			CreatedAt: createdAt,
			EndDate:   endDate,
			Status:    status,
			Group:     groupName,
		})
	}
	return notes, nil
}

func (db *DB) UpdateStatus(id int, status string) error {
	tx, err := db.conn.Begin()
	if err != nil {
		log.Errorf("Error initializing transaction: %v", err)
		return err
	}

	createNoteStmt := `
		update notes set status=? where id=?;
	`
	stmt, err := tx.Prepare(createNoteStmt)
	if err != nil {
		log.Errorf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(status, id); err != nil {
		log.Errorf("Error updating status for id %d: %v", id, err)
		return err
	}
	tx.Commit()
	return nil
}
