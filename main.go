package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var port string
var dir string

func init() {
	flag.StringVar(&port, "p", "8080", "Port for the server to run on.")
	flag.StringVar(&dir, "d", "", "Path to delete.")
}

func main() {
	flag.Parse()
	if dir == "" {
		log.Panic("directory not set")
	}

	http.HandleFunc("/", clean)

	log.Printf("Serve on port %v", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Panicf("ListenAndServe: %v", err)
	}
}

func clean(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("can't read dir: %v", err)
		http.Error(w, "Can't read dir", http.StatusInternalServerError)
		return
	}

	for _, file := range files {
		if err := os.RemoveAll(filepath.Join(dir, file.Name())); err != nil {
			http.Error(w, "Can't suppress dir", http.StatusInternalServerError)
			// try to suppress as much as we can
		}
	}
}
