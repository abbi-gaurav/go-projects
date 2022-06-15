package main

import (
	app "example.com/react-go-embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var uiFS fs.FS

func init() {
	var err error
	uiFS, err = fs.Sub(app.UI, "_ui/build")
	if err != nil {
		log.Fatal("failed to get ui fs", err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api", handleApi)
	mux.HandleFunc("/", handleStatic)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Println("server failed...", err)
	}
}

func handleStatic(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Error(writer, http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed)
		return
	}
	path := filepath.Clean(request.URL.Path)
	if path == "/" {
		path = "index.html"
	}
	path = strings.TrimPrefix(path, "/")

	file, err := uiFS.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("file", path, "not found:", err)
			http.NotFound(writer, request)
			return
		}
		log.Println("file", path, "cannot be read:", err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	contentType := mime.TypeByExtension(filepath.Ext(path))
	writer.Header().Set("Content-Type", contentType)
	if strings.HasPrefix(path, "static/") {
		writer.Header().Set("Cache-Control", "public, max-age=31536000")
	}
	stat, err := file.Stat()
	if err == nil && stat.Size() > 0 {
		writer.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	}

	num, _ := io.Copy(writer, file)
	log.Println("file", path, "copied", num, "bytes")
}

func handleApi(writer http.ResponseWriter, request *http.Request) {

}

func handleHealth(writer http.ResponseWriter, request *http.Request) {

}
