package repository_test

import (
	"context"
	repository "github.com/jaztec/gorm-repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strconv"
	"testing"
	"time"
)

type testMasterModel struct {
	repository.Model

	Name    string
	Details []testDetailModel `gorm:"foreignKey:MasterID;references:ID"`
}

type testDetailModel struct {
	repository.Model

	Name     string
	MasterID *int
	Master   *testMasterModel `gorm:"foreignKey:MasterID"`
}

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

func TestRelationDatabaseCommands(t *testing.T) {
	db := getDb(t)
	r, err := repository.NewRepository[testMasterModel](repository.NewGORMDatabase(db))
	if err != nil {
		t.Fatalf("Error creating DB: %v", err)
	}

	t.Run("Make sure preloads are executed", func(t *testing.T) {
		m := testMasterModel{Name: "master"}
		_, _ = r.Create(context.Background(), &m)
		for i := 1; i < 4; i++ {
			db.Create(&testDetailModel{
				Name:   "detail " + strconv.Itoa(i),
				Master: &m,
			})
		}

		r.AddPreload("Details", nil)

		q, err := r.FindBy(context.Background(), 0, 1, repository.NewWhereCondition("ID = ?", m.ID))
		if err != nil {
			t.Fatalf("Error fetching master record: %v", err)
		}

		if len(q) != 1 {
			t.Errorf("Expect exactly 1 result, got: %d", len(q))
			return
		}

		q1 := q[0]
		if q1.Details == nil {
			t.Error("Details should have been preloaded")
			return
		}

		if len(q1.Details) != 3 {
			t.Errorf("Expect exactly 3 results, got: %d", len(q1.Details))
		}

		r.ClearPreloads()
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
	db := sqlite.Open("gorm-test.db")
	gdb, err := gorm.Open(db, &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	var m testModel
	var tm testMasterModel
	var td testDetailModel
	if err := gdb.Migrator().AutoMigrate(&m, &tm, &td); err != nil {
		t.Fatal(err)
	}

	return gdb
}
