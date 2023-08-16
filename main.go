//go:generate go run gen.go
package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

//go:embed www/*.html
//go:embed www/*.css
//go:embed www/*.js
var www embed.FS

type server struct {
	indexHtml *template.Template
}

func (s *server) getRoot(w http.ResponseWriter, r *http.Request) {
	s.indexHtml.ExecuteTemplate(w, "index.html", struct{}{})
}

func (s *server) getW3(w http.ResponseWriter, r *http.Request) {
	file, err := www.ReadFile(strings.TrimPrefix(r.RequestURI, "/"))
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "NOT FOUND: %s", r.RequestURI)
		return
	}

	contentType := "text/javascript"
	if strings.HasSuffix(r.RequestURI, ".css") {
		contentType = "text/css"
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(file)
}

func main() {
	f, err := www.ReadFile("www/index.html")
	if err != nil {
		log.Fatal(err)
	}
	t, err := template.New("index.html").Parse(string(f))
	if err != nil {
		log.Fatal(err)
	}
	s := &server{t}
	http.HandleFunc("/www/", s.getW3)
	http.HandleFunc("/", s.getRoot)
	if err := http.ListenAndServe(":3333", nil); err != nil {
		log.Fatal(err)
	}
}
