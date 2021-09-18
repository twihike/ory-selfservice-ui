// Copyright (c) 2021 twihike. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package authui

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/twihike/ory-selfservice-ui/selfservice"
)

var (
	tmplDashboard = template.Must(template.ParseFS(selfservice.Views, "views/main.html", "views/partials/*.html", "views/dashboard.html"))
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	var b bytes.Buffer
	fmt.Fprintln(&b, "----")
	for k, v := range r.Header {
		fmt.Fprintln(&b, k, ":", v)
	}
	fmt.Fprintln(&b, "----")
	for _, v := range r.Cookies() {
		fmt.Fprintln(&b, v)
	}
	fmt.Fprintln(&b, "----")
	if _, s, err := ToSession(w, r); err == nil {
		fmt.Fprint(&b, s)
		fmt.Fprintln(&b, "----")
	} else {
		log.Println(err)
	}

	data := map[string]interface{}{"Flow": b.String()}
	if cookie := r.Header.Get("Cookie"); cookie != "" {
		logoutURL, _, err := KratosPublicClient.V0alpha1Api.CreateSelfServiceLogoutFlowUrlForBrowsers(context.Background()).Cookie(cookie).Execute()
		if err == nil {
			data["LogoutURL"] = logoutURL.LogoutUrl
		}
	}

	if err := tmplDashboard.Execute(w, data); err != nil {
		log.Println("dashboard-tmpl:", err)
		http.Redirect(w, r, "/auth/login", http.StatusFound)
		return
	}
}
