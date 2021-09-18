// Copyright (c) 2021 twihike. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package authui

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func KratosProxy(prefix string) http.Handler {
	k, err := url.Parse(KratosPublicURL)
	if err != nil {
		panic(err)
	}
	director := func(request *http.Request) {
		request.URL.Scheme = k.Scheme
		request.URL.Host = k.Host
		request.URL.Path = strings.TrimPrefix(request.URL.Path, prefix)
	}
	return &httputil.ReverseProxy{Director: director}
}
