package main

import (
	"context"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/spanner"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"github.com/reiot777/spansqlx-example/spanadmin"
)

func init() {
	parser.AddCommand("migrate",
		"Migrate spanner database (beta)",
		"The migrate command (beta) will schema to spanner database",
		&migrateCommand{},
	)
}

type migrateCommand struct {
	Database string `short:"d" long:"database" description:"spanner emulator database url" default:"projects/sandbox/instances/sandbox/databases/sandbox" ENV:"DATABASE"`
	File     string `short:"f" long:"file" description:"Full path migrate file" default:"migrations"`
}

func (cmd *migrateCommand) Execute(args []string) error {
	if os.Getenv("SPANNER_EMULATOR_HOST") != "" {
		if err := spanadmin.CreateInstanceAndDatabase(context.TODO(), cmd.Database, false); err != nil {
			return errors.Wrap(err, "cannot create spanner instance and database.")
		}
	}

	if err := migrateFiles(cmd.Database, cmd.File); err != nil {
		return errors.Wrap(err, "apply migrate")
	}

	fmt.Println("Migration successful!")
	return nil
}

func migrateFiles(uri string, path string) error {
	s := &spanner.Spanner{}
	d, err := s.Open(uri + "?x-clean-statements=true")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+path, uri, d)
	if err != nil {
		return err
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		// Already up-tp-date
		return nil
	}

	return err
}
