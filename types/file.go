package types

import (
	"errors"

	"github.com/bakageddy/tog/util"
)

// File must be managed and must exist on the file system
func (t *TogManager) IsPresent(file string) (bool, error) {
	// FIX: SELECT 1 is not getting in to the resultset
	row := t.Db.QueryRow("SELECT 1 FROM managed_filepaths WHERE filepath = ?;", file)
	var result int
	if err := row.Scan(&result); err != nil {
		return false, err
	}

	if result != 1 {
		return false, TogFileNotManaged
	}

	if !util.FileExists(file) {
		return false, TogFileDeleted
	}

	return true, nil
}

func (t *TogManager) GetFile(file string) (TogFile, error) {
	present, err := t.IsPresent(file)
	if err != nil {
		return TogFile{}, err
	}

	if !present {
		return TogFile{}, TogUnreachable
	}

	result := TogFile{}
	row := t.Db.QueryRow("SELECT file_id, filepath FROM managed_filepaths WHERE filepath = ?", file)
	if err := row.Scan(&result.Id, &result.Path); err != nil {
		return result, err
	}
	return result, nil
}

func (t *TogManager) ManageFile(file string) error {
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

func (t *TogManager) ReleaseFile(file string) error {
	tx, err := t.Db.Begin()
	if err != nil {
		return err
	}

	togfile, err := t.GetFile(file)
	if err != nil && errors.Is(err, TogFileNotManaged) {
		return err
	}

	_, err = tx.Exec("DELETE * FROM managed_filepaths WHERE file_id = ?;", togfile.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE * FROM file_tags WHERE file_id = ?;", togfile.Id)
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

func (t *TogManager) AssociateTag(files []TogFile, tag TogTag) error {
	tx, err := t.Db.Begin()
	if err != nil {
		return err
	}

	query := "INSERT INTO file_tags (file_id, tag_id) VALUES (?, ?);"
	for _, file := range files {
		if _, err := tx.Exec(query, file.Id, tag.Id); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (t *TogManager) DisassociateTag(files []TogFile, tag TogTag) error {
	tx, err := t.Db.Begin()
	if err != nil {
		return err
	}

	query := "DELETE FROM file_tags WHERE file_id = ? AND tag_id = ?;"
	for _, file := range files {
		if _, err := tx.Exec(query, file.Id, tag.Id); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
