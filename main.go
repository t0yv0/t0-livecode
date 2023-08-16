//go:generate go run gen.go
package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

//go:embed www/*.html
//go:embed www/*.css
//go:embed www/*.js
var www embed.FS

type codeRepo interface {
	Store(id, code string) error
	Load(id string) (string, error)
}

type server struct {
	indexHtml *template.Template
	appHtml   *template.Template
	codeRepo  codeRepo
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

func (s *server) postUpdate(w http.ResponseWriter, r *http.Request) {
	code, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "SERVER ERROR")
		return
	}
	if err := s.codeRepo.Store("thecode", string(code)); err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "SERVER ERROR")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "OK"}`))
}

func (s *server) getPreview(w http.ResponseWriter, r *http.Request) {
	code, err := s.codeRepo.Load("thecode")
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "SERVER ERROR")
		return
	}
	if err := s.appHtml.ExecuteTemplate(w, "app.html", struct {
		Code string
	}{
		Code: code,
	}); err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "SERVER ERROR")
		return
	}
}

type inMemCodeStore struct {
	codes map[string]string
}

func (i *inMemCodeStore) Load(id string) (string, error) {
	if c, ok := i.codes[id]; ok {
		return c, nil
	}
	return "", fmt.Errorf("NOT FOUND")
}

func (i *inMemCodeStore) Store(id, code string) error {
	i.codes[id] = code
	return nil
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
	f1, err := www.ReadFile("www/app.html")
	if err != nil {
		log.Fatal(err)
	}
	t1, err := template.New("app.html").Parse(string(f1))
	if err != nil {
		log.Fatal(err)
	}
	s := &server{
		indexHtml: t,
		appHtml:   t1,
		codeRepo: &inMemCodeStore{
			codes: map[string]string{
				"thecode": "// TODO",
			},
		},
	}
	http.HandleFunc("/www/", s.getW3)
	http.HandleFunc("/preview/", s.getPreview)
	http.HandleFunc("/update/", s.postUpdate)
	http.HandleFunc("/", s.getRoot)
	if err := http.ListenAndServe(":3333", nil); err != nil {
		log.Fatal(err)
	}
}
