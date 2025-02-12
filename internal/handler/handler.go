package handler

import (
	"Todo_go/internal/storage/sqlitebase"
	"Todo_go/internal/utils"
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi"
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
	chi := chi.NewRouter()

	chi.Get("/", http.HandlerFunc(http.FileServer(http.Dir("../web")).ServeHTTP)) // на маке путь ../web , на  ПК - C:\\KKO11\\Golang\\Todo_go\\web
	chi.Get("/api/nextdate", rt.NextDateHandler_Get)
	chi.Post("/api/task", rt.AddTaskHandler_Post)
	//chi.Delete("/api/task", .....)
	//chi.Get("/api/task", ....)
	return chi
}

// ...
func (rt *Router) NextDateHandler_Get(w http.ResponseWriter, r *http.Request) {

	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	parseNow, err := time.Parse("20060102", now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	response, err := utils.NextDate(parseNow, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Write([]byte(response))

}

// ...
func (rt *Router) AddTaskHandler_Post(w http.ResponseWriter, r *http.Request) {

}
