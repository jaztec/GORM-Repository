package repository_test

import (
	"context"
	repository "github.com/jaztec/gorm-repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

type testModel struct {
	repository.Model

	Name string
}

func TestCRUDCommands(t *testing.T) {
	db := getDb(t)
	r, err := repository.NewRepository[testModel](db)
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
		if ne.ID == "" {
			t.Errorf("No ID was filled")
		}
		if ne.Name != "Test" {
			t.Errorf("Name does not comply: %s", ne.Name)
		}
	})
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
