package repository

import (
	"time"
)

type DateFields struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

type Model struct {
	ID int `gorm:"primary_key"`
	DateFields
}

type UUIDModel struct {
	ID string `gorm:"default:uuid_generate_v4();primary_key;"`
	DateFields
}

type Interface interface {
	GetID() any
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetDeletedAt() *time.Time
}

func (m Model) GetID() any {
	return m.ID
}

func (m UUIDModel) GetID() any {
	return m.ID
}

func (m DateFields) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func (m DateFields) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}

func (m DateFields) GetDeletedAt() *time.Time {
	return m.DeletedAt
}
