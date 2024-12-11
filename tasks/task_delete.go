package tasks

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/kosttiik/todo-list/database"
)

func DeleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
			return
		}

		_, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, `{"error":"Неверный формат идентификатора задачи"}`, http.StatusBadRequest)
			return
		}

		err = database.DeleteTaskFromDB(db, id)
		if err != nil {
			http.Error(w, `{"error":"Ошибка удаления задачи"}`, http.StatusInternalServerError)
			return
		}

		w.Write([]byte("{}"))
	}
}
