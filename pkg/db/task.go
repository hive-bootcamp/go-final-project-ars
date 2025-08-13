package db

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Arsadidas/go-final-project/pkg/helpers"
	"github.com/Arsadidas/go-final-project/pkg/models"
)

func AddTask(task *models.Task) (int64, error) {
	var id int64
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES(:date, :title, :comment, :repeat)"
	result, err := Db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err == nil {
		id, err = result.LastInsertId()
	}

	return id, err
}

func Tasks(w http.ResponseWriter, r *http.Request, limit int) ([]*models.Task, error) {
	var tasks []*models.Task
	searchValue := r.URL.Query().Get("search")

	var isDate bool
	var parsedDate time.Time
	var err error
	queryArgs := []any{}
	var query string = ""

	if searchValue != "" {
		parsedDate, err = time.Parse("02.01.2006", searchValue)
		if err == nil {
			isDate = true
		}
	} else {
		query = "SELECT * FROM scheduler ORDER BY date ASC LIMIT :limit"
	}

	if isDate {
		formattedDate := parsedDate.Format("20060102")
		query = "SELECT * from scheduler WHERE date = :date ORDER BY date ASC LIMIT :limit"
		queryArgs = append(queryArgs, sql.Named("date", formattedDate))
	} else {
		query = "SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date ASC LIMIT :limit"
		queryArgs = append(queryArgs, sql.Named("search", "%"+searchValue+"%"))
	}
	queryArgs = append(queryArgs, sql.Named("limit", limit))

	rows, err := Db.Query(query, queryArgs...)
	if err != nil {
		helpers.WriteJSONError(w, "cannot get tasks by limit")
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		task := models.Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			helpers.WriteJSONError(w, "cannot scan values of task columns")
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func GetTask(w http.ResponseWriter, id string) (*models.Task, error) {
	task := &models.Task{}
	row := Db.QueryRow("SELECT * from scheduler WHERE id = :id", sql.Named("id", id))

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		helpers.WriteJSONError(w, "cannot get task")
		return nil, err
	}
	return task, err
}

func UpdateTask(w http.ResponseWriter, task *models.Task) {
	query := `
		UPDATE scheduler
		SET date = :date, title = :title, comment = :comment, repeat = :repeat
		WHERE id = :id
		`
	result, err := Db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID),
	)
	if err != nil {
		helpers.WriteJSONError(w, "cannot update task")
		return
	}
	count, err := result.RowsAffected()
	if err != nil {
		helpers.WriteJSONError(w, "cannot get updated tasks")
		return
	}
	if count == 0 {
		helpers.WriteJSONError(w, "incorrect id for updating task")
		return
	}
}

func UpdateDate(w http.ResponseWriter, next string, id string) {
	query := `
		UPDATE scheduler
		SET date = :date
		WHERE id = :id
		`
	result, err := Db.Exec(query, sql.Named("date", next), sql.Named("id", id))
	if err != nil {
		helpers.WriteJSONError(w, "cannot update date")
		return
	}
	count, err := result.RowsAffected()
	if err != nil {
		helpers.WriteJSONError(w, "cannot get updated date")
		return
	}
	if count == 0 {
		helpers.WriteJSONError(w, "incorrect id for updating task date")
		return
	}
}

func DeleteTask(w http.ResponseWriter, id string) error {
	_, err := Db.Exec("DELETE from scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		helpers.WriteJSONError(w, "cannot delete task")
		return err
	}

	return nil
}
