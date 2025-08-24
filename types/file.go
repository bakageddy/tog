package types

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/bakageddy/tog/util"
)

func (t *TogManager) IsManaged(file string) (bool, error) {
	row := t.Db.QueryRow("SELECT 1 as ROW_COUNT FROM managed_filepaths WHERE filepath = ? LIMIT 1;", file)
	var result uint8
	err := row.Scan(&result);
	if errors.Is(sql.ErrNoRows, err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return result == 1, nil
	}
}

// File must be managed and must exist on the file system
func (t *TogManager) IsPresent(file string) (bool, error) {
	managed, err := t.IsManaged(file)
	if err != nil || !managed {
		return false, err
	}

	if !util.FileExists(file) {
		return false, TogFileDeleted
	}

	return true, nil
}

func (t *TogManager) SearchFile(fileglob string) ([]TogFile, error) {
	if !strings.HasSuffix(fileglob, "*") {
		fileglob += "%"
	}

	query := "SELECT file_id, filepath FROM managed_filepaths WHERE filepath LIKE ?;"
	rows, err := t.Db.Query(query, fileglob)
	if err != nil {
		return nil, err
	}

	result := make([]TogFile, 0)
	for rows.Next() {
		temp := TogFile{}
		err := rows.Scan(&temp.Id, &temp.Path)
		if err != nil {
			return result, err
		}
		result = append(result, temp)
	}

	return result, nil
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
