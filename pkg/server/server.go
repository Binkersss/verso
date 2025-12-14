package server

import (
	"log"
	"net/http"
)

func Serve(addr, rootDir string) error {
	fs := http.FileServer(http.Dir(rootDir))
	http.Handle("/", fs)

	log.Printf("Server running at http://localhost%s\n", addr)
	log.Println("Watching for changes...")
	return http.ListenAndServe(addr, nil)
}
