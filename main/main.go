package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	http.Handle("/", new(RequestHandler))
	http.Handle("/download", new(DownloadHandler))

	http.Handle("/img/", new(StaticHandler))
	http.Handle("/css/", new(StaticHandler))
	http.Handle("/js/", new(StaticHandler))

	http.ListenAndServe(":8000", nil)
}

// RequestHandler implements the http.Handler interface
type RequestHandler struct {
	http.Handler
}

// ServeHome serves files from public directory
func (handler *RequestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("%s - %s\n", req.Method, req.URL.Path)

	path := "public/home.html"
	file, err := os.Open(path)

	if err == nil {
		bufferedReader := bufio.NewReader(file)

		w.Header().Set("Content-Type", "text/html")
		bufferedReader.WriteTo(w)
	} else {
		w.WriteHeader(404)
		w.Write([]byte("404 - " + http.StatusText(404)))
	}
}

// DownloadHandler will download a file from the user specified link
type DownloadHandler struct {
	http.Handler
}

func (handler *DownloadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	stdout, err := exec.Command("youtube-dl", "--get-filename", "--restrict-filename", req.FormValue("url")).Output()
	videoName := strings.TrimSpace(string(stdout))
	musicName := videoName[:len(videoName)-1] + "3"

	exec.Command("youtube-dl", "--restrict-filename", req.FormValue("url")).Run()

	exec.Command("ffmpeg", "-i", videoName, "-vn", musicName).Run()

	filename, err := os.Open(musicName)

	if err == nil {
		br := bufio.NewReader(filename)
		w.Header().Set("Content-Disposition", "attachment; filename="+musicName)
		w.Header().Set("Content-Type", "audio/mp3")

		br.WriteTo(w)
	} else {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
	}

	exec.Command("rm", "-rf", videoName, musicName).Run()
}

// StaticHandler will download a file from the user specified link
type StaticHandler struct {
	http.Handler
}

func (handler *StaticHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("%s - %s\n", req.Method, req.URL.Path)

	path := "public/" + req.URL.Path
	file, err := os.Open(path)

	if err == nil {
		bufferedReader := bufio.NewReader(file)
		var contentType string

		if strings.HasSuffix(path, ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(path, ".css") {
			contentType = "text/css"
		} else if strings.HasSuffix(path, ".js") {
			contentType = "application/javascript"
		} else {
			contentType = "text/plain"
		}

		fmt.Println("contentType:")
		fmt.Println(contentType)

		w.Header().Set("Content-Type", contentType)
		bufferedReader.WriteTo(w)
	} else {
		w.WriteHeader(404)
		w.Write([]byte("404 - " + http.StatusText(404)))
	}
}
