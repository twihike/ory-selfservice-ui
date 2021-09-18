// Copyright (c) 2021 twihike. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package oauthui

import (
	"log"
	"net/http"
	"text/template"

	"github.com/twihike/ory-selfservice-ui/selfservice"
)

var (
	tmplError = template.Must(template.ParseFS(selfservice.Views, "views/main.html", "views/partials/*.html", "views/error_oauth.html"))
)

func ErrorPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]map[string]interface{}{"Flow": {}}
	if e := r.URL.Query().Get("error"); e != "" {
		data["Flow"]["Error"] = e
	}
	if e := r.URL.Query().Get("error_description"); e != "" {
		data["Flow"]["ErrorDescription"] = e
	}

	if err := tmplError.Execute(w, data); err != nil {
		log.Println("error-page-tmpl:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
