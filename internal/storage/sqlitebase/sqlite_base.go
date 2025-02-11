package sqlitebase

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// CheckingBase проверяет существование в дирректории файла scheduler.db
func CheckingBase() (*sql.DB, error) {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}
	var db *sql.DB
	if install {
		db, err := sql.Open("sqlite", appPath+"scheduler.db")
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		// defer db.Close()

		_, err = db.Exec("CREATE TABLE scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date TEXT NOT NULL DEFAULT '', title VARCHAR(128) NOT NULL DEFAULT '', comment VARCHAR(256) NOT NULL DEFAULT '', repeat VARCHAR(128) NOT NULL DEFAULT '')")
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		_, err = db.Exec("CREATE INDEX date_sort ON scheduler (date)")
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

	}
	return db, nil
}
