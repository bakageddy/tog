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

func (t *TagManager) IsPresent(tag_name string) (bool, error) {
	row := t.Db.QueryRow("SELECT 1 FROM tags_definition WHERE tag_name = ?", tag_name)
	var result uint8
	if err := row.Scan(&result); err != nil {
		return false, err
	}
	return result == 1, nil
}

func (t *TagManager) SerializeTag(tag_name string, tag_desc string) error {
	tx, err := t.Db.Begin()
	if err != nil {
		return err
	}

	present, err := t.IsPresent(tag_name)
	if err != nil {
		return err
	}

	if present {
		return TogTagExists
	}

	_, err = tx.Exec(
		"INSERT INTO tags_definition (tag_name, tag_description) VALUES(?, ?);",
		tag_name,
		tag_desc,
	)

	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (t *TagManager) SearchTag(tag_name string) ([]Tag, error) {
	rows, err := t.Db.Query(
		"SELECT tag_name, tag_description FROM tags_definition WHERE tag_name LIKE ?",
		tag_name,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag = make([]Tag, 1)

	for rows.Next() {
		var tag_name, tag_description string
		var tag Tag = Tag{}
		err = rows.Scan(&tag_name, &tag_description)
		if err != nil {
			return tags, err
		}

		tag.Name = tag_name
		tag.Description = tag_description
		tags = append(tags, tag)
	}

	return tags, nil
}

func (t *TagManager) GetTag(tag_name string) (Tag, error) {
	row := t.Db.QueryRow(
		"SELECT tag_name, tag_description FROM tag_description WHERE tag_name = ?;",
		tag_name,
	)

	var tag_result_name, tag_description string
	if err := row.Scan(&tag_name, &tag_description); err != nil {
		return Tag{}, err
	}

	if tag_result_name != tag_name {
		return Tag{}, TogTagNotFound
	} else {
		return Tag{Name: tag_name, Description: tag_description}, nil
	}
}

func (t *TagManager) GetTagID(tag Tag) (uint64, error) {
	row := t.Db.QueryRow("SELECT tag_id FROM tags_definition WHERE tag_name = ?;", tag.Name)
	var result uint64
	err := row.Scan(&result)
	if err != nil {
		return 0, err
	}
	return result, nil
}
