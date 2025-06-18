package tasks

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/kosttiik/todo-list/ds"
)

func GetTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
			return
		}

		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
		row := db.QueryRow(query, id)

		var task ds.Task
		err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, `{"error":"Ошибка получения задачи"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(task)
	}
}
