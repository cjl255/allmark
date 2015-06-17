// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"allmark.io/modules/common/route"
	"allmark.io/modules/web/view/viewmodel"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"text/template"
)

func getRouteFromRequest(r *http.Request) route.Route {
	return route.NewFromRequest(r.URL.Path)
}

func getBaseUrlFromRequest(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	return scheme + "://" + r.Host
}

func renderTemplate(template *template.Template, model interface{}, writer io.Writer) error {
	return template.Execute(writer, model)
}

func renderViewModelAsJSON(viewModel viewmodel.Model, writer io.Writer) error {
	bytes, err := json.MarshalIndent(viewModel, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}

func getPageParameterFromUrl(url url.URL) (page int, parameterIsAvailable bool) {
	pageParam := url.Query().Get("page")
	if pageParam == "" {
		return 0, false
	}

	page64, err := strconv.ParseInt(pageParam, 10, 64)
	if err != nil {
		return 0, true
	}

	if page64 < 1 {
		return 0, true
	}

	return int(page64), true
}

func getQueryParameterFromUrl(url url.URL) (query string, parameterIsAvailable bool) {
	queryParam := url.Query().Get("q")
	if queryParam == "" {
		return "", false
	}

	return queryParam, true
}