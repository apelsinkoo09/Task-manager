package task

import (
	"time"

	_ "github.com/lib/pq"
)

type Task struct {
	Id          int
	Title       string
	Description string
	Status      string
	Priority    int
	Created_at  time.Time
	Updated_at  time.Time
	User_id     int
}

func (task *Task) Save() {

}
