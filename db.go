// Licensed under GPL, 2016
// Refer to LICENSE for more details
// Refer to README for structural information

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

const entries = 14 // entries per page, should be even

var (
	db                         *sql.DB
	create, getsol, getone     *sql.Stmt
	subm, getsolc, create_word *sql.Stmt
	findn, findp, findd        *sql.Stmt
	findnr, findpr, finddr     *sql.Stmt
	getan, getap, getad        *sql.Stmt
	getrnd, getwrd, inccou     *sql.Stmt
)

func init() {
	if os.Getenv("redb_conn") == "" {
		log.Fatal("\"$redb_conn\" not defined.")
	}

	var err error
	db, err = sql.Open("postgres", os.Getenv("redb_conn"))
	if err != nil {
		log.Fatal(err)
	}

	mustExec := func(stmt string) {
		_, err := db.Exec(stmt)
		if err != nil {
			log.Fatal(err)
		}
	}

	mustPrepare := func(stmt string) *sql.Stmt {
		r, err := db.Prepare(stmt)
		if err != nil {
			log.Fatal(err)
		}
		return r
	}

	genGetall := func(sort string) *sql.Stmt {
		c := `SELECT id, title, called, level, created
                        FROM tasks ORDER BY %s LIMIT %d OFFSET $1`
		stmt, err := db.Prepare(fmt.Sprintf(c, sort, entries))
		if err != nil {
			log.Fatal(err)
		}
		return stmt
	}

	genFind := func(sort string) *sql.Stmt {
		c := `SELECT id, title, called, level, created
                        FROM tasks WHERE title ~* $1
                    ORDER BY %s LIMIT %d OFFSET $2`
		stmt, err := db.Prepare(fmt.Sprintf(c, sort, entries))
		if err != nil {
			log.Fatal(err)
		}
		return stmt
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	mustExec(`CREATE TABLE IF NOT EXISTS tasks (
                         id       SERIAL PRIMARY KEY UNIQUE,
                         title    VARCHAR(100) NOT NULL,
                         author   VARCHAR(100),
                         discrip  TEXT,
                         called   INT DEFAULT 0,
                         level    INT NOT NULL,
                         created  TIMESTAMP DEFAULT NOW());`)

	mustExec(`CREATE TABLE IF NOT EXISTS solutions (
		         suggested   INT DEFAULT 1,
                         solves      INT NOT NULL REFERENCES tasks (id) ON DELETE CASCADE,
                         solution    VARCHAR(256) NOT NULL,
                         first       TIMESTAMP DEFAULT NOW(),
                         last        TIMESTAMP DEFAULT NOW(),
                         PRIMARY KEY (solution, solves))`)

	mustExec(`CREATE TABLE IF NOT EXISTS words (
                         word    VARCHAR(75) NOT NULL,
                         matches BOOL NOT NULL,
                         task    INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE)`)

	getad = genGetall("level")
	getan = genGetall("created DESC")
	getap = genGetall("(SELECT COUNT(1) FROM solutions WHERE solves = id) DESC")

	findd = genFind("level")
	findn = genFind("created DESC")
	findp = genFind("(SELECT COUNT(1) FROM solutions WHERE solves = id) DESC")

	finddr = genFind("level DESC")
	findnr = genFind("created")
	findpr = genFind("(SELECT COUNT(1) FROM solutions WHERE solves = id)")

	create = mustPrepare(`INSERT INTO tasks (title, author, discrip, level)
                                   VALUES ($1, $2, $3, $4) RETURNING id`)

	create_word = mustPrepare(`INSERT INTO words (word, matches, task)
                                        VALUES ($1, $2, $3)`)

	subm = mustPrepare(`INSERT INTO solutions (solution, solves)
                                 VALUES ($1, $2)
                            ON CONFLICT (solution, solves)
                          DO UPDATE SET suggested = solutions.suggested + 1,
                                             last = NOW()`)

	getsolc = mustPrepare(`SELECT COUNT(1)
	                       FROM solutions
                               WHERE solves = $1`)

	getsol = mustPrepare(`SELECT suggested, solution, first, last
                                FROM solutions
                               WHERE solves = $1
                            ORDER BY -suggested`)

	getwrd = mustPrepare(`SELECT word, matches
                                FROM words
                               WHERE task = $1`)

	getone = mustPrepare(`SELECT title, author, discrip, called, level, created
                                FROM tasks
                               WHERE id = $1`)

	getrnd = mustPrepare(`SELECT id
                                FROM tasks
                               WHERE level = 0 OR level >= $1
                            ORDER BY called / EXTRACT(EPOCH FROM created) DESC
                              OFFSET RANDOM() * (SELECT COUNT(1) - 1 FROM tasks)
                               LIMIT 1`)

	inccou = mustPrepare(`UPDATE tasks
                                 SET called = called + 1
                               WHERE id = $1`)
}
