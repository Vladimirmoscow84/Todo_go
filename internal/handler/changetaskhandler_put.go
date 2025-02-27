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

// ChangeTaskHandler_Put - ручка для изменения значений задач
func (rt *Router) ChangeTaskHandler_Put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	task := Task{}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//десериализация тела запроса в структуру
	err = json.Unmarshal(b, &task)
	if err != nil {
		sendError(w, "ошибка десериализации JSON", err)
		return
	}

	// конверитирование значения строчного поля id структуры task в цифровое, чтобы оперировать в БД
	idInt, err := strconv.Atoi(task.ID)
	if err != nil {
		sendError(w, "ошибка в конвертированиии поля ID", err)
		return
	}
	fmt.Printf("id: %d\n", idInt)
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

	// проверка наличия id в БД
	ro := rt.DB.QueryRow(`
		SELECT id FROM scheduler
		WHERE id = :id
	`, sql.Named("id", idInt))

	var f int
	err = ro.Scan(&f)
	if err != nil {
		sendError(w, "ошибка: id отсутствует в базе", err)
		return
	}

	// обновление записи  в БД
	_, err = rt.DB.Exec(`
		UPDATE scheduler
		SET date = :date,
			title = :title,
			comment = :comment,
			repeat = :repeat
			WHERE id = :id`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", idInt),
	)
	if err != nil {
		sendError(w, "ошибка добавления в БД: %s\n", err)
		return
	}
	fmt.Println("Запись в БД выполнена успешно")
	var emptytask Task

	resp, err := json.Marshal(emptytask)
	if err != nil {
		fmt.Printf("ошибка сериализации ответа: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(resp)

}
