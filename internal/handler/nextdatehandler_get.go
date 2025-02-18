package handler

import (
	"Todo_go/internal/utils"
	"net/http"
	"time"
)

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
