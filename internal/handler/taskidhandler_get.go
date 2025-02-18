package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// TaskIDhandler_Get  - ручка для получения задач по id
func (rt *Router) TaskIDhandler_Get(w http.ResponseWriter, r *http.Request) {

	qp := r.URL.Query()
	idStr := qp.Get("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, "Задача не найдена", err)
		return
	}

	task := Task{}

	//обращение к БД
	row := rt.DB.QueryRow(`
		SELECT id, date, title, comment, repeat FROM scheduler 
		WHERE id = :idInt
		`, sql.Named("idInt", idInt),
	)
	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		sendError(w, "ошибка БД", err)
		return
	}
	resp, err := json.MarshalIndent(task, "", " ")
	if err != nil {
		fmt.Printf("ошибка сериализации ответа: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(resp)
}
