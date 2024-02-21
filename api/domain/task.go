package domain

import "github.com/music-tribe/uuid"

type Task struct {
	Id          uuid.UUID `json:"id" bson:"_id"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Completed   bool      `json:"completed"`
}
