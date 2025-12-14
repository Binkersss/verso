package builder

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Binkersss/verso/pkg/client"
	"github.com/Binkersss/verso/pkg/parser"
)

type ContentManifest struct {
	Pages map[string]parser.Page `json:"pages"`
}

type Builder struct {
	contentDir  string
	templateDir string
	staticDir   string
	outputDir   string
	parser      *parser.Parser
}

func New(contentDir, templateDir, staticDir, outputDir string) *Builder {
	return &Builder{
		contentDir:  contentDir,
		templateDir: templateDir,
		staticDir:   staticDir,
		outputDir:   outputDir,
		parser:      parser.New(),
	}
}

func (b *Builder) Build() error {
	log.Println("Building site...")

	if err := os.MkdirAll(b.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	pages, err := b.parseContent()
	if err != nil {
		return fmt.Errorf("failed to parse content: %w", err)
	}

	if err := b.writeClient(pages); err != nil {
		return fmt.Errorf("failed to write client: %w", err)
	}

	if err := b.writeManifest(pages); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	if err := b.copyTemplate(); err != nil {
		return fmt.Errorf("failed to copy template: %w", err)
	}

	if err := b.copyStatic(); err != nil {
		return fmt.Errorf("failed to copy static: %w", err)
	}

	log.Println("Build complete!")
	return nil
}

func (b *Builder) parseContent() (map[string]parser.Page, error) {
	pages := make(map[string]parser.Page)

	err := filepath.Walk(b.contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		metadata, html, err := b.parser.Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		relPath, _ := filepath.Rel(b.contentDir, path)
		route := strings.TrimSuffix(relPath, ".md")
		route = strings.ReplaceAll(route, string(os.PathSeparator), "/")

		pages[route] = parser.Page{
			Route:    route,
			Content:  html,
			Metadata: metadata,
		}

		log.Printf("Parsed: %s -> /%s\n", path, route)
		return nil
	})

	return pages, err
}

func (b *Builder) writeManifest(pages map[string]parser.Page) error {
	manifest := ContentManifest{Pages: pages}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(b.outputDir, "content.json")
	return os.WriteFile(path, data, 0644)
}

func (b *Builder) copyTemplate() error {
	src := filepath.Join(b.templateDir, "index.html")
	dst := filepath.Join(b.outputDir, "index.html")

	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, 0644)
}

func (b *Builder) writeClient(pages map[string]parser.Page) error {
	// Convert parser.Page to client.Page
	clientPages := make(map[string]client.Page)
	for route, page := range pages {
		clientPages[route] = client.Page{
			Route:    page.Route,
			Content:  page.Content,
			Metadata: page.Metadata,
		}
	}

	// Generate client JS
	clientConfig := client.DefaultConfig()
	js, err := client.Generate(client.ContentManifest{Pages: clientPages}, clientConfig)
	if err != nil {
		return err
	}

	// Write to dist/app.js
	jsPath := filepath.Join(b.outputDir, "app.js")
	if err := os.WriteFile(jsPath, []byte(js), 0644); err != nil {
		return err
	}

	return nil
}

func (b *Builder) copyStatic() error {
	if _, err := os.Stat(b.staticDir); os.IsNotExist(err) {
		return nil
	}

	return filepath.Walk(b.staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(b.staticDir, path)
		destPath := filepath.Join(b.outputDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		return copyFile(path, destPath)
	})
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
