// Copyright (c) 2021 twihike. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package oauthui

import (
	"context"
	"log"
	"net/http"
	"net/url"

	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"github.com/twihike/ory-selfservice-ui/selfservice/authui"
)

func OAuthLogin(w http.ResponseWriter, r *http.Request) {
	loginChallenge := r.URL.Query().Get("login_challenge")
	if loginChallenge == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	glrParams := &admin.GetLoginRequestParams{
		LoginChallenge: loginChallenge,
		Context:        context.Background(),
	}
	glrOK, err := HydraAdminClient.Admin.GetLoginRequest(glrParams)
	if err != nil {
		log.Println("glr:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if glrOK.Payload.Skip != nil && *glrOK.Payload.Skip {
		alrParams := &admin.AcceptLoginRequestParams{
			Body: &models.AcceptLoginRequest{
				Subject: glrOK.Payload.Subject,
			},
			LoginChallenge: loginChallenge,
			Context:        context.Background(),
		}
		alrOK, err := HydraAdminClient.Admin.AcceptLoginRequest(alrParams)
		if err != nil {
			log.Println("alr-skip:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, *alrOK.Payload.RedirectTo, http.StatusFound)
		return
	}

	if c, err := r.Cookie(authui.KratosCookieName); err != nil || c.Value == "" {
		log.Println(authui.KratosCookieName)
		redirectLogin(w, r)
		return
	}

	sess, _, err := authui.ToSession(w, r)
	if err != nil {
		log.Println("whoami:", err)
		redirectLogin(w, r)
		return
	}

	alrParams := &admin.AcceptLoginRequestParams{
		Body: &models.AcceptLoginRequest{
			Subject: &sess.Identity.Id,
		},
		LoginChallenge: loginChallenge,
		Context:        context.Background(),
	}
	alrOK, err := HydraAdminClient.Admin.AcceptLoginRequest(alrParams)
	if err != nil {
		log.Println("alr:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, *alrOK.Payload.RedirectTo, http.StatusFound)
}

func redirectLogin(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(authui.KratosBrowserURL + "/self-service/login/browser")
	q := u.Query()
	q.Add("refresh", "true")
	q.Add("return_to", r.URL.String())
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}
