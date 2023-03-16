package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
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

	mux := http.NewServeMux()
	mux.HandleFunc("/", clean)

	srv := http.Server{
		Addr:    address,
		Handler: mux,
	}

	idleConnsClosed := make(chan struct{})
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-ctx.Done()

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("Serve on address %s", address)

	if err := srv.ListenAndServe(); err != nil {
		log.Printf("ListenAndServe: %v", err)
		cancel()
	}

	<-idleConnsClosed
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
