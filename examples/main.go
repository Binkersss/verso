package main

import (
	"github.com/Binkersss/verso"
	"log"
)

func main() {
	site := verso.New(verso.Config{
		ContentDir:  "content",
		TemplateDir: "templates",
		StaticDir:   "static",
		OutputDir:   "dist",
		SiteTitle:   "Nathaniel Chappelle",
	})

	if err := site.Serve(":3000"); err != nil {
		log.Fatal(err)
	}
}
