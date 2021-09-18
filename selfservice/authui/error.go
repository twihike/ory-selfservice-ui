// Copyright (c) 2021 twihike. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package authui

import (
	"context"
	"log"
	"net/http"
	"text/template"

	"github.com/twihike/ory-selfservice-ui/selfservice"
)

var (
	tmplError = template.Must(template.ParseFS(selfservice.Views, "views/main.html", "views/partials/*.html", "views/error.html"))
)

func ErrorPage(w http.ResponseWriter, r *http.Request) {
	e := r.URL.Query().Get("error")
	if e == "" {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("An error occurred"))
		return
	}

	flow, _, err := KratosPublicClient.V0alpha1Api.GetSelfServiceError(context.Background()).Id(e).Execute()
	if err != nil {
		log.Println("error-page:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Flow": flow}
	if err := tmplError.Execute(w, data); err != nil {
		log.Println("error-page-tmpl:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
