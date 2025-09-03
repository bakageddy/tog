package types

import "errors"

type TogError error

var (
	TogUnreachable         TogError = errors.New("Unreachable")
	TogUnrecognizedCommand TogError = errors.New("Unrecognized Command")

	TogFileExists     TogError = errors.New("File Exists")
	TogFileNotFound   TogError = errors.New("File Not Found")
	TogFileNotManaged TogError = errors.New("File Not Managed")
	TogFileDeleted    TogError = errors.New("Managed File Deleted")

	TogTagExists   TogError = errors.New("Tag Exists/Managed")
	TogTagNotFound TogError = errors.New("Tag Not Found")
	TogErrDatabase TogError = errors.New("File Database Error")
)
