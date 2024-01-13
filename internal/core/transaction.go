package core

import (
	"sync"

	"github.com/google/uuid"
)

type Transaction struct {
	Id     uuid.UUID
	Mu     sync.RWMutex
	Writer string
	Logs   []uuid.UUID
}
