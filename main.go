package main

import (
	"database/sql"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"

	_ "github.com/glebarez/go-sqlite"
)

var (
	file       string
	tag        string
	flag_count uint8 = 0
)

func main() {
	// TODO: I do not know any good default file path
	flag.StringVar(&file, "file", ".", "Set the file path")
	flag.StringVar(&tag, "tag", "default", "Add Tag to the file")

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() != f.DefValue {
			flag_count += 1
		}
	})

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
	log.Println("PARSED: ", flag_count)
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

func add_tag(db *sql.DB, file string, tag string) error {
	if db == nil {
		return errors.New("Uninitialized DB")
	}

	// DEFAULT flag value
	if file == "." {
		return errors.New("Flag `file` must be provided")
	}

	_, err := os.Stat(file)
	if errors.Is(err, os.ErrNotExist) {
		return err
	}

	file, err = filepath.Abs(file)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	tx.Exec("INSERT INTO managed_filepaths(filepath) VALUES (?);", file)
	tx.Exec("INSERT INTO tags_definition(tag_name, tag_description) VALUES (?, ?);", tag, tag_description)

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
