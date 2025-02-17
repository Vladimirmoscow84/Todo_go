package handler

import (
	"Todo_go/internal/storage/sqlitebase"
	"Todo_go/internal/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"

	_ "github.com/mattn/go-sqlite3"
)

type Router struct {
	DB *sql.DB
}

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// структура для вывода ответа при добавлении новой задачи
type respTask struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// NewRouter - контсруктор роутера
func NewRouter() (*Router, error) {

	db, err := sqlitebase.CheckingBase()
	if err != nil {
		return nil, err
	}
	return &Router{
		DB: db,
	}, nil
}

func (rt *Router) Routers() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/api/nextdate", rt.NextDateHandler_Get)
	router.Post("/api/task", rt.AddTaskHandler_Post)
	router.Get("/api/tasks", rt.NextTaskHandler_Get)
	router.Get("/api/task", rt.TaskIDhandler_Get)
	router.Put("/api/task", rt.ChangeTaskHandler_Put)

	router.Handle("/*", http.FileServer(http.Dir("C:\\KKO11\\Golang\\Todo_go\\web"))) // на маке путь ../web , на  ПК - C:\\KKO11\\Golang\\Todo_go\\web
	// router.Handle("/", http.FileServer(http.Dir("../web"))) // на маке путь ../web , на  ПК - C:\\KKO11\\Golang\\Todo_go\\web
	return router
}

// sendError - сериализация и отправка ошибки в формате JSON
func sendError(w http.ResponseWriter, errText string, err error) {
	var resptaskErr respTask
	resptaskErr.Error = errText
	fmt.Printf("%s: %s\n", errText, err.Error())

	resp, err := json.Marshal(resptaskErr)
	if err != nil {
		fmt.Printf("ошибка сериализации ошибки: %s\n", err.Error())
		http.Error(w, fmt.Sprintf("%s: %s", errText, err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write(resp)
}

// NextDateHandler_Get - ручка для возврата даты следующего выполнения задачи
func (rt *Router) NextDateHandler_Get(w http.ResponseWriter, r *http.Request) {

	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	parseNow, err := time.Parse("20060102", now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := utils.NextDate(parseNow, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(response))

}

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

// NextTaskHandler_Get ручка для получения списка ближайших задач
func (rt *Router) NextTaskHandler_Get(w http.ResponseWriter, r *http.Request) {
	//эземпляр структуры для формирования ответа при ошибке
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var resptask respTask
	tasks := make([]Task, 0)
	now := time.Now().Format("20060102")
	//обращение к БД
	rows, err := rt.DB.Query(`
		SELECT date, title, comment, repeat 
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

		err := rows.Scan(&tempTask.Date, &tempTask.Title, &tempTask.Comment, &tempTask.Repeat)
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

// TaskIDhandler_Get  - ручка для получения задач по id
func (rt *Router) TaskIDhandler_Get(w http.ResponseWriter, r *http.Request) {

	idStr := r.FormValue("id") //значение параметра id  в строке запроса api/task?id=<...>
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, "Задача не найдена", err)
		return
	}

	task := Task{}

	//обращение к БД
	row := rt.DB.QueryRow(`
		SELECT id, date, title, comment, repeat FROM scheduler 
		WHERE id == :idInt
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

// ChangeTaskHandler_Put - ручка для изменения значений задачь
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

	// ro := rt.DB.QueryRow(`
	// 	SELECT title FROM scheduler
	// 	WHERE id = :id
	// `, sql.Named("id", idInt))
	// var f string
	// ro.Scan(&f)
	// if f == "" {
	// 	sendError(w, "ошибка: id отсутствует в базе", err)
	// 	return
	// }

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
