package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

// Хранение конфига БД из файла считанного из файла .yaml
type Config struct {
	Database struct {
		Host     string // Адрес хоста бд
		Port     int    // Порт бд
		User     string
		Password string
		DBName   string
	}
}

// Функция для чтения конфигурационного файла
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var configuration Config

	err = yaml.Unmarshal(data, &configuration)
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}

// Функция для подключения к бд, передача в метод указатель!!! на пул соединений
func DBConnection(c Config) (*sql.DB, error) {
	// Формирование строки подключения на основе данных из конфига
	connectionString := fmt.Sprintf("user=%s password=%s dbname='%s' host=%s port=%d sslmode=disable",
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.Host,
		c.Database.Port)

	//fmt.Println("Connection string:", connectionString)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
	}

	// Проверка подключения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	fmt.Println("Successful connection to the database")
	return db, nil
}

func main() {
	config, err := LoadConfig("configs/config.yaml")
	if err != nil {
		log.Printf("Failed to load configuration: %v", err)
	}
	db, err := DBConnection(*config)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
	}
	defer db.Close()
}
