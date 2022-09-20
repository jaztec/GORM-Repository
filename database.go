package repository

import (
	"context"
	"gorm.io/gorm"
)

type Database interface {
	DB(context.Context) Database
	Preload(query string, args ...any) Database
	Where(query string, args ...any) Database
	Joins(query string, args ...any) Database
	Limit(limit int) Database
	Offset(offset int) Database
	Order(value any) Database

	Create(value any) error
	Find(dest any, conds ...any) error
	First(dest any, conds ...any) error
	Save(value any) error
	Delete(value any, conds ...any) error

	AutoMigrate(dst ...any) error
}

type gormDBImpl struct {
	db *gorm.DB
}

func (impl *gormDBImpl) DB(ctx context.Context) Database {
	return impl.wrap(impl.db.WithContext(ctx))
}

func (impl *gormDBImpl) Preload(query string, args ...any) Database {
	return impl.wrap(impl.db.Preload(query, args...))
}

func (impl *gormDBImpl) Where(query string, args ...any) Database {
	return impl.wrap(impl.db.Where(query, args...))
}

func (impl *gormDBImpl) Joins(query string, args ...any) Database {
	return impl.wrap(impl.db.Joins(query, args...))
}

func (impl *gormDBImpl) Limit(limit int) Database {
	return impl.wrap(impl.db.Limit(limit))
}

func (impl *gormDBImpl) Offset(offset int) Database {
	return impl.wrap(impl.db.Offset(offset))
}

func (impl *gormDBImpl) Order(value any) Database {
	return impl.wrap(impl.db.Order(value))
}

func (impl *gormDBImpl) Delete(value any, conds ...any) error {
	return impl.db.Delete(value, conds...).Error
}

func (impl *gormDBImpl) Create(value any) error {
	return impl.db.Create(value).Error
}

func (impl *gormDBImpl) Find(dest any, conds ...any) error {
	return impl.db.Find(dest, conds...).Error
}

func (impl *gormDBImpl) First(dest any, conds ...any) error {
	return impl.db.First(dest, conds...).Error
}

func (impl *gormDBImpl) Save(value any) error {
	return impl.db.Save(value).Error
}

func (impl *gormDBImpl) AutoMigrate(dst ...any) error {
	return impl.db.AutoMigrate(dst...)
}

func (impl *gormDBImpl) wrap(db *gorm.DB) Database {
	return &gormDBImpl{db: db}
}

func NewGORMDatabase(db *gorm.DB) Database {
	return &gormDBImpl{db: db}
}
