package domain

type Task struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Completed   bool   `json:"completed"`
}
