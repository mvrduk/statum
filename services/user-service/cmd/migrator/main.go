package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mvrduk/statum/services/user-service/config"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return errors.New("usage: migrator up | down | version")
	}

	cfg, err := config.Load("config/app.ini")

	if err != nil {
		return err
	}

	dbURL := strings.Replace(cfg.DatabaseURL, "postgres://", "pgx5://", 1)

	m, err := migrate.New("file://migrations", dbURL)

	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}

	defer func() { _, _ = m.Close() }()
	switch os.Args[1] {
	case "up":
		err = m.Up()
	case "down":
		err = m.Steps(-1)
	case "version":
		return printVersion(m)
	default:
		return fmt.Errorf("unknown command %q, use up | down | version", os.Args[1])
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate %s: %w", os.Args[1], err)
	}

	fmt.Println("migrate", os.Args[1], "done")
	return nil
}

func printVersion(m *migrate.Migrate) error {
	version, dirty, err := m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		fmt.Println("no migrations applied")
		return nil
	}

	if err != nil {
		return err
	}

	fmt.Printf("version: %d, dirty: %t\n", version, dirty)
	return nil
}
