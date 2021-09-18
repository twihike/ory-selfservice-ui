// Copyright (c) 2021 twihike. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package oauthui

import (
	"context"
	"log"
	"net/http"
	"text/template"

	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"github.com/twihike/ory-selfservice-ui/selfservice"
	"github.com/twihike/ory-selfservice-ui/selfservice/authui"
)

var (
	tmplConsent = template.Must(template.ParseFS(selfservice.Views, "views/main.html", "views/partials/*.html", "views/consent.html"))
)

func Consent(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		consentGet(w, r)
		return
	}
	if r.Method == http.MethodPost {
		consentPost(w, r)
		return
	}
	http.NotFound(w, r)
}

func consentGet(w http.ResponseWriter, r *http.Request) {
	consentChallenge := r.URL.Query().Get("consent_challenge")
	if consentChallenge == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	gcrParams := &admin.GetConsentRequestParams{
		ConsentChallenge: consentChallenge,
		Context:          context.Background(),
	}
	gcrOK, err := HydraAdminClient.Admin.GetConsentRequest(gcrParams)
	if err != nil {
		log.Println("gcr:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if gcrOK.Payload.Skip {
		sess, err := getSession(w, r, gcrOK.Payload.RequestedScope)
		if err != nil {
			return
		}
		acrParams := &admin.AcceptConsentRequestParams{
			Body: &models.AcceptConsentRequest{
				GrantScope:               gcrOK.Payload.RequestedScope,
				GrantAccessTokenAudience: gcrOK.Payload.RequestedAccessTokenAudience,
				Session: &models.ConsentRequestSession{
					AccessToken: sess,
					IDToken:     sess,
				},
			},
			ConsentChallenge: consentChallenge,
			Context:          context.Background(),
		}
		acrOK, err := HydraAdminClient.Admin.AcceptConsentRequest(acrParams)
		if err != nil {
			log.Println("acr-skip:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, *acrOK.Payload.RedirectTo, http.StatusFound)
		return
	}

	if err := tmplConsent.Execute(w, map[string]interface{}{"Flow": gcrOK}); err != nil {
		log.Println("consent-tmpl:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func consentPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("consentPost ParseForm:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	consentChallenge := r.FormValue("challenge")
	if consentChallenge == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	if r.FormValue("submit") != "Allow access" {
		rcrParams := &admin.RejectConsentRequestParams{
			Body: &models.RejectRequest{
				Error:            "access_denied",
				ErrorDescription: "The resource owner denied the request",
			},
			ConsentChallenge: consentChallenge,
			Context:          ctx,
		}
		rcrOK, err := HydraAdminClient.Admin.RejectConsentRequest(rcrParams)
		if err != nil {
			log.Println("rcr:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, *rcrOK.Payload.RedirectTo, http.StatusFound)
		return
	}

	gcrParams := &admin.GetConsentRequestParams{
		ConsentChallenge: consentChallenge,
		Context:          ctx,
	}
	gcrOK, err := HydraAdminClient.Admin.GetConsentRequest(gcrParams)
	if err != nil {
		log.Println("gcr:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sess, err := getSession(w, r, r.Form["grant_scope"])
	if err != nil {
		return
	}
	acrParams := &admin.AcceptConsentRequestParams{
		Body: &models.AcceptConsentRequest{
			GrantScope:               r.Form["grant_scope"],
			GrantAccessTokenAudience: gcrOK.Payload.RequestedAccessTokenAudience,
			Remember:                 r.FormValue("remember") == "1",
			Session: &models.ConsentRequestSession{
				AccessToken: sess,
				IDToken:     sess,
			},
		},
		ConsentChallenge: consentChallenge,
		Context:          ctx,
	}
	acrOK, err := HydraAdminClient.Admin.AcceptConsentRequest(acrParams)
	if err != nil {
		log.Println("acr:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, *acrOK.Payload.RedirectTo, http.StatusFound)
}

func getSession(w http.ResponseWriter, r *http.Request, scope []string) (map[string]interface{}, error) {
	sess, _, err := authui.ToSession(w, r)
	if err != nil {
		log.Println("whoami:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}
	result := map[string]interface{}{}
	for _, s := range scope {
		if s == "email" {
			result["email"] = sess.Identity.VerifiableAddresses[0].Value
			result["email_verified"] = sess.Identity.VerifiableAddresses[0].Verified
		}
	}
	return result, nil
}
