package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.Handle("/", new(RequestHandler))

	http.ListenAndServe(":8000", nil)
}

// RequestHandler implements the http.Handler interface
type RequestHandler struct {
	http.Handler
}

// ServeHome serves files from public directory
func (handler *RequestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("%s - %s\n", req.Method, req.URL.Path)

	path := "public/" + req.URL.Path
	file, err := os.Open(path)

	if err == nil {
		bufferedReader := bufio.NewReader(file)
		var contentType string

		if strings.HasSuffix(path, ".css") {
			contentType = "text/css"
		} else if strings.HasSuffix(path, ".html") {
			contentType = "text/html"
		} else if strings.HasSuffix(path, ".js") {
			contentType = "application/javascript"
		} else if strings.HasSuffix(path, ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(path, ".mp4") {
			contentType = "video/mp4"
		} else {
			contentType = "text/plain"
		}

		w.Header().Add("Content Type", contentType)
		bufferedReader.WriteTo(w)
	} else {
		w.WriteHeader(404)
		w.Write([]byte("404 - " + http.StatusText(404)))
	}
}
