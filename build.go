// Copyright (C) 2017 File Maps Backend Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	log "github.com/Sirupsen/logrus"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	target    = "github.com/filemaps/filemaps-backend/cmd/filemaps"
	webuiPath = "filemaps-webui/src/"
)

var tmpl = template.Must(template.New("assets").Parse(`package httpd

import "encoding/base64"

func GetAssets() map[string][]byte {
	var assets = make(map[string][]byte, {{.Assets | len}})
{{range $asset := .Assets}}
	assets["{{$asset.Name}}"], _ = base64.StdEncoding.DecodeString("{{$asset.Content}}"){{end}}
	return assets
}
`))

type asset struct {
	Name    string
	Content string
}

type tmplVars struct {
	Assets []asset
}

var (
	assets []asset
)

func main() {
	log.Info("Building and installing File Maps")
	bundleWebUI(webuiPath, "pkg/httpd/webui.go")
	run("go", "install", target)
}

func run(cmd string, args ...string) {
	cmdh := exec.Command(cmd, args...)
	cmdh.Stdout = os.Stdout
	cmdh.Stderr = os.Stderr
	err := cmdh.Run()
	if err != nil {
		log.WithFields(log.Fields{
			"cmd":  cmd,
			"args": args,
		}).Error(err)
	}
}

func getWalkFunc(base string) filepath.WalkFunc {
	return func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasPrefix(filepath.Base(name), ".") {
			// ignore files beginning with dot
			return nil
		}

		if info.Mode().IsRegular() {
			f, err := os.Open(name)
			if err != nil {
				return err
			}

			// read file contents and gzip it to buffer
			var buf bytes.Buffer
			g := gzip.NewWriter(&buf)
			io.Copy(g, f)
			f.Close()
			g.Flush()
			g.Close()

			// create asset struct and append it to vars
			name, _ = filepath.Rel(base, name)
			assets = append(assets, asset{
				Name:    filepath.ToSlash(name),
				Content: base64.StdEncoding.EncodeToString(buf.Bytes()),
			})
		}
		return nil
	}
}

// packageWebUI packages Web UI files into a single go file.
// All files are gzipped and base64 encoded into static strings
// which are served by HTTP server.
func bundleWebUI(path string, out string) {
	log.Info("Bundling Web UI")
	filepath.Walk(path, getWalkFunc(path))

	if len(assets) == 0 {
		log.WithFields(log.Fields{
			"path": path,
		}).Fatal("No Web UI files found. Make sure you have installed them into path")
		return
	}
	var buf bytes.Buffer
	tmpl.Execute(&buf, tmplVars{
		Assets: assets,
	})

	bs, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(out, bs, 0644)
	if err != nil {
		panic(err)
	}
	log.WithFields(log.Fields{
		"output": out,
	}).Info("Web UI bundled")
}
