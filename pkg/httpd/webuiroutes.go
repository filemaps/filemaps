// Copyright (c) 2017, CodeBoy. All rights reserved.
//
// This Source Code Form is subject to the terms of the
// license that can be found in the LICENSE file.

package httpd

import (
	"bytes"
	"compress/gzip"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	assets map[string][]byte
	eTags  map[string]string
)

func routeWebUI(r *httprouter.Router, webUIPath string) {
	routePath := UIURL + "/*filepath"
	if webUIPath != "" {
		// web UI path provided, serve files from there
		r.ServeFiles(routePath, http.Dir(webUIPath))
	} else {
		// use bundled web UI files (built webui.go)
		r.GET(routePath, webUI)
	}
}

func webUI(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// remove preceding slash
	path := ps.ByName("filepath")[1:]
	if path == "" {
		path = "index.html"
	}

	// use ETag for efficient caching
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/ETag
	if r.Header.Get("If-None-Match") == eTags[path] {
		w.WriteHeader(304)
		return
	}

	bs, ok := assets[path]
	if !ok {
		w.WriteHeader(404)
		fmt.Fprint(w, "Not found")
		log.WithFields(log.Fields{
			"path": path,
		}).Error("Web UI file not found")
		return
	}

	w.Header().Set("Content-Type", getContentType(path))
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
	} else {
		// client browser does not support gzip encoding,
		// so need to unzip content
		var gr *gzip.Reader
		gr, _ = gzip.NewReader(bytes.NewReader(bs))
		bs, _ = ioutil.ReadAll(gr)
		gr.Close()
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(bs)))
	w.Header().Set("ETag", eTags[path])

	w.Write(bs)
}

func getContentType(file string) string {
	ext := filepath.Ext(file)
	switch ext {
	case ".css":
		return "text/css"
	case ".html":
		return "text/html"
	case ".js":
		return "application/javascript"
	case ".png":
		return "image/png"
	case ".ttf":
		return "application/x-font-ttf"
	default:
		return mime.TypeByExtension(ext)
	}
}

// setAssets sets static assets.
// Called by auto-generated webui.go
func setAssets(a map[string][]byte) {
	assets = a
}

// setETags sets entity tags for static assets.
// Called by auto-generated webui.go
func setETags(a map[string]string) {
	eTags = a
}
