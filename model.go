package repository

import (
	"time"
)

type Model struct {
	ID        string `gorm:"default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

type Interface interface {
	GetID() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetDeletedAt() *time.Time
}

func (m Model) GetID() string {
	return m.ID
}

func (m Model) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func (m Model) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}

func (m Model) GetDeletedAt() *time.Time {
	return m.DeletedAt
}
