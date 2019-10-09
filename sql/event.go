package sql

import (
	"database/sql"
	"fmt"
	"github.com/krlvi/github-devstats/event"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Repository struct {
	db       *sql.DB
	migrator *migrate.Migrate
}

func NewRepository(db *sql.DB) (*Repository, error) {
	migrator, err := newMigrator(db)
	if err != nil {
		return nil, err
	}
	repo := &Repository{
		db:       db,
		migrator: migrator,
	}
	return repo, nil
}

func (r *Repository) MigrateUp() error {
	return r.migrator.Up()
}

// Test only
func (r *Repository) migrateDown() error {
	return r.migrator.Down()
}

func newMigrator(db *sql.DB) (*migrate.Migrate, error) {
	var migrationsDir string
	if srcdir, ok := os.LookupEnv("TEST_SRCDIR"); ok {
		migrationsDir = srcdir + "/__main__/sql/migrations"
	} else if wd, ok := os.LookupEnv("BUILD_WORKING_DIRECTORY"); ok {
		migrationsDir = wd + "/sql/migrations"
	}
	log.Printf("Loading migrations from: %s", migrationsDir)
	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	return migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsDir),
		"mysql",
		driver,
	)
}

func (m *Repository) Save(e event.Event) error {
	_, err := m.db.Exec("INSERT INTO events (`repository`, `pr_number`) VALUES ($1, $2)", e.Repository, e.PrNumber)
	if err != nil {
		return err
	}
	return nil
}
