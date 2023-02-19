package storage

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	PasswordHash []byte
	PasswordSalt string
	Login        string
	IsAdmin      bool
	LastOnline   time.Time
}

type Company struct {
	ID         uuid.UUID
	Owner      uuid.UUID
	Name       string
	TimeOffset time.Duration
}

type Device struct {
	ID         uuid.UUID
	Company    uuid.UUID
	Name       string
	Resource   int
	LastOnline time.Time
}

type Schedule struct {
	ID        uuid.UUID
	Device    uuid.UUID
	Week      int
	TimeStart int
	TimeStop  int
}

type Journal struct {
	ID     uuid.UUID
	Device uuid.UUID
	Start  time.Time
	End    time.Time
	Done   bool
}
