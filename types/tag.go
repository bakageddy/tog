package types

import (
	"database/sql"
	"errors"
)

func (t *TogManager) IsTagPresent(tag_name string) (bool, error) {
	row := t.Db.QueryRow(
		"SELECT 1 FROM tags_definition WHERE tag_name = ?",
		tag_name,
	)

	var result uint8
	err := row.Scan(&result)
	if errors.Is(sql.ErrNoRows, err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return result == 1, nil
	}
}

// Creates an Instance on the Database
func (t *TogManager) NewTag(tag_name string, tag_desc string) error {
	tx, err := t.Db.Begin()
	if err != nil {
		return err
	}

	present, err := t.IsTagPresent(tag_name)
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

func (t *TogManager) FetchTag(tag_name string) (TogTag, error) {
	present, err := t.IsTagPresent(tag_name)
	if err != nil {
		return TogTag{}, err
	}

	if !present {
		return TogTag{}, TogTagNotFound
	}

	result := TogTag{}
	row := t.Db.QueryRow("SELECT tag_id, tag_name, tag_description FROM tags_definition WHERE tag_name = ? LIMIT 1", tag_name)
	if err := row.Scan(&result.Id, &result.Name, &result.Description); err != nil {
		return result, err
	}
	return result, nil
}

func (t *TogManager) RemoveTag(tag_name string) error {
	tag, err := t.FetchTag(tag_name)
	if err != nil {
		return err
	}

	tx, err := t.Db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM tags_definition WHERE tag_id = ?", tag.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM file_tags WHERE tag_id = ?", tag.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (t *TogManager) SearchTag(tag_glob string) ([]TogTag, error) {
	rows, err := t.Db.Query(
		"SELECT tag_name, tag_description FROM tags_definition WHERE tag_name LIKE ?",
		tag_glob,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []TogTag = make([]TogTag, 0)

	for rows.Next() {
		var tag_name, tag_description string
		var tag TogTag = TogTag{}
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
