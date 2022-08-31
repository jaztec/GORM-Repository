package repository

import (
	"context"
)

type DBProvider interface {
	DB(ctx context.Context) Database
}

type Migrator interface {
	Migrate(ctx context.Context) error
}

type Repository[T Interface] struct {
	db       Database
	preloads map[string][]any
}

func (r *Repository[T]) GetByID(ctx context.Context, id string) (T, error) {
	var e T
	err := r.addPreloads(r.DB(ctx)).
		Where("id = ?", id).
		First(&e)
	return e, err
}

func (r *Repository[T]) FindAll(ctx context.Context) ([]T, error) {
	var e []T
	err := r.addPreloads(r.DB(ctx)).
		Find(&e)
	return e, err
}

func (r *Repository[T]) Create(ctx context.Context, e *T) (*T, error) {
	err := r.db.DB(ctx).Create(e)
	return e, err
}

func (r *Repository[T]) Update(ctx context.Context, e *T) (*T, error) {
	err := r.DB(ctx).Save(e)
	return e, err
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

	err := tx.Find(&ts)
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func (r *Repository[T]) Model() []Interface {
	var e T
	return []Interface{e}
}

func (r *Repository[T]) DB(ctx context.Context) Database {
	return r.db.DB(ctx)
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

func (r *Repository[T]) addPreloads(tx Database) Database {
	for p, args := range r.preloads {
		tx.Preload(p, args...)
	}
	return tx
}

func NewRepository[T Interface](db Database) (Repository[T], error) {
	return Repository[T]{db: db}, nil
}

type MigrateUtil struct {
	db     Database
	models []any
}

func NewMigrateUtil(db Database, models []Interface) Migrator {
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
	return e.db.DB(ctx).AutoMigrate(e.models...)
}
