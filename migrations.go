package celeritas

import (
	// support for mysql/mariadb
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gobuffalo/pop"
	"github.com/golang-migrate/migrate/v4"

	// support for mysql/mariadb
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	// support for postgres
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// support for file based migrations
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (c *Celeritas) PopConnect() (*pop.Connection, error) {
	tx, err := pop.Connect("development")
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (c *Celeritas) CreatePopMigration(up, down []byte, migrationName, migrationType string) error {
	var migrationPath = c.RootPath + "/migrations"
	err := pop.MigrationCreate(migrationPath, migrationName, migrationType, up, down)
	if err != nil {
		return err
	}
	return nil
}

func (c *Celeritas) RunPopMigrations(tx *pop.Connection) error {
	var migrationPath = c.RootPath + "/migrations"

	fm, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return err
	}

	err = fm.Up()
	if err != nil {
		return err
	}

	return nil
}

func (c *Celeritas) PopMigrateDown(tx *pop.Connection, steps ...int) error {
	var migrationPath = c.RootPath + "/migrations"

	step := 1
	if len(steps) > 0 {
		step = steps[0]
	}

	fm, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return err
	}

	err = fm.Down(step)
	if err != nil {
		return err
	}
	return nil

}

func (c *Celeritas) PopMigrateReset(tx *pop.Connection) error {
	var migrationPath = c.RootPath + "/migrations"
	fm, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return err
	}
	err = fm.Reset()
	if err != nil {
		return err
	}
	return nil
}

// MigrateUp runs an up migration.
func (c *Celeritas) MigrateUp(dsn string) error {
	rootPath := filepath.ToSlash(c.RootPath)
	m, err := migrate.New("file://"+rootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err = m.Up(); err != nil {
		log.Println("Error running migration")
		return err
	}

	return nil
}

// MigrateDownAll runs all down migrations
func (c *Celeritas) MigrateDownAll(dsn string) error {
	rootPath := filepath.ToSlash(c.RootPath)
	m, err := migrate.New("file://"+rootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Down(); err != nil {
		return err
	}

	return nil
}

// Steps runs n migrations. When n is positive, up migrations are run; when negative,
// down migrations are run.
func (c *Celeritas) Steps(n int, dsn string) error {
	rootPath := filepath.ToSlash(c.RootPath)
	m, err := migrate.New("file://"+rootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Steps(n); err != nil {
		return err
	}

	return nil
}

// MigrateForce sets the migration version, and sets the dirty state to false.
func (c *Celeritas) MigrateForce(dsn string) error {
	rootPath := filepath.ToSlash(c.RootPath)
	m, err := migrate.New("file://"+rootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err = m.Force(-1); err != nil {
		log.Println("Error running migration")
		return err
	}

	return nil
}
