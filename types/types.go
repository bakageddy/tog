package types

import "database/sql"

type TogManager struct {
	Db *sql.DB
}

type TogFile struct {
	Id   uint64
	Path string
}

type TogTag struct {
	Id          uint64
	Name        string
	Description string
}

type TogInstance struct {
	File TogFile
	Tags []TogTag
}
