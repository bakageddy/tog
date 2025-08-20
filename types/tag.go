package types

func (t *TogManager) IsTagPresent(tag_name string) (bool, error) {
	row := t.Db.QueryRow(
		"SELECT 1 FROM tags_definition WHERE tag_name = ?",
		tag_name,
	)

	var result uint8
	if err := row.Scan(&result); err != nil {
		return false, err
	}
	return result == 1, nil
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

func (t *TogManager) SearchTag(tag_name string) ([]TogTag, error) {
	rows, err := t.Db.Query(
		"SELECT tag_name, tag_description FROM tags_definition WHERE tag_name LIKE ?",
		tag_name,
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
