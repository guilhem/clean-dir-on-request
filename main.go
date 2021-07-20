package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var address string
var dir string

func init() {
	flag.StringVar(&address, "a", ":8080", "Address for the server to run on.")
	flag.StringVar(&dir, "d", "", "Path to delete.")
}

func main() {
	flag.Parse()
	if dir == "" {
		log.Panic("directory not set")
	}

	http.HandleFunc("/", clean)

	log.Printf("Serve on address %s", address)

	if err := http.ListenAndServe(address, nil); err != nil {
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
