package parser

import (
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

	return metadata, buf.String(), nil
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
