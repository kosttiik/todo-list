package tasks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/kosttiik/todo-list/ds"
)

func AddTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task ds.Task
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Ошибка чтения тела запроса: %v", err)
			http.Error(w, `{"error":"Ошибка чтения тела запроса"}`, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(body, &task)
		if err != nil {
			log.Printf("Ошибка преобразования формата: %v", err)
			http.Error(w, `{"error":"Ошибка преобразования формата"}`, http.StatusBadRequest)
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
			log.Printf("Дата указана в неверном формате: %v", err)
			http.Error(w, `{"error":"Дата указана в неверном формате"}`, http.StatusBadRequest)
			return
		}

		parsedDate = NormalizeToDate(parsedDate)
		now = NormalizeToDate(now)

		if parsedDate.Before(now) && task.Repeat == "" {
			task.Date = now.Format(ds.DateFormat)
		}

		if parsedDate.Before(now) && task.Repeat != "" {
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				log.Printf("Ошибка в правиле повторения: %v", err)
				http.Error(w, fmt.Sprintf(`{"error":"Ошибка в правиле повторения: %s"}`, err), http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		}

		query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
		res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
		if err != nil {
			log.Printf("Ошибка добавления задачи: %v", err)
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка добавления задачи: %s"}`, err), http.StatusInternalServerError)
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			log.Printf("Ошибка получения ID задачи: %v", err)
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка получения ID задачи: %s"}`, err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	}
}
