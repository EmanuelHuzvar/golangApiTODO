package models

import "time"

type StateEnum string
type PriorityEnum string

const (
	StateNew                    StateEnum    = "New"
	StateInProgress             StateEnum    = "In Progress"
	StateCompleted              StateEnum    = "Completed"
	PriorityCritical            PriorityEnum = "Critical"
	PriorityHigh                PriorityEnum = "High"
	PriorityMedium              PriorityEnum = "Medium"
	PriorityLow                 PriorityEnum = "Low"
	PriorityCosmeticEnhancement PriorityEnum = "Cosmetic/Enhancement"
)

type User struct {
	Id       int                  `json:"id"`
	Username string               `json:"username"`
	Password string               `json:"password"`
	Tasks    map[StateEnum][]Task `json:"tasks"`
	States   []string             `json:"states"`
}

type Task struct {
	Id                    int          `json:"id"`
	Title                 string       `json:"title"`
	Description           string       `json:"description"`
	State                 StateEnum    `json:"state"`
	Priority              PriorityEnum `json:"priority"`
	CreatedAt             time.Time    `json:"createdAt"`
	StartedAt             *time.Time   `json:"startedAt,omitempty"`
	CompletedAt           *time.Time   `json:"completedAt,omitempty"`
	EstimatedCompletionAt *time.Time   `json:"estimatedCompletionAt"`
	EstimatedStartAt      *time.Time   `json:"estimatedStartAt"`
	Position              int          `json:"position"`
}

type UserForGet struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type TaskPosition struct {
	IdOfTask int `json:"IdOfTask"`
	Position int `json:"Position"`
}

type TaskForGet struct {
	ID                    int        `json:"id"`
	Title                 string     `json:"title"`
	Description           string     `json:"description"`
	Priority              string     `json:"priority"`
	State                 string     `json:"state"`
	CreatedAt             time.Time  `json:"createdAt"`
	StartedAt             *time.Time `json:"startedAt,omitempty"` // Pointer allows for nil value
	CompletedAt           *time.Time `json:"completedAt,omitempty"`
	EstimatedCompletionAt *time.Time `json:"estimatedCompletionAt"`
	EstimatedStartAt      *time.Time `json:"estimatedStartAt"`
	Position              int        `json:"position"`
}
type RemoveTask struct {
	ID int `json:"id"`
}
type SamoUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}
