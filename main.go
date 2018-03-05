/*
Tests
$ docker run --rm --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -d postgres
$ docker cp MBrand e5de5c99802a:/home/
$ docker container exec -it e5de5c99802a bash
$ ./FirstSetup.sh localhost menable
*/
package main

import (
	// other imports

	"database/sql"
	_ "github.com/lib/pq" // "empty import" to register driver with database/sql
	"net/url"
	"testing"
	"time"
)

type DBTest func(t *testing.T, conn *sql.DB, db *PostgresDB)

func RunDBTest(t *testing.T, dbVersion string, test DBTest) {
	/*
		Password string
		Username string // defaults to "postgres"
		Database string // defaults to "username"
		Version  string // defaults to "latest"
	*/
	c := PostgresConfig{"pass", "postgres", "postgres", dbVersion}

	// create a postgres container
	db, err := NewPostgresDB(c)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close() // destroy the postgres container after the test

	// create a connection URL
	// http://www.postgresql.org/docs/current/static/libpq-connect.html#LIBPQ-CONNSTRING
	connURL := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.Username, c.Password),
		Host:     db.Host,
		Path:     "/" + c.Database,
		RawQuery: "sslmode=disable",
	}

	// connect to database
	conn, err := sql.Open("postgres", connURL.String())
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	// ping the database every 100ms until it comes up
	// wait max 20 sec till return Error
	timeout := time.Now().Add(time.Second * 20)
	for time.Now().Before(timeout) {
		err = conn.Ping()
		if err == nil {
			// yay! we've connected to the database, time to run the test
			test(t, conn, db)
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Errorf("failed to connect to database: %v", err)
}
