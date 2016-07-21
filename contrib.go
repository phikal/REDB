// Licensed under GPL, 2016
// Refer to LICENSE for more details
// Refer to README for structural information

package main

import (
	"fmt"
	"github.com/dchest/captcha"
	"log"
	"net/http"
)

func contrib(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var sug Task
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			t.ExecuteTemplate(w, "error.gtml", err)
			return
		}

		cid := r.FormValue("cid")
		csol := r.FormValue("csol")
		if !captcha.VerifyString(cid, csol) {
			w.WriteHeader(http.StatusNotAcceptable)
			t.ExecuteTemplate(w, "error.gtml", errorPage{
				"406 - Not Acceptable",
				"The Captcha solution was invalid",
			})
			return
		}

		sug.Title = r.FormValue("title")
		sug.Author = r.FormValue("author")
		sug.Discription = r.FormValue("discr")
		sug.Match = r.Form["match"]
		sug.Dmatch = r.Form["dmatch"]

		if len(sug.Title) == 0 || len(sug.Title) > 100 {
			w.WriteHeader(http.StatusNotAcceptable)
			t.ExecuteTemplate(w, "error.gtml", errorPage{
				"406 - Not Acceptable",
				"The proposed Title was either too long or too short (min 1 char, max 100 chars)",
			})
		} else if len(sug.Author) > 100 {
			w.WriteHeader(http.StatusNotAcceptable)
			t.ExecuteTemplate(w, "error.gtml", errorPage{
				"406 - Not Acceptable",
				"The proposed Author name was too long (max 100 chars)",
			})
		} else if len(sug.Discription) > 512 {
			w.WriteHeader(http.StatusNotAcceptable)
			t.ExecuteTemplate(w, "error.gtml", errorPage{
				"406 - Not Acceptable",
				"The proposed Description was either too long (max 512 chars)",
			})
		} else if err := sug.isAcceptable(); err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			t.ExecuteTemplate(w, "error.gtml", errorPage{
				"406 - Not Acceptable",
				err.Error(),
			})
		} else {
			var ment, dent uint // amount of entries to be match and not

			for _, e := range sug.Match {
				if len(e) > 0 {
					ment++
				}
			}
			for _, e := range sug.Dmatch {
				if len(e) > 0 {
					dent++
				}
			}

			if ment < 2 || ment > 12 {
				w.WriteHeader(http.StatusPartialContent)
				t.ExecuteTemplate(w, "error.gtml", errorPage{
					"406 - Not Acceptable",
					"Too many or too few entries to match (min 2, max 12)",
				})
			} else if dent < 1 || dent > 12 {
				w.WriteHeader(http.StatusPartialContent)
				t.ExecuteTemplate(w, "error.gtml", errorPage{
					"406 - Not Acceptable",
					"Too many or too few entries not to match (min 1, max 12)",
				})
			} else if c, err := sug.insert(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				t.ExecuteTemplate(w, "error.gtml", err)
			} else {
				w.Header().Add("Refresh", fmt.Sprintf("0; url=/r/%x", c))
			}
		}
	default:
		if err := t.ExecuteTemplate(w, "contrib.gtml", captcha.New()); err != nil {
			log.Panicln(err)
		}
	}
}
