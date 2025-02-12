package handler

import (
	"Todo_go/internal/storage/sqlitebase"
	"Todo_go/internal/utils"
	"database/sql"
	"encoding/json"
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
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
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

	router.Get("/", http.HandlerFunc(http.FileServer(http.Dir("C:\\KKO11\\Golang\\Todo_go\\web")).ServeHTTP)) // на маке путь ../web , на  ПК - C:\\KKO11\\Golang\\Todo_go\\web
	router.Get("/api/nextdate", rt.NextDateHandler_Get)
	router.Post("/api/task", rt.AddTaskHandler_Post)
	//chi.Delete("/api/task", .....)
	//chi.Get("/api/task", ....)
	return router
}

// ...
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

// ...
func (rt *Router) AddTaskHandler_Post(w http.ResponseWriter, r *http.Request) {
	var err error
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
		//отправить ошибку в json
		resptask.Error = "ошибка десериализации JSON"

		resp, err := json.Marshal(resptask)
		if err != nil {
			fmt.Printf("ошибка сериализации ошибки: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Printf("ошибка десериализации: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}
	fmt.Printf("Структура: %v\n", task)

	if task.Title == "" {
		resptask.Error = "Не указан заголовок задачи"

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

	var parseDate time.Time
	var nextDate string

	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	} else {
		parseDate, err = time.Parse("20060102", task.Date)
		if err != nil {
			//отправить ошибку в JSON
			resptask.Error = "ошибка формата времени"

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

		if parseDate.Before(time.Now()) {
			switch {
			case task.Repeat == "":
				task.Date = time.Now().Format("20060102")
			default:
				nextDate, err = utils.NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					//отправить ошибку в JSON
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

	resp, err := json.Marshal(resptask)
	if err != nil {
		fmt.Printf("ошибка сериализации номера записи: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}
