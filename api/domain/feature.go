package domain

import "github.com/music-tribe/uuid"

type Feature struct {
	Id          uuid.UUID   `json:"id" bson:"_id" example:"f6e7f8c4-3af6-4028-ac7c-30c9d79a3fa7"`
	UserId      uuid.UUID   `json:"userId" param:"userId" bson:"userId" validate:"required" example:"effe01ec-7f09-4a1c-9453-794212a8ac26"`
	Name        string      `json:"name" validate:"required" example:"My New Feature Request"`
	Description string      `json:"description" validate:"required" example:"Could we have this new feature please?"`
	Votes       []uuid.UUID `json:"votes" example:"['155dccaa-0299-4018-ab6b-90b9ee448943','ef2a27c4-b03d-4190-86f2-b1dc2538243e']"`
}
