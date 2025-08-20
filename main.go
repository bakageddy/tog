package main

import (
	"database/sql"
	"flag"
	"io"
	"log"
	"os"

	"github.com/bakageddy/tog/types"
	"github.com/bakageddy/tog/util"
	_ "github.com/glebarez/go-sqlite"
)

var (
	file string
	tag  string
	cmd  string
)

func main() {
	// FIX:
	// TODO: I do not know any good default file path
	flag.StringVar(&file, "file", ".", "Set the file path")
	flag.StringVar(&tag, "tag", "default", "Add Tag to the file")
	flag.StringVar(&cmd, "cmd", "add", util.CommandDescription)

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	db, err := sql.Open("sqlite", "tags.db")
	if err != nil {
		log.Fatalln("failed to open database: ", err)
	}
	defer db.Close()

	err = database_init(db, "./schema.sql")
	if err != nil {
		log.Fatalln(err.Error())
	}

	flag.Parse()

	tfm := types.TogManager{Db: db}
	cmd_type := util.Mux(cmd)
	switch cmd_type {
	case util.AddFile:
		tfm.ManageFile(file)
	case util.RemoveFile:
	tfm.ReleaseFile(file)	
	}
}

func database_init(db *sql.DB, schema_path string) error {
	file, err := os.Open(schema_path)
	if err != nil {
		return err
	}
	defer file.Close()
	body, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(string(body))
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
