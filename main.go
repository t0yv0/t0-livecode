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

type server struct {
	indexHtml *template.Template
	appHtml   *template.Template
	codeRepo  codeRepo
}

func (s *server) getRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Location", "/p/"+string(s.codeRepo.Default()))
	w.WriteHeader(http.StatusMovedPermanently)
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

func (s *server) fail(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	fmt.Fprintf(w, "SERVER ERROR")
	log.Println(err)
}

func (s *server) pHandler(w http.ResponseWriter, r *http.Request) {
	pid, err := s.parsePID(r)
	if err != nil {
		s.fail(w, err)
		return
	}
	if _, err := s.codeRepo.Load(pid); err != nil {
		if err := s.codeRepo.Store(pid, starterCode); err != nil {
			s.fail(w, err)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	s.indexHtml.ExecuteTemplate(w, "index.html", struct{ CurrentProgram string }{
		CurrentProgram: string(pid),
	})
}

func (s *server) programHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	switch {
	case r.Method == http.MethodPost:
		err = s.programUpdate(w, r)
	case strings.HasSuffix(r.RequestURI, ".js"):
		err = s.programScript(w, r)
	default:
		err = s.programPage(w, r)
	}
	if err != nil {
		s.fail(w, err)
	}
}

func (s *server) programPage(w http.ResponseWriter, r *http.Request) error {
	pid, err := s.parsePID(r)
	if err != nil {
		return err
	}
	return s.appHtml.ExecuteTemplate(w, "app.html",
		struct{ Pid Pid }{Pid: pid})
}

func (s *server) programScript(w http.ResponseWriter, r *http.Request) error {
	pid, err := s.parsePID(r)
	if err != nil {
		return err
	}
	code, err := s.codeRepo.Load(pid)
	if err != nil {
		return err
	}
	contentType := "text/javascript"
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(code))
	return nil
}

func (s *server) programUpdate(w http.ResponseWriter, r *http.Request) error {
	pid, err := s.parsePID(r)
	if err != nil {
		return err
	}
	code, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err := s.codeRepo.Store(pid, string(code)); err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "OK"}`))
	return nil
}

func (*server) parsePID(r *http.Request) (Pid, error) {
	parts := strings.Split(r.RequestURI, "/")
	if len(parts) < 3 || (parts[1] != "program" && parts[1] != "p") {
		return "", fmt.Errorf("Parsing PID failed: invalid URI %q", r.RequestURI)
	}
	return ParsePid(parts[2])
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
		codeRepo:  newInMemCodeStore(),
	}
	http.HandleFunc("/www/", s.getW3)
	http.HandleFunc("/program/", s.programHandler)
	http.HandleFunc("/p/", s.pHandler)
	http.HandleFunc("/", s.getRoot)
	if err := http.ListenAndServe(":3333", nil); err != nil {
		log.Fatal(err)
	}
}
