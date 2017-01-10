// Licensed under GPL, 2016
// Refer to LICENSE for more details
// Refer to README for structural information

package main

import (
	"database/sql"
	"net/http"
	"strconv"
)

func search(w http.ResponseWriter, r *http.Request) {
	var (
		before, next, sort string
		tasks              []Task
		page               = 1
		values             = r.URL.Query()
		query              *sql.Stmt
	)

	page, err := strconv.Atoi(values.Get("p"))
	if err != nil || page < 1 {
		page = 1
	}

	// prepare query, if `q` specified
	if values.Get("q") != "" {
		sort = values.Get("s")

		// reverse
		if values.Get("r") != "" {
			switch sort {
			case pop:
				query = findp
			case diff:
				query = findd
			default:
				sort = new
				fallthrough
			case new:
				query = findn
			}
		} else {
			switch sort {
			case pop:
				query = findpr
			case diff:
				query = finddr
			default:
				sort = new
				fallthrough
			case new:
				query = findnr

			}
		}

		// query database
		data, err := query.Query(values.Get("q"), (page-1)*25)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			t.ExecuteTemplate(w, "error.gtml", errorPage{
				"505 - Internal Server Error",
				err.Error(),
			})
			return
		}
		defer data.Close()

		tx, err := db.Begin()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			t.ExecuteTemplate(w, "error.gtml", errorPage{
				"505 - Internal Server Error",
				err.Error(),
			})
			return
		}

		// read into slice
		tasks = make([]Task, 0, entries)
		for data.Next() {
			var T Task
			if err = data.Scan(
				&T.Id,
				&T.Title,
				&T.Called,
				&T.Created,
			); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				t.ExecuteTemplate(w, "error.gtml", errorPage{
					"505 - Internal Server Error",
					err.Error(),
				})
				return
			}
			T.loadWords(tx)
			tasks = append(tasks, T)
		}

		tx.Commit()
	}

	// generate URL for the previous page, if possible
	if page > 1 {
		values.Set("p", strconv.Itoa(page-1))
		u := r.URL
		u.RawQuery = values.Encode()
		before = u.String()
	}

	// generate URL for the next page, if necessary
	if len(tasks) >= entries {
		values.Set("p", strconv.Itoa(page+1))
		u := r.URL
		u.RawQuery = values.Encode()
		next = u.String()
	}

	t.ExecuteTemplate(w, "search.gtml", struct {
		T          []Task
		P          int
		Q, B, N, S string
		R          bool
	}{
		tasks,
		page,
		values.Get("q"),
		before,
		next,
		sort,
		values.Get("r") == "y",
	})
}
