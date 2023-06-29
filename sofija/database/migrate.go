package database

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Bloxico/exchange-gateway/sofija/log"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type MigrationProcess struct {
	db *DB

	logger log.Logger
}

func NewMigrationProcess(db *DB, logger log.Logger) *MigrationProcess {
	return &MigrationProcess{
		db:     db,
		logger: logger,
	}
}

// Performs all migrations from the given directory
func (m *MigrationProcess) Migrate(migrationsDir string) (err error) {
	m.logger.Info("running migrations",
		"dir", migrationsDir)

	// Find all migration files in the given dir
	files, err := findMigrationFiles(migrationsDir)
	if err != nil {
		return errors.Wrap(err, "find migrations")
	}

	return m.db.TxContext(context.Background(), func(ctx context.Context) error {
		for _, file := range files {
			m.logger.Info("executing migration",
				"name", file.Name())

			tx, _ := FromContext(ctx)
			// Execute each migration, one by one
			_, err = sqlx.LoadFile(tx, filepath.Join(migrationsDir, file.Name()))
			if err != nil {
				return errors.Wrap(err, "execute migration")
			}
		}

		m.logger.Info("migration done")
		return nil
	})
}

// Given a path to the folder containing all the migration files, return the list of all migrations in the directory
func findMigrationFiles(dir string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "read migrations directory")
	}

	var migrationFiles []os.FileInfo
	for _, fi := range files {
		if fi.IsDir() || filepath.Ext(fi.Name()) != ".sql" {
			continue
		}

		migrationFiles = append(migrationFiles, fi)
	}

	if len(migrationFiles) == 0 {
		return nil, errors.New("no migrations found")
	}

	return migrationFiles, nil
}

// Drops the given schema from the database
func (m *MigrationProcess) DropSchema(name string) error {
	m.logger.Info("dropping schema",
		"schema", name)

	if name == "" {
		return errors.New("schema name empty")
	}

	_, err := m.db.Exec(context.Background(), fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", name))
	if err != nil {
		return errors.Wrap(err, "drop schema")
	}

	m.logger.Info("schema dropped")
	return nil
}
