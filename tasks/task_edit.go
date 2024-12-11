package tasks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kosttiik/todo-list/ds"
)

func EditTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
			return
		}

		var task ds.Task
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, `{"error":"Ошибка чтения тела запроса"}`, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(body, &task)
		if err != nil {
			http.Error(w, `{"error":"Ошибка преобразования формата"}`, http.StatusBadRequest)
			return
		}

		if task.ID == 0 {
			http.Error(w, `{"error":"Идентификатор задачи не может быть пустым"}`, http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
			return
		}

		now := time.Now().UTC()
		if task.Date == "" {
			task.Date = now.Format(ds.DateFormat)
		}

		parsedDate, err := time.Parse(ds.DateFormat, task.Date)
		if err != nil {
			http.Error(w, `{"error":"Дата указана в неверном формате"}`, http.StatusBadRequest)
			return
		}

		if parsedDate.Before(now) && task.Repeat != "" {
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error":"Ошибка в правиле повторения: %s"}`, err.Error()), http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		}

		query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
		res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
		if err != nil {
			http.Error(w, `{"error":"Ошибка обновления задачи"}`, http.StatusInternalServerError)
			return
		}

		affected, err := res.RowsAffected()
		if err != nil || affected == 0 {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte(`{}`))
	}
}
