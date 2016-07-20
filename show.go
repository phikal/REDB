// Licensed under GPL, 2016
// Refer to LICENSE for more details
// Refer to README for structural information

package main

import (
	"fmt"
	"log"
	"net/http"
)

func showRegex(w http.ResponseWriter, r *http.Request) {
	tx, err := db.Begin()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		t.ExecuteTemplate(w, "error.gtml", errorPage{
			"500 - Internal Server Error",
			err.Error(),
		})
		return
	}

	var T Task
	_, err = fmt.Sscanf(r.RequestURI, "/r/%x", &T.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		t.ExecuteTemplate(w, "error.gtml", errorPage{
			"400 - Bad Request",
			"Could not resolve " + r.RequestURI,
		})
		return
	}

	if err := tx.Stmt(getone).QueryRow(T.Id).Scan(
		&T.Title,
		&T.Author,
		&T.Discription,
		&T.Called,
		&T.Level,
		&T.Created,
	); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		t.ExecuteTemplate(w, "error.gtml", errorPage{
			"404 - Not Found",
			fmt.Sprintf("Could not find task with the ID %x", T.Id),
		})
	} else {
		T.loadWords(tx)
		if err := T.loadSolutions(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			t.ExecuteTemplate(w, "error.gtml", errorPage{
				"500 - Internal Server Error",
				err.Error(),
			})
			return
		} else if err := tx.Commit(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			t.ExecuteTemplate(w, "error.gtml", errorPage{
				"500 - Internal Server Error",
				err.Error(),
			})
			return
		} else if err := t.ExecuteTemplate(w, "show.gtml", T); err != nil {
			log.Println(err)
		}
	}
}
