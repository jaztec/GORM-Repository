package repository_test

import (
	"context"
	repository "github.com/jaztec/gorm-repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"
)

type testModel struct {
	repository.Model

	Name string
}

func TestCRUDCommands(t *testing.T) {
	db := getDb(t)
	r, err := repository.NewRepository[testModel](repository.NewGORMDatabase(db))
	if err != nil {
		t.Fatalf("Error creating DB: %v", err)
	}

	t.Run("Create one", func(t *testing.T) {
		n := testModel{
			Name: "Test",
		}

		ne, err := r.Create(context.Background(), &n)
		if err != nil {
			t.Errorf("Error creating one: %+v", err)
		}
		if ne.GetID() == "" {
			t.Errorf("No ID was filled")
		}
		if ne.Name != "Test" {
			t.Errorf("Name does not comply: %s", ne.Name)
		}
	})

	t.Run("Find one", func(t *testing.T) {
		var err error
		models := []testModel{
			{Name: "Test"},
			{Name: "Test2"},
		}

		for _, n := range models {
			_, err = r.Create(context.Background(), &n)
			if err != nil {
				t.Fatalf("Error creating one: %+v", err)
			}
		}

		m, err := r.FindBy(context.Background(), 1, 1, repository.NewWhereCondition("name = ?", "Test"))
		if err != nil {
			t.Fatalf("Error finding one: %+v", err)
		}
		if len(m) != 1 {
			t.Errorf("Invalid number of records returned, expect %d but received %d", 1, len(m))
		}
		ne := m[0]
		testModelParts(t, ne, false)
	})
}

func testModelParts(t *testing.T, m testModel, withDeleted bool) {
	if m.GetID() == "" {
		t.Errorf("No ID was filled")
	}
	if m.GetCreatedAt().Equal(time.Time{}) {
		t.Errorf("Incorrect create time")
	}
	if m.GetUpdatedAt().Equal(time.Time{}) || !m.GetUpdatedAt().Equal(m.GetCreatedAt()) {
		t.Errorf("Incorrect update time")
	}
	if m.GetDeletedAt() == nil && withDeleted {
		t.Errorf("Incorrect delete time")
	} else if m.GetDeletedAt() != nil && !withDeleted {
		t.Errorf("Incorrect delete time")
	}
	if m.GetID() == "" {
		t.Errorf("No ID was filled")
	}
	if m.Name != "Test" {
		t.Errorf("Name does not comply: %s", m.Name)
	}
}

func getDb(t *testing.T) *gorm.DB {
	db := postgres.Open("postgres://test:test@localhost:5432")
	gdb, err := gorm.Open(db, &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	var m testModel
	if err := gdb.Migrator().AutoMigrate(&m); err != nil {
		t.Fatal(err)
	}

	return gdb
}
