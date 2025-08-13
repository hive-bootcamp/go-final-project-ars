package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Arsadidas/go-final-project/pkg/db"
	"github.com/Arsadidas/go-final-project/pkg/helpers"
	"github.com/Arsadidas/go-final-project/pkg/models"
)

type TasksResp struct {
	Tasks []*models.Task `json:"tasks"`
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(w, r, 50)
	if err != nil {
		helpers.WriteJSONError(w, "cannot get tasks by limit")
		return
	}

	if tasks == nil {
		tasks = []*models.Task{}
	}

	helpers.WriteJson(w, map[string][]*models.Task{"tasks": tasks})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Query().Get("id")
	task, err := db.GetTask(w, taskId)
	if err != nil {
		helpers.WriteJSONError(w, "cannot get task by id")
		return
	}
	helpers.WriteJson(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var task models.Task
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		helpers.WriteJSONError(w, fmt.Sprintf("Ошибка чтения тела запроса: %v", err))
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		helpers.WriteJSONError(w, fmt.Sprintf("Ошибка при десериализации тела запроса: %v", err))
		return
	}
	if task.Title == "" {
		helpers.WriteJSONError(w, "Не указан заголовок задачи")
		return
	}
	if err := checkDate(&task); err != nil {
		helpers.WriteJSONError(w, fmt.Sprintf("Неверный формат даты: %v", err))
		return
	}
	db.UpdateTask(w, &task)
	helpers.WriteJson(w, struct{}{})
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Query().Get("id")
	task, err := db.GetTask(w, taskId)
	if err != nil {
		helpers.WriteJSONError(w, "cannot get task by id")
		return
	}

	if task.Repeat == "" {
		db.DeleteTask(w, taskId)
	} else {
		nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			helpers.WriteJSONError(w, "cannot get next date")
			return
		}
		db.UpdateDate(w, nextDate, taskId)
	}
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Query().Get("id")
	err := db.DeleteTask(w, taskId)
	if err != nil {
		helpers.WriteJSONError(w, "cannot delete task by id")
		return
	}
	helpers.WriteJson(w, struct{}{})
}
