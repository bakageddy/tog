package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
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

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile | log.Ldate)

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
	case util.AddFile: {
		tfm.IsManaged(file)
		err := tfm.ManageFile(file)
		if errors.Is(types.TogFileNotManaged, err) {
			fmt.Fprintf(os.Stderr, "%s not managed by tog, consider adding it\n", file)
		} else if errors.Is(types.TogFileDeleted, err) {
			// TODO: fix this issue, add isPresent, isManaged, isDeleted
			fmt.Fprintf(os.Stderr, "%s seems to be deleted/not found, consider removing it\n", file)
		} else if err != nil {
			log.Println(err.Error())
		} else {
			log.Printf("Managing %s\n", file)
		}
	}
		
	case util.RemoveFile: {
		err := tfm.ReleaseFile(file)	
		if errors.Is(types.TogFileNotManaged, err) {
			fmt.Fprintf(os.Stderr, "%s not managed by tog", file)
		} else if errors.Is(types.TogFileDeleted, err) {
			fmt.Fprintf(os.Stderr, "%s already deleted", file)
		} else if err != nil {
			log.Println(err.Error())
		}
	}
	case util.SearchFile:
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
