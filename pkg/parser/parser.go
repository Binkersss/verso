package parser

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

type Page struct {
	Route    string                 `json:"route"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata"`
}

type Parser struct {
	md goldmark.Markdown
}

func New() *Parser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			extension.Strikethrough,
			extension.Linkify,
			extension.TaskList,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)

	return &Parser{md: md}
}

func (p *Parser) Parse(content string) (map[string]interface{}, string, error) {
	metadata, body := p.parseFrontmatter(content)

	var buf strings.Builder
	if err := p.md.Convert([]byte(body), &buf); err != nil {
		return nil, "", err
	}

	htmlContent := buf.String()

	metaHeader := buildMetadataHeader(metadata)
	if metaHeader != "" {
		htmlContent = metaHeader + htmlContent
	}

	return metadata, htmlContent, nil
}

func (p *Parser) parseFrontmatter(content string) (map[string]interface{}, string) {
	metadata := make(map[string]interface{})

	if !strings.HasPrefix(content, "---\n") {
		return metadata, content
	}

	parts := strings.SplitN(content[4:], "\n---\n", 2)
	if len(parts) != 2 {
		return metadata, content
	}

	_ = yaml.Unmarshal([]byte(parts[0]), &metadata)
	return metadata, strings.TrimSpace(parts[1])
}

func buildMetadataHeader(metadata map[string]interface{}) string {
	var parts []string

	// Add title as h1
	if title, ok := metadata["title"].(string); ok && title != "" {
		parts = append(parts, fmt.Sprintf(`<h1>%s</h1>`, title))
	}

	// Add date and authors in metadata section
	var metaParts []string
	if date, ok := metadata["date"].(string); ok && date != "" {
		metaParts = append(metaParts, fmt.Sprintf(`<time class="page-date">%s</time>`, date))
	}

	if authors := getAuthors(metadata); authors != "" {
		metaParts = append(metaParts, fmt.Sprintf(`<span class="page-authors">%s</span>`, authors))
	}

	if len(metaParts) > 0 {
		parts = append(parts, `<div class="page-meta-info">`+strings.Join(metaParts, " • ")+`</div>`)
	}

	if len(parts) == 0 {
		return ""
	}

	return `<div class="page-header">` + strings.Join(parts, "") + `</div>` + "\n"
}

func getAuthors(metadata map[string]interface{}) string {
	// Handle single author
	if author, ok := metadata["author"].(string); ok && author != "" {
		return author
	}

	// Handle multiple authors as array
	if authors, ok := metadata["authors"].([]interface{}); ok {
		var authorStrs []string
		for _, a := range authors {
			if s, ok := a.(string); ok {
				authorStrs = append(authorStrs, s)
			}
		}
		if len(authorStrs) > 0 {
			return strings.Join(authorStrs, ", ")
		}
	}

	return ""
}
