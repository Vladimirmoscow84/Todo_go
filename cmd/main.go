package main

import (
	"Todo_go/internal/storage/sqlitebase"
	"fmt"
	"net/http"

	_ "modernc.org/sqlite"
)

func main() {
	sqlitebase.CheckingBase()
	fmt.Println("Запуск сервера")
	http.Handle("/", http.FileServer(http.Dir("../web")))
	err := http.ListenAndServe(":7540", nil)
	if err != nil {
		panic(err)
	}

}
