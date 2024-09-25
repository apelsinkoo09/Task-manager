package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Task struct {
	Id          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    int       `json:"priority"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
	User_id     int       `json:"user_id"`
}

func ReadAll(db *sql.DB) ([]Task, error) {
	result, err := db.Query("select * from tasks") //Query возвращает итерратор, надо пройти по всем!!! строчкам
	if err != nil {
		log.Println(err)
	}
	defer result.Close()

	var tasks []Task

	for result.Next() { // next() подготавливает следующую строку для чтения с помощью метода scan()

		var task Task

		err := result.Scan( // проверка на ошибки в конкретной строке. Scan сканирует строки и назначает им переменные соответствующего типа
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.Created_at,
			&task.Updated_at,
			&task.User_id,
		)
		if err != nil {
			log.Printf("failed to scan row: %v", err)
		}
		tasks = append(tasks, task)
		log.Printf("Task: %+v\n", task)
	}
	err = result.Err() // поиск общих ошибок
	if err != nil {
		log.Println(err)
	}
	return tasks, nil
}

func Read(db *sql.DB, id int64) (Task, error) {
	request := "select * from tasks where id = $1"
	result := db.QueryRow(request, id)
	var task Task
	err := result.Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.Created_at,
		&task.Updated_at,
		&task.User_id,
	)
	if err != nil {
		log.Println(err)
	}
	return task, nil
}

func (task *Task) Create(db *sql.DB) error {
	insert := `
	insert into tasks (title, description, status, priority, created_at)
	values ($1, $2, $3, $4, $5)
	`
	_, err := db.Exec(insert, task.Title, task.Description, task.Status, task.Priority, time.Now())
	if err != nil {
		return fmt.Errorf("failed to create task: %v", err)
	}
	return nil
}

func (task *Task) Update(db *sql.DB) error {
	update := `
	Update tasks 
	Set title = $1, description = $2, status = $3, priority = $4, updated_at = $5
	Where id = $6
	`
	_, err := db.Exec(update, task.Title, task.Description, task.Status, task.Priority, time.Now(), task.Id)
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}
	return nil
}

func Delete(db *sql.DB, id int64) error {
	del := "Delete from tasks Where id = $1"
	_, err := db.Exec(del, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}
	return nil
}
