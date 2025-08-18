package util

import (
	"database/sql"
	"errors"
)

type TogFileManager struct {
	Db *sql.DB
}

type TogFile struct {
	File string
	Tags []Tag
}

func (t *TogFileManager) SerializeFile(file string) error {
	present, err := t.IsPresent(file)
	if err != nil {
		return err
	}

	if present {
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

func (t *TogFileManager) GetManagedFile(file string) (TogFile, error) {
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
		file,
	)

	if err != nil {
		return TogFile{}, err
	}
	defer rows.Close()

	tog_file := TogFile{
		File: file,
		Tags: make([]Tag, 0),
	}

	for rows.Next() {
		var file_path string
		var tag_name string
		rows.Scan(&file_path, &tag_name)

		tag := Tag{
			Name:        tag_name,
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

func (t *TogFileManager) GetFileID(file string) (uint64, error) {
	row := t.Db.QueryRow("SELECT file_id FROM managed_filepaths WHERE file_path = ?;", file)
	var result uint64
	if err := row.Scan(&result); err != nil {
		return 0, errors.Join(err, TogFileNotFound)
	}
	return result, nil
}

func (t *TogFileManager) IsPresent(file string) (bool, error) {
	row := t.Db.QueryRow("SELECT 1 FROM managed_filepaths WHERE file_path = ?;", file)
	var result int
	if err := row.Scan(&result); err != nil {
		return false, err
	}
	return result == 1, nil
}

func (t *TogFileManager) AssociateTag(file string, tags []Tag) error {
	file_id, err := t.GetFileID(file)
	if err != nil {
		return err
	}

	tag_mgr := TagManager{
		Db: t.Db,
	}

	tx, err := tag_mgr.Db.Begin()
	if err != nil {
		return err
	}

	// PERF: check whether batching or discrete db calls is better
	for _, tag := range tags {
		tag_id, err := tag_mgr.GetTagID(tag)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.Exec("INSERT INTO file_tags (file_id, tag_id) VALUES(?, ?);", file_id, tag_id)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return nil
}
