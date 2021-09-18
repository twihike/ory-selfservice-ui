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
	tmplRecovery = template.Must(template.ParseFS(selfservice.Views, "views/main.html", "views/partials/*.html", "views/recovery.html"))
)

func Recovery(w http.ResponseWriter, r *http.Request) {
	flowID := r.URL.Query().Get("flow")
	if flowID == "" {
		http.Redirect(w, r, KratosBrowserURL+"/self-service/recovery/browser", http.StatusFound)
		return
	}
	cookie := r.Header.Get("Cookie")
	if cookie == "" {
		http.Redirect(w, r, KratosBrowserURL+"/self-service/recovery/browser", http.StatusFound)
		return
	}

	flow, _, err := KratosPublicClient.V0alpha1Api.GetSelfServiceRecoveryFlow(context.Background()).Id(flowID).Cookie(cookie).Execute()
	if err != nil {
		log.Println("recovery:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Flow": flow}
	logoutURL, _, err := KratosPublicClient.V0alpha1Api.CreateSelfServiceLogoutFlowUrlForBrowsers(context.Background()).Cookie(cookie).Execute()
	if err == nil {
		data["LogoutURL"] = logoutURL.LogoutUrl
	}
	if err := tmplRecovery.Execute(w, data); err != nil {
		log.Println("recovery-tmpl:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
