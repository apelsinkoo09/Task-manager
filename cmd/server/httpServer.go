package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
	log.Println("Connection succesfull")

	// Настройка максимального количества открытых соединений
	db.SetMaxOpenConns(25)

	// Настройка максимального количества свободных соединений
	db.SetMaxIdleConns(25)

	// Настройка максимального времени ожидания перед разрывом соединения
	db.SetConnMaxLifetime(5 * time.Minute)

	// Отправка ответа клиенту
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully connected to the database!"))
}

func main() {
	// Настройка маршрутов
	http.HandleFunc("/connect", connectHandler)
	http.HandleFunc("/api/v1.1/tasks", func(w http.ResponseWriter, r *http.Request) { // обязательно проверять на существование подключения к бд, хандлер вызывать в теле маршрутизатора
		if db == nil {
			http.Error(w, "Database not connected", http.StatusInternalServerError)
			return
		}
		handlers.GetAllTasksHandler(db)(w, r)
	})
	http.HandleFunc("/api/v1.1/task", func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			http.Error(w, "Database not connected", http.StatusInternalServerError)
			return
		}
		handlers.GetIdTaskHandler(db)(w, r)
	})
	http.HandleFunc("/api/v1.1/task/create", func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			http.Error(w, "Database not connected", http.StatusInternalServerError)
			return
		}
		handlers.CreateTaskHandler(db)(w, r)
	})
	http.HandleFunc("/api/v1.1/task/update", func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			http.Error(w, "Database not connected", http.StatusInternalServerError)
			return
		}
		handlers.UpdateTaskHandler(db)(w, r)
	})
	http.HandleFunc("/api/v1.1/task/delete", func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			http.Error(w, "Database not connected", http.StatusInternalServerError)
			return
		}
		handlers.DeleteTaskHandler(db)(w, r)
	})

	// Запуск сервера
	log.Println("Server is running on port 8081...")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
