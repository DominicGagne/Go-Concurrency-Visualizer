//go:generate go-bindata-assetfs page/...
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

type PageInfo struct {
	Title string
}

// indexTmpl is a html template for index page.
var (
	devMode   bool
	indexTmpl *template.Template
)

func init() {
	devMode = os.Getenv("GOTRACE_DEVMODE") == "1"
	if devMode {
		indexTmpl = template.Must(template.New("index.html").Parse("page/index.html"))
	} else {
		data, err := Asset("page/index.html")
		if err != nil {
			panic(err)
		}
		indexTmpl = template.Must(template.New("index.html").Parse(string(data)))
	}
}

// StartServer generates webpage, serves it via http
// and tries to open it using default browser.
func StartServer(bind string, data []byte, params *Params) error {
	// Serve data as data.js
	http.HandleFunc("/data.js", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("var data = "))
		w.Write(data)
	}))

	// Serve params as params.js
	http.HandleFunc("/params.js", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("var params = "))
		data, err := json.MarshalIndent(params, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, "[ERROR] failed to render params:", err)
			return
		}
		w.Write(data)
	}))

	// Handle static files
	if devMode {
		var fs http.FileSystem
		fs = http.Dir("page")
		http.Handle("/", http.FileServer(fs))
	} else {
		http.Handle("/", http.FileServer(assetFS()))
	}
	go StartBrowser("http://localhost" + bind)

	return http.ListenAndServe(bind, nil)
}

// handler handles index page.
func handler(w http.ResponseWriter, r *http.Request, info *PageInfo) {
	err := indexTmpl.Execute(w, info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, "[ERROR] failed to render template:", err)
		return
	}
}

// StartBrowser tries to open the URL in a browser
// and reports whether it succeeds.
//
// Orig. code: golang.org/x/tools/cmd/cover/html.go
func StartBrowser(url string) bool {
	// try to start the browser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	fmt.Println("If browser window didn't appear, please go to this url:", url)
	return cmd.Start() == nil
}
