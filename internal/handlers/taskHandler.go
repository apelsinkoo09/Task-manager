package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/apelsinkoo09/task-manager/internal/models"
	_ "github.com/lib/pq"
)

func GetAllTasksHandler(db *sql.DB) http.HandlerFunc { //db - соединение с базой
	return func(w http.ResponseWriter, r *http.Request) { // хендлер
		//  w http.ResponseWriter - интерфейс для записи ответа клиенту
		//  r *http.Request - структура принимаемого запроса от клиента
		tasks, err := models.ReadAll(db)
		if err != nil {
			http.Error(w, "Unable to retrieve tasks", http.StatusInternalServerError)
			// Сообщение клиенту об ошибке
			// http.StatusInternalServerError - 500 статус
			return
		}
		w.Header().Set("Content-Type", "application/json") // установка заголовка http ответа в формате ключ, значения, запись в карту. Формат отправляемых значений json
		json.NewEncoder(w).Encode(tasks)                   // кодирование в формат json
	}
}

func GetIdTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := r.URL.Query().Get("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			http.Error(w, "Incorrect input", http.StatusBadRequest)
			return
		}
		task, err := models.Read(db, id)
		if err != nil {
			http.Error(w, "Task not exist", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	}
}

func CreateTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { // метод предназначен только для обработки post запросово
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed) // возврат ошибки 405 при несоответствии метода
		}
		var newTask models.Task
		body, err := io.ReadAll(r.Body) // Чтение тела запроса
		if err != nil {
			http.Error(w, "Incrorrect input", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &newTask) // Парсинг данных из json в тело задачи newTask
		if err != nil {
			http.Error(w, "Invalid input format", http.StatusBadRequest)
		}
		err = newTask.Create(db) // newTask имеет тип  Task, для нее определен метод Create, так что можно вызвать метод из под переменной
		if err != nil {
			http.Error(w, "Failed to create task", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTask)
	}
}

func UpdateTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

		idParam := r.URL.Query().Get("id")           // Получение id из строки запроса
		id, err := strconv.ParseInt(idParam, 10, 64) //Преобразование строки в целое число, второ параметр основание системы счисления
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		var updatedTask models.Task

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Incrorrect input", http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(body, &updatedTask)
		if err != nil {
			http.Error(w, "Invalid input format", http.StatusBadRequest)
		}

		updatedTask.Id = id // присвоение id из запроса в структуру задачи
		err = updatedTask.Update(db)
		if err != nil {
			http.Error(w, "Failed to create task", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedTask)
	}
}

func DeleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		idParam := r.URL.Query().Get("id")           // получения id из тела запроса. r.URL.Query() парсит строку запроса, Get("id") извлекает значения id
		id, err := strconv.ParseInt(idParam, 10, 64) //Преобразование строки в число, 10 основание системы счисления, 64 битность(не фурычит)
		if err != nil {
			http.Error(w, "Incorrect input", http.StatusBadRequest)
		}

		var DeleteTask models.Task

		DeleteTask.Id = id
		err = models.Delete(db, id)
		if err != nil {
			http.Error(w, "Failed to delete task", http.StatusInternalServerError)
			return
		}
		w.Header()
	}
}
