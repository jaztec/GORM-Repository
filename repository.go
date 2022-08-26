package repository

import (
	"context"
	"gorm.io/gorm"
)

type DBProvider interface {
	DB(ctx context.Context) *gorm.DB
}

type Migrator interface {
	Migrate(ctx context.Context) error
}

type Repository[T Interface] struct {
	db       *gorm.DB
	preloads map[string][]any
}

func (r *Repository[T]) GetByID(ctx context.Context, id string) (T, error) {
	var e T
	tx := r.addPreloads(r.DB(ctx)).
		Where("id = ?", id).
		First(&e)
	return e, tx.Error
}

func (r *Repository[T]) FindAll(ctx context.Context) ([]T, error) {
	var e []T
	tr := r.addPreloads(r.DB(ctx)).
		Find(&e)
	return e, tr.Error
}

func (r *Repository[T]) Create(ctx context.Context, e *T) (*T, error) {
	tr := r.db.WithContext(ctx).Create(e)
	return e, tr.Error
}

func (r *Repository[T]) Update(ctx context.Context, e *T) (*T, error) {
	tr := r.DB(ctx).Save(e)
	return e, tr.Error
}

func (r *Repository[T]) FindBy(ctx context.Context, after, pageSize int, conditions ...Condition) ([]T, error) {
	var ts []T
	tx := r.DB(ctx).
		Limit(pageSize).
		Offset(after).
		Order("created_at ASC")

	tx = r.addPreloads(tx)

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

func (r *Repository[T]) Model() []Interface {
	var e T
	return []Interface{e}
}

func (r *Repository[T]) DB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *Repository[T]) AddPreload(preload string, args []any) *Repository[T] {
	if r.preloads == nil {
		r.preloads = make(map[string][]any, 5)
	}
	r.preloads[preload] = args

	return r
}

func (r *Repository[T]) ClearPreloads() *Repository[T] {
	r.preloads = nil

	return r
}

func (r *Repository[T]) addPreloads(tx *gorm.DB) *gorm.DB {
	for p, args := range r.preloads {
		tx.Preload(p, args)
	}
	return tx
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
