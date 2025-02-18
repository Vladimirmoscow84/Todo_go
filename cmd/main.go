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

	fmt.Println("Запуск сервера")

	err = http.ListenAndServe(":7540", router.Routers())
	if err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s\n", err.Error())
		return
	}

	router.DB.Close()
	time.Sleep(time.Second * 3)
}
