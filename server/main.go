package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"secure-file-sync/server/handlers"
)

func main() {
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	// Register handlers
	http.HandleFunc("/upload", handlers.UploadHandler)
	http.HandleFunc("/download", handlers.DownloadHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	addr := ":" + *port
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

