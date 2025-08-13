package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

var CommonDateFormat = "20060102"

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")
	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(CommonDateFormat, nowStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid now date: %v", err), http.StatusBadRequest)
			return
		}
	}
	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, next)
}

func Init(r *chi.Mux) {
	r.Get("/api/nextdate", Auth(nextDayHandler))
	r.Get("/api/tasks", Auth(getTasksHandler))
	r.Get("/api/task", Auth(getTaskHandler))
	r.Post("/api/task", Auth(addTaskHandler))
	r.Post("/api/task/done", Auth(taskDoneHandler))
	r.Post("/api/signin", signinHandler)
	r.Put("/api/task", Auth(updateTaskHandler))
	r.Delete("/api/task", Auth(deleteTaskHandler))
}
