// Copyright (c) 2021 twihike. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package oauthui

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"os"

	httptransport "github.com/go-openapi/runtime/client"
	hydra "github.com/ory/hydra-client-go/client"
)

var (
	HydraAdminURL    = os.Getenv("HYDRA_ADMIN_URL")
	HydraAdminClient = NewHydraClient(HydraAdminURL)
)

func NewHydraClient(endpoint string) *hydra.OryHydra {
	hydraURL, err := url.Parse(endpoint)
	if err != nil {
		panic("invalid hydra admin url")
	}
	// config := hydra.NewHTTPClientWithConfig(nil, &hydra.TransportConfig{
	// 	Schemes:  []string{hydraURL.Scheme},
	// 	Host:     hydraURL.Host,
	// 	BasePath: hydraURL.Path,
	// })
	skipTlsClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 0,
	}
	transport := httptransport.NewWithClient(
		hydraURL.Host,
		hydraURL.Path,
		[]string{hydraURL.Scheme},
		skipTlsClient,
	)
	config := hydra.New(transport, nil)
	return config
}
