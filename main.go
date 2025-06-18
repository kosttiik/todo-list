package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kosttiik/todo-list/database"
	"github.com/kosttiik/todo-list/tasks"
)

func main() {
	webDir := "./web"

	db := database.CreateDB()
	defer db.Close()

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", tasks.HandleNextDate)

	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			tasks.AddTaskHandler(db)(w, r)
		case http.MethodPut:
			tasks.EditTaskHandler(db)(w, r)
		case http.MethodGet:
			tasks.GetTaskHandler(db)(w, r)
		case http.MethodDelete:
			tasks.DeleteTaskHandler(db)(w, r)
		default:
			http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/tasks", tasks.TasksHandler(db))
	http.HandleFunc("/api/task/done", tasks.DoneTaskHandler(db))

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = ":7540"
	}

	log.Printf("Сервер запущен с портом %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
