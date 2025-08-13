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

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var task models.Task

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		helpers.WriteJSONError(w, fmt.Sprintf("Ошибка чтения тела запроса: %v", err))
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		helpers.WriteJSONError(w, fmt.Sprintf("Ошибка при десериализации тела запроса: %v", err))
		return
	}
	if task.Title == "" {
		helpers.WriteJSONError(w, "Не указан заголовок задачи")
		return
	}
	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}
	if err := checkDate(&task); err != nil {
		helpers.WriteJSONError(w, fmt.Sprintf("Неверный формат даты: %v", err))
		return
	}
	id, err := db.AddTask(&task)
	if err != nil {
		helpers.WriteJSONError(w, fmt.Sprintf("Ошибка добавления задачи: %v", err))
		return
	}
	helpers.WriteJson(w, map[string]int64{"id": id})
}

func checkDate(task *models.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(CommonDateFormat)
	}
	parsedDate, err := time.Parse(CommonDateFormat, task.Date)
	if err != nil {
		return err
	}
	var next string
	var nextDateErr error
	if task.Repeat != "" {
		next, nextDateErr = NextDate(now, task.Date, task.Repeat)
	}
	if nextDateErr != nil {
		return nextDateErr
	}
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if parsedDate.Before(todayStart) {
		if task.Repeat == "" {
			task.Date = now.Format(CommonDateFormat)
		} else {
			task.Date = next
		}
	}
	return nil
}
