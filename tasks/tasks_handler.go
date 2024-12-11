package tasks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kosttiik/todo-list/ds"
)

func TasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")
		limit := 50

		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE 1=1`
		args := []interface{}{}

		if search != "" {
			if parsedDate, err := time.Parse("02.01.2006", search); err == nil {
				query += ` AND date = ?`
				args = append(args, parsedDate.Format(ds.DateFormat))
			} else {
				query += ` AND (title LIKE ? OR comment LIKE ?)`
				searchPattern := "%" + strings.ToLower(search) + "%"
				args = append(args, searchPattern, searchPattern)
			}
		}

		query += ` ORDER BY date ASC LIMIT ?`
		args = append(args, limit)

		rows, err := db.Query(query, args...)
		if err != nil {
			log.Printf("Ошибка выполнения запроса: %s", err.Error())
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка выполнения запроса: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var tasks []ds.Task
		for rows.Next() {
			var task ds.Task
			if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
				log.Printf("Ошибка чтения данных: %s", err.Error())
				http.Error(w, fmt.Sprintf(`{"error":"Ошибка чтения данных: %s"}`, err.Error()), http.StatusInternalServerError)
				return
			}
			tasks = append(tasks, task)
		}

		if err := rows.Err(); err != nil {
			log.Printf("Ошибка итерации по данным: %s", err.Error())
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка итерации по данным: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		if tasks == nil {
			tasks = []ds.Task{}
		}

		response := map[string]interface{}{
			"tasks": tasks,
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(response)
	}
}
