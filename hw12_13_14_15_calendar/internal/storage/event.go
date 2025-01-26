package storage

import (
	"errors"
	"time"
)

var (
	ErrEventAlreadyExist = errors.New("event with this id already exist")
	ErrEventDoesNotExist = errors.New("event does not exist")
)

type Event struct {
	ID          string        `json:"id" db:"id"`
	Title       string        `json:"title" db:"title"`
	StartDate   time.Time     `json:"startDate" db:"start_date"`
	Duration    time.Duration `json:"duration" db:"duration"`
	Description string        `json:"description" db:"description"`
	Owner       string        `json:"owner" db:"owner"`
	RemindAt    int64         `json:"remindAt" db:"remind_at"`
}
