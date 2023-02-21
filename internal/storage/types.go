package storage

import (
	"time"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/utils"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	PasswordHash []byte    `json:"-"`
	PasswordSalt string    `json:"-"`
	Login        string    `json:"login"`
	IsAdmin      bool      `json:"is_admin"`
	LastOnline   time.Time `json:"last_online"`
}

type Company struct {
	ID         uuid.UUID    `json:"id"`
	Owner      uuid.UUID    `json:"owner_id"`
	Name       string       `json:"name"`
	TimeOffset utils.Offset `json:"time_offset"`
}

type Device struct {
	ID               uuid.UUID `json:"id"`
	Company          uuid.UUID `json:"company_id"`
	Name             string    `json:"name"`
	Resource         int       `json:"resource"`
	MinutesRemaining int       `json:"minutes_remaining"`
	LastOnline       time.Time `json:"last_online"`
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
