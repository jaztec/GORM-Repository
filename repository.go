package repository

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DBProvider interface {
	DB(ctx context.Context) *gorm.DB
}

type Migrator interface {
	Migrate(ctx context.Context) error
}

type Repository[T Interface] struct {
	db *gorm.DB
}

func (r *Repository[T]) GetByID(ctx context.Context, id string) (T, error) {
	var e T
	result := r.DB(ctx).
		Preload(clause.Associations).
		Where("id = ?", id).
		First(&e)
	return e, result.Error
}

func (r *Repository[T]) FindAll(ctx context.Context) ([]T, error) {
	var e []T
	tr := r.DB(ctx).
		Preload(clause.Associations).
		Find(&e)
	return e, tr.Error
}

func (r *Repository[T]) Create(ctx context.Context, e *T) (*T, error) {
	tr := r.db.WithContext(ctx).Create(e)
	return e, tr.Error
}

func (r *Repository[T]) FindBy(ctx context.Context, after, pageSize int, conditions ...Condition) ([]T, error) {
	var ts []T
	tx := r.DB(ctx).
		Limit(pageSize).
		Offset(after).
		Order("created_at ASC")

	for _, c := range conditions {
		switch c.Type() {
		case TypeWhere:
			tx = tx.Where(c.Query(), c.Args()...)
		case TypeJoin:
			tx = tx.Joins(c.Query(), c.Args()...)
		}
	}

	tx.Find(&ts)

	if err := tx.Error; err != nil {
		return nil, err
	}

	return ts, nil
}

func (r *Repository[T]) Model() T {
	var e T
	return e
}

func (r *Repository[T]) DB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func NewRepository[T Interface](db *gorm.DB) (Repository[T], error) {
	return Repository[T]{db: db}, nil
}

type MigrateUtil struct {
	db     *gorm.DB
	models []any
}

func NewMigrateUtil(db *gorm.DB, models []Interface) Migrator {
	m := make([]any, 0, len(models))
	for _, i := range models {
		m = append(m, i)
	}
	return &MigrateUtil{
		db:     db,
		models: m,
	}
}

func (e *MigrateUtil) Migrate(ctx context.Context) error {
	return e.db.WithContext(ctx).AutoMigrate(e.models...)
}
