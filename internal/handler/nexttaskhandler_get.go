package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// NextTaskHandler_Get ручка для получения списка ближайших задач
func (rt *Router) NextTaskHandler_Get(w http.ResponseWriter, r *http.Request) {
	//эземпляр структуры для формирования ответа при ошибке
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var resptask respTask
	tasks := make([]Task, 0)
	now := time.Now().Format("20060102")
	//обращение к БД
	rows, err := rt.DB.Query(`
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE date >= :now
		ORDER BY date
	`, sql.Named("now", now),
	)
	if err != nil {
		resptask.Error = err.Error()

		resp, err := json.Marshal(resptask)
		if err != nil {
			fmt.Printf("ошибка сериализации ошибки: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return

	}
	defer rows.Close()

	//вызов функции НЕКСТ() для парсинга множества строк из БД построчно
	for rows.Next() {
		tempTask := Task{}

		err := rows.Scan(&tempTask.ID, &tempTask.Date, &tempTask.Title, &tempTask.Comment, &tempTask.Repeat)
		if err != nil {
			resptask.Error = err.Error()

			resp, err := json.Marshal(resptask)
			if err != nil {
				fmt.Printf("ошибка сериализации ошибки: %s\n", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(resp)
			return
		}
		tasks = append(tasks, tempTask)
	}

	if err := rows.Err(); err != nil {
		resptask.Error = err.Error()

		resp, err := json.Marshal(resptask)
		if err != nil {
			fmt.Printf("ошибка сериализации ошибки: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}
	mapa := make(map[string][]Task)
	mapa["tasks"] = tasks
	resp, err := json.MarshalIndent(mapa, "", " ")
	if err != nil {
		fmt.Printf("ошибка сериализации ответа: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}
