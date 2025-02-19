package handler

import (
	"Todo_go/internal/storage/sqlitebase"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

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
	root, _ := os.Getwd()
	router.Get("/api/nextdate", rt.NextDateHandler_Get)
	router.Post("/api/task", rt.AddTaskHandler_Post)
	router.Get("/api/tasks", rt.NextTaskHandler_Get)
	router.Get("/api/task", rt.TaskIDhandler_Get)
	router.Put("/api/task", rt.ChangeTaskHandler_Put)
	router.Post("/api/task/done", rt.DoneTaskHandler_Post)
	router.Delete("/api/task", rt.DeleteTaskHandler_Delete)

	router.Handle("/*", http.FileServer(http.Dir(filepath.Join(root, "web"))))

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
