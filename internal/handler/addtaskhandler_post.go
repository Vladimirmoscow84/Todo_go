package handler

import (
	"Todo_go/internal/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// AddTaskHandler_Post - ручка добавления задачи
func (rt *Router) AddTaskHandler_Post(w http.ResponseWriter, r *http.Request) {
	var err error
	//эземпляр структуры для формирования ответа
	var resptask respTask
	//тело запроса,сформированное в байтовый срез
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	//экземпляр структуры
	var task Task
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	//десериализация тела запроса в структуру
	err = json.Unmarshal(b, &task)
	if err != nil {
		sendError(w, "ошибка десериализации JSON", err)
		return
	}
	fmt.Printf("Структура: %v\n", task)

	if task.Title == "" {
		sendError(w, "не указан заголовок задачи", errors.New("не указан заголовок задачи"))
		return
	}

	var parseDate time.Time
	var nextDate string

	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	} else {
		parseDate, err = time.Parse("20060102", task.Date)
		if err != nil {
			sendError(w, "ошибка формата времени", err)
			return
		}

		if task.Date == time.Now().Format("20060102") {
			task.Date = time.Now().Format("20060102")
		} else if parseDate.Before(time.Now()) {
			switch {
			case task.Repeat == "":
				task.Date = time.Now().Format("20060102")
			default:
				nextDate, err = utils.NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					sendError(w, "ошибка вычисления следующей даты", err)
					return
				}
				task.Date = nextDate
			}
		}

	}
	fmt.Println("Запись в БД...")
	fmt.Printf("date: %s\n", task.Date)
	fmt.Printf("title: %s\n", task.Title)
	fmt.Printf("comment: %s\n", task.Comment)
	fmt.Printf("repeat: %s\n", task.Repeat)
	resDb, err := rt.DB.Exec(`
		INSERT INTO scheduler
			(date, title, comment, repeat)
			VALUES
			(:date, :title, :comment, :repeat)`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		fmt.Printf("ошибка добавления в БД: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Запись в БД выпонена успешно")

	id, err := resDb.LastInsertId()
	if err != nil {
		fmt.Printf("ошибка получения номера записи в БД: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resptask.ID = strconv.Itoa(int(id))
	fmt.Printf("id: %s\n", resptask.ID)

	resp, err := json.MarshalIndent(resptask, "", " ")
	if err != nil {
		fmt.Printf("ошибка сериализации номера записи: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Println(string(resp))
	w.Write(resp)
}
