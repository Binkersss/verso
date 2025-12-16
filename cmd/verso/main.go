package main

import (
	"flag"
	"log"
	"os"

	"github.com/Binkersss/verso"
)

func main() {
	var (
		build       = flag.Bool("build", false, "Build the site")
		serve       = flag.String("serve", "", "Serve the site (e.g., :3000)")
		contentDir  = flag.String("content", "content", "Content directory")
		templateDir = flag.String("templates", "templates", "Templates directory")
		staticDir   = flag.String("static", "static", "Static assets directory")
		outputDir   = flag.String("output", "dist", "Output directory")
		siteTitle   = flag.String("title", "Verso Site", "Default Site Title")
	)
	flag.Parse()

	cfg := verso.Config{
		ContentDir:  *contentDir,
		TemplateDir: *templateDir,
		StaticDir:   *staticDir,
		OutputDir:   *outputDir,
		SiteTitle:   *siteTitle,
	}

	site := verso.New(cfg)

	if *serve != "" {
		if err := site.Serve(*serve); err != nil {
			log.Fatal(err)
		}
	} else if *build {
		if err := site.Build(); err != nil {
			log.Fatal(err)
		}
	} else {
		flag.Usage()
		os.Exit(1)
	}
}
