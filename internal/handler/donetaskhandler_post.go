package handler

import (
	"Todo_go/internal/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// DoneTaskHandler_Post - ручка, которая делает задачу выполненой
func (rt *Router) DoneTaskHandler_Post(w http.ResponseWriter, r *http.Request) {

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
	// Проверка на условие повторяемости отмечаемой задачи
	// если задача не повторяется, то ее необходимо удалить
	if task.Repeat == "" {
		_, err := rt.DB.Exec(`
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

	} else {
		//если задача повторяется, то необходимо рассчитать и изменить следующую дату следующего выполнения
		now := time.Now()

		newDate, err := utils.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		task.Date = newDate

		//вносим изменения в таблицу
		_, err = rt.DB.Exec(`
		UPDATE scheduler
		SET date = :date
		WHERE id = :id 
		`,
			sql.Named("date", task.Date),
			sql.Named("id", task.ID),
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

}
