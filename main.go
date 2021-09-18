// Copyright (c) 2021 twihike. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"

	"github.com/twihike/ory-selfservice-ui/selfservice/authui"
	"github.com/twihike/ory-selfservice-ui/selfservice/oauthui"
)

func main() {
	http.HandleFunc("/auth/registration", authui.Registraction)
	http.HandleFunc("/auth/login", authui.Login)
	http.HandleFunc("/auth/dashboard", authui.Dashboard)
	http.HandleFunc("/auth/settings", authui.Settings)
	http.HandleFunc("/auth/verify", authui.Verification)
	http.HandleFunc("/auth/recovery", authui.Recovery)
	http.HandleFunc("/auth/error", authui.ErrorPage)
	http.Handle("/auth/api/", authui.KratosProxy("/auth/api"))
	http.HandleFunc("/oauth/login", oauthui.OAuthLogin)
	http.HandleFunc("/oauth/consent", oauthui.Consent)
	http.HandleFunc("/oauth/error", oauthui.ErrorPage)
	log.Println("started")
	log.Fatal(http.ListenAndServe(":4455", nil))
}
