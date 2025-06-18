package database

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/kosttiik/todo-list/ds"
	_ "github.com/mattn/go-sqlite3"
)

func CreateDB() *sql.DB {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		appPath, err := os.Getwd()
		if err != nil {
			log.Fatalf("Не удалось получить текущую директорию приложения: %v", err)
		}
		dbFile = filepath.Join(appPath, "scheduler.db")
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Ошибка при открытии БД: %v", err)
	}

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		createTableQuery := `
        CREATE TABLE IF NOT EXISTS scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date VARCHAR(8) NOT NULL,
            title TEXT NOT NULL,
            comment TEXT,
            repeat VARCHAR(128)
        );
        CREATE INDEX IF NOT EXISTS indexDate ON scheduler(date);
        `
		if _, err := db.Exec(createTableQuery); err != nil {
			log.Fatalf("Ошибка создания таблицы: %v", err)
		}
		log.Println("База данных создана")
	} else {
		log.Println("База данных уже существует")
	}

	return db
}

func GetTaskByID(db *sql.DB, id string) (ds.Task, error) {
	var task ds.Task

	getTaskQuery := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	row := db.QueryRow(getTaskQuery, id)

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if errors.Is(err, sql.ErrNoRows) {
		return ds.Task{}, errors.New("Задача не найдена")
	} else if err != nil {
		return ds.Task{}, err
	}

	return task, nil
}

func DeleteTaskFromDB(db *sql.DB, id string) error {
	deleteTaskQuery := `DELETE FROM scheduler WHERE id = ?`

	_, err := db.Exec(deleteTaskQuery, id)

	return err
}

func UpdateTaskInDB(db *sql.DB, task ds.Task) error {
	updateTaskQuery := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	_, err := db.Exec(updateTaskQuery, task.Date, task.Title, task.Comment, task.Repeat, task.ID)

	return err
}
