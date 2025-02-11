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

	chi.Get("/", http.HandlerFunc(http.FileServer(http.Dir("C:\\KKO11\\Golang\\Todo_go\\web")).ServeHTTP)) // на маке путь нужно откорректировать ../web
	chi.Get("/api/nextdate", rt.NextDateHandler_Get)
	return chi
}

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
