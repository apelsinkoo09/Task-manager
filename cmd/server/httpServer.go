package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/apelsinkoo09/task-manager/internal/handlers"
	_ "github.com/lib/pq"
)

type User struct {
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

var db *sql.DB

func connectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Чтение данных из запроса
	var conectionInfo User
	err := json.NewDecoder(r.Body).Decode(&conectionInfo)
	if err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	//Строка подключения
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", conectionInfo.User, conectionInfo.Password, conectionInfo.DBName)
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}

	// Проверка подключения
	if err := db.Ping(); err != nil {
		log.Printf("failed to ping database: %v", err)
		return
	}
	// Отправка ответа клиенту
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully connected to the database!"))
}

func main() {
	// Настройка маршрутов
	http.HandleFunc("/connect", connectHandler)
	http.HandleFunc("/api/v1/tasks", handlers.GetAllTasksHandler(db))
	http.HandleFunc("/api/v1/task", handlers.GetIdTaskHandler(db))
	http.HandleFunc("/api/v1/task/create", handlers.CreateTaskHandler(db))
	http.HandleFunc("/api/v1/task/update", handlers.UpdateTaskHandler(db))
	http.HandleFunc("/api/v1/task/delete", handlers.DeleteTaskHandler(db))

	// Запуск сервера
	log.Println("Server is running on port 8081...")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
