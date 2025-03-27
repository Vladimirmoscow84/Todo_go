package sqlitebase

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// CheckingBase проверяет существование в дирректории файла scheduler.db
func CheckingBase() (*sql.DB, error) {
	appPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("ошибка executable: %s", err.Error())
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	fmt.Println(appPath)
	fmt.Println(dbFile)
	_, err = os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}
	var db *sql.DB
	if install {
		db, err = sql.Open("sqlite3", dbFile)
		if err != nil {
			return nil, fmt.Errorf("ошибка открытия БД: %s", err.Error())
		}
		fmt.Println("Подключено к БД")

		// defer db.Close()

		_, err = db.Exec("CREATE TABLE scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date TEXT NOT NULL DEFAULT '', title VARCHAR(128) NOT NULL DEFAULT '', comment VARCHAR(256) NOT NULL DEFAULT '', repeat VARCHAR(128) NOT NULL DEFAULT '')")
		if err != nil {
			return nil, fmt.Errorf("ошибка создания таблицы scheduler: %s", err.Error())
		}
		_, err = db.Exec("CREATE INDEX date_sort ON scheduler (date)")
		if err != nil {
			return nil, fmt.Errorf("ошибка создания индексов таблицы scheduler: %s", err.Error())
		}

		fmt.Println("Созданы таблицы и индексы")
	} else {
		db, err = sql.Open("sqlite3", dbFile)
		if err != nil {
			return nil, fmt.Errorf("ошибка открытия БД: %s", err.Error())
		}
		fmt.Println("Подключено к БД")
	}
	return db, nil
}
