package database

import (
	log "github.com/sirupsen/logrus"
)

type Group struct {
	ID        int
	GroupName string
}

func (db *DB) GetGroups() ([]Group, error) {
	getGroupsStmt := `
		SELECT id, group_name from groups;
	`
	rows, err := db.conn.Query(getGroupsStmt)
	if err != nil {
		log.Errorf("Error getting groups: %v", err)
		return []Group{}, err
	}
	defer rows.Close()

	groups := []Group{}
	for rows.Next() {
		var id int
		var groupName string
		err = rows.Scan(&id, &groupName)
		if err != nil {
			log.Errorf("Error reading value: %v", err)
			return []Group{}, err
		}
		groups = append(groups, Group{
			ID:        id,
			GroupName: groupName,
		})
	}
	return groups, nil
}

func (db *DB) AddGroup(groupName string) error {
	tx, err := db.conn.Begin()
	if err != nil {
		log.Errorf("Error initializing transaction: %v", err)
		return err
	}

	createGroupStmt := `
		INSERT INTO group(group_name)
		VALUES(?)
	`
	stmt, err := tx.Prepare(createGroupStmt)
	if err != nil {
		log.Errorf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(groupName); err != nil {
		log.Errorf("Error inserting group: %v", err)
		return err
	}
	tx.Commit()
	return nil
}
