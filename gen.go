//go:build ignore

package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

const (
	codeMirror = "https://cdnjs.cloudflare.com/ajax/libs/codemirror/6.65.7/"
	p5         = "https://cdnjs.cloudflare.com/ajax/libs/p5.js/1.7.0/"
)

var downloads = map[string]string{
	codeMirror + "codemirror.min.css":                "www/codemirror.min.css",
	codeMirror + "codemirror.min.js":                 "www/codemirror.min.js",
	codeMirror + "mode/javascript/javascript.min.js": "www/javascript.min.js",
	p5 + "p5.min.js":                                 "www/p5.min.js",
}

func main() {
	for k, v := range downloads {
		if err := downloadFile(k, v); err != nil {
			log.Fatal(err)
		}
	}
}

func downloadFile(src, dest string) error {
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	resp, err := client.Get(src)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}
