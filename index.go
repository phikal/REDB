// Licensed under GPL, 2016
// Refer to LICENSE for more details
// Refer to README for structural information

package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

func index(order string, start bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			before, next string
			tasks        = make([]Task, 0, entries)
			page         = 1
			values       = r.URL.Query()
			query        *sql.Stmt
		)

		page, err := strconv.Atoi(values.Get("p"))
		if err != nil || page < 1 {
			page = 1
		}

		switch order {
		case diff:
			query = getad
		case pop:
			query = getap
		case new:
			fallthrough
		default:
			query = getan
		}

		rows, err := query.Query((page - 1) * entries)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			t.ExecuteTemplate(w, "error.gtml", err.Error())
			return
		}
		defer rows.Close()

		tx, err := db.Begin()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			t.ExecuteTemplate(w, "error.gtml", err.Error())
			return
		}

		for rows.Next() {
			var T Task
			if err = rows.Scan(
				&T.Id,
				&T.Title,
				&T.Called,
				&T.Created,
			); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				t.ExecuteTemplate(w, "error.gtml", err.Error())
				return
			}
			if err = T.loadWords(tx); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				t.ExecuteTemplate(w, "error.gtml", err.Error())
				return
			}
			tasks = append(tasks, T)
		}
		tx.Commit()

		// define "before" URL, if possible
		if page > 1 {
			values.Set("p", strconv.Itoa(page-1))
			u := r.URL
			u.RawQuery = values.Encode()
			before = u.String()
		}

		// define "next" URL, if necessary
		if len(tasks) >= entries {
			values.Set("p", strconv.Itoa(page+1))
			u := r.URL
			u.RawQuery = values.Encode()
			next = u.String()
		}

		err = t.ExecuteTemplate(w, "index.gtml", struct {
			T []Task // list of tasks
			P int    // current page and item number
			S bool   // show start dialog
			B string // brevious page
			N string // next page
			O string // ordering
		}{
			tasks,
			page,
			start,
			before,
			next,
			order,
		})

		if err != nil {
			fmt.Println(err)
		}
	}
}
