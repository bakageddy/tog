package util

import "database/sql"

type TagManager struct {
	Db *sql.DB
}

type Tag struct {
	Name        string
	Description string
}

func NewTag(tag_name string, tag_desc string) Tag {
	return Tag{
		Name:        tag_name,
		Description: tag_desc,
	}
}

func (t *TagManager) SerializeTag(tag_name string, tag_desc string) error {
	tx, err := t.Db.Begin()
	if err != nil {
		return err
	}

	tx.Exec(
		"INSERT INTO tags_definition (tag_name, tag_description) VALUES(?, ?);",
		tag_name,
		tag_desc,
	)

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
