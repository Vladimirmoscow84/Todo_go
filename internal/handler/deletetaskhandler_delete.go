package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// DeleteTaskHandler_Delete -   ручка, которая удаляет задачу по ее  id
func (rt *Router) DeleteTaskHandler_Delete(w http.ResponseWriter, r *http.Request) {

	qp := r.URL.Query()
	idStr := qp.Get("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, "Задача не найдена", err)
		return
	}

	//обращение к БД (удаление по указанному id)
	_, err = rt.DB.Exec(`
			DELETE FROM scheduler
			WHERE id = :idInt
			`, sql.Named("idInt", idInt),
	)

	if err != nil {
		sendError(w, "ошибка БД", err)
		return
	}

	var emptyJSON respTask

	resp, err := json.Marshal(emptyJSON)
	if err != nil {
		fmt.Printf("ошибка сериализации ответа: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(resp)
}
