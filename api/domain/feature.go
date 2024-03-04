package domain

import "github.com/music-tribe/uuid"

type Feature struct {
	Id          uuid.UUID `json:"id" bson:"_id"`
	UserId      uuid.UUID `json:"userId" param:"userId" bson:"userId" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Completed   bool      `json:"completed"`
}
