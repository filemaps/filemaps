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
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

const (
	target       = "github.com/filemaps/filemaps/cmd/filemaps"
	webuiPath    = "filemaps-webui/"
	webuiVersion = "0.1.5"
)

var tmpl = template.Must(template.New("assets").Parse(`package httpd

import "encoding/base64"

func init() {
	var a = make(map[string][]byte, {{.Assets | len}})
{{range $asset := .Assets}}
	a["{{$asset.Name}}"], _ = base64.StdEncoding.DecodeString("{{$asset.Content}}"){{end}}
	setAssets(a)
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
	assets  []asset
	goos    string
	goarch  string
	pkgdir  string
	version string
)

func init() {
	flag.StringVar(&goarch, "goarch", runtime.GOARCH, "Target architecture")
	flag.StringVar(&goos, "goos", runtime.GOOS, "Target OS")
	flag.StringVar(&pkgdir, "pkgdir", "", "Pkgdir")
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		// no command provided
		runCommand("install")
	} else {
		runCommand(flag.Arg(0))
	}
}

func runCommand(cmd string) {
	switch cmd {
	case "clean":
		clean()
	case "install":
		install()
	case "maketar":
		buildPkg("tar")
	case "makezip":
		buildPkg("zip")
	case "setup":
		setup()
	case "test":
		test()
	}
}

func clean() {
	os.RemoveAll("build")
}

func install() {
	log.Info("Building and installing File Maps")
	version = readVersion()

	bundleWebUI(webuiPath, "pkg/httpd/webui.go")

	os.Setenv("GOOS", goos)
	os.Setenv("GOARCH", goarch)

	args := []string{"install", "-ldflags", ldflags()}
	if pkgdir != "" {
		args = append(args, "-pkgdir")
		args = append(args, pkgdir)
	}
	args = append(args, target)
	exe("go", args...)
}

func setup() {
	downloadWebUI()
}

func downloadWebUI() {
	url := "https://github.com/filemaps/filemaps-webui/releases/download/v" + webuiVersion + "/filemaps-webui-build.tar.gz"
	exe("curl", "-L", "-O", url)
	exe("tar", "xf", "filemaps-webui-build.tar.gz")
}

func buildPkg(format string) {
	install()
	version = readVersion()
	name := fmt.Sprintf("filemaps-%s-%s-%s", goos, goarch, version)
	targetPath := "build/" + name
	os.MkdirAll(targetPath, 0755)

	// copy license and readme files
	exe("cp", "LICENSE", "README.md", targetPath)

	// copy executable
	binPath := os.Getenv("GOPATH") + "/bin/"
	if goos != runtime.GOOS || goarch != runtime.GOARCH {
		binPath = binPath + goos + "_" + goarch + "/"
	}
	binPath = binPath + "filemaps"
	if goos == "windows" {
		binPath = binPath + ".exe"
	}
	exe("cp", binPath, targetPath)

	// create archive
	os.Chdir("build")
	a := name
	if format == "zip" {
		a = a + ".zip"
		exe("zip", "-r", a, name)
	} else if format == "tar" {
		a = a + ".tar.xz"
		exe("tar", "cJf", a, name)
	}
	os.Chdir("..")
	log.WithFields(log.Fields{
		"file": a,
	}).Info("Archive created")
}

func test() {
	exe("go", "test", "-v", "./...")
}

// ldflags sets variables for building process, such as version number
func ldflags() string {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "-X main.Version=%s", version)
	return b.String()
}

func readVersion() string {
	o, err := exeRead("git", "describe", "--always", "--dirty")
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Could not read version using git")
		return ""
	}
	v := string(o)
	log.WithFields(log.Fields{
		"version": v,
	}).Info("Version detected from git")
	return v
}

// exeRead executes given command and returns output
func exeRead(cmd string, args ...string) ([]byte, error) {
	cmdh := exec.Command(cmd, args...)
	bs, err := cmdh.CombinedOutput()
	return bytes.TrimSpace(bs), err
}

// exe executes given command
func exe(cmd string, args ...string) {
	cmdh := exec.Command(cmd, args...)
	cmdh.Stdout = os.Stdout
	cmdh.Stderr = os.Stderr
	err := cmdh.Run()
	if err != nil {
		log.WithFields(log.Fields{
			"cmd":  cmd,
			"args": args,
		}).Fatal(err)
	}
}

// getWalkFunc is walker function for bundling web ui files.
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
