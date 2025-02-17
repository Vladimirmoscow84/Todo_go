package main

import (
	"Todo_go/internal/handler"
	"fmt"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

func main() {

	router, err := handler.NewRouter()
	if err != nil {
		fmt.Printf("Ошибка роутера: %s\n", err.Error())
		return
	}

	// для windows
	// http.Handle("/", http.FileServer(http.Dir("C:\\KKO11\\Golang\\Todo_go\\web")))
	// для mac
	// http.Handle("/", http.FileServer(http.Dir("../web"))) // на маке путь нужно откорректировать ../web

	fmt.Println("Запуск сервера")

	err = http.ListenAndServe(":7540", router.Routers())
	if err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s\n", err.Error())
		return
	}

	router.DB.Close()
	time.Sleep(time.Second * 3)

}
