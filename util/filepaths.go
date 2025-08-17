package util

import (
	"database/sql"
)

type TogFileManager struct {
	Db *sql.DB
}

type TogFile struct {
	File string
	Tags []Tag
}

func (t *TogFileManager) NewFile(file string) error {
	is_present, err := t.IsPresent(file)
	if err != nil {
		return err
	}

	if is_present {
		return TogFileExists
	}

	tx, err := t.Db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO managed_filepaths(filepath) VALUES(?)", file)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (t *TogFileManager) GetFile(file string) (TogFile, error) {
	present, err := t.IsPresent(file)
	if err != nil {
		return TogFile{}, err
	}

	if !present {
		return TogFile{}, TogFileNotFound
	}

	tx, err := t.Db.Begin()
	if err != nil {
		return TogFile{}, err
	}

	rows, err := tx.Query(`
		SELECT filepath, tag_name 
		FROM managed_filepaths m
		JOIN file_tags f ON m.file_id = f.file_id
		JOIN tags_definition t on f.tag_id = t.tag_id
		WHERE filepath = ?;`, 
		file)

	if err != nil {
		return TogFile{}, err
	}
	defer rows.Close()


	tog_file := TogFile{
		File: file,
		Tags: make([]Tag, 0),
	}

	for ; rows.Next(); {
		var file_path string
		var tag_name string
		rows.Scan(&file_path, &tag_name)

		tag := Tag {
			Name: tag_name,
			Description: "",
		}

		tog_file.Tags = append(tog_file.Tags, tag)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return TogFile{}, err
	}

	return tog_file, nil
}

func (t *TogFileManager) IsPresent(file string) (bool, error) {
	tx, err := t.Db.Begin()
	if err != nil {
		return false, err
	}

	row := tx.QueryRow("SELECT 1 FROM managed_filepaths WHERE file_path = ?", file)
	var result int
	if err := row.Scan(&result); err != nil {
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return false, err
	}

	return result == 1, nil
}

func (t *TogFileManager) AssociateTag(file string, tags []Tag) error {
	return nil
}
