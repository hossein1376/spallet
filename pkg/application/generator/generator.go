package generator

import (
	"github.com/google/uuid"
)

type Generate struct {
}

func New() Generate {
	return Generate{}
}

func (g Generate) NewUUID() uuid.UUID {
	return uuid.New()
}
