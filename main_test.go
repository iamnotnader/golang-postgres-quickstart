package main

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func connectToTestDB(t *testing.T) (*sqlx.DB, *migrate.Migrate) {
	require.NoError(t, godotenv.Load("test.env"))

	psqlURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))
	fmt.Printf("DB URL: %s\n", psqlURL)

	db, err := sqlx.Open("postgres", psqlURL)
	require.NoError(t, err)

	// Give the db some time to come up...
	dbStarted := false
	pingStartTime := time.Now()
	for (time.Now().Sub(pingStartTime)) < 1*time.Second {
		err = db.Ping()
		if err == nil {
			dbStarted = true
			fmt.Printf("Connected to db after %v millis\n",
				time.Now().Sub(pingStartTime).Milliseconds())
			break
		}
		time.Sleep(time.Millisecond * 10)
	}
	require.True(t, dbStarted,
		"Error could not connect to database after 1 second of pinging...")

	// Run migrations
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Error starting sql migration: %v", err)
	}
	mm, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", "migrations"), // file://path/to/directory
		"postgres", driver)
	require.NoError(t, err, "Error migration failed")

	err = mm.Up()
	require.Truef(t, err == nil || err == migrate.ErrNoChange,
		"Error occurred while syncing db: %v", err)

	return db, mm
}

func tearDownTestDB(t *testing.T, db *sqlx.DB, mm *migrate.Migrate) {
	err := mm.Down()
	require.NoError(t, err, "Error reversing migrations")
}

func TestDB(t *testing.T) {
	db, mm := connectToTestDB(t)
	defer db.Close()

	db.MustExec("INSERT INTO users (linkedin_url) VALUES ($1)", "bbb")

	tearDownTestDB(t, db, mm)
}
