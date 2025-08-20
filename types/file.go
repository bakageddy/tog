package types

import (
	"github.com/bakageddy/tog/util"
	"golang.org/x/tools/present"
)

// File must be managed and must exist on the file system
func (t *TogManager) IsPresent(file string) (bool, error) {
	row := t.Db.QueryRow("SELECT 1 FROM managed_filepaths WHERE file_path = ?;", file)
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
	present, err := t.IsPresent(file)
	if err != nil {
		return err
	}
}

func (t *TogManager) GetFile(file string) (TogFile, error) {
}

func (t *TogManager) AssociateTag(file []string, tags TogTag) error {
}
