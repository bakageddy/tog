package util

import "errors"

type TogError error

var (
	TogFileExists   TogError = errors.New("File Exists/Managed")
	TogFileNotFound TogError = errors.New("File Not Found")
	TogErrDatabase  TogError = errors.New("File Database Error")
)
