package verso

import (
	"log"

	"github.com/Binkersss/verso/pkg/builder"
	"github.com/Binkersss/verso/pkg/server"
	"github.com/Binkersss/verso/pkg/watcher"
)

type Config struct {
	ContentDir  string
	TemplateDir string
	StaticDir   string
	OutputDir   string
}

type Site struct {
	config  Config
	builder *builder.Builder
}

func New(cfg Config) *Site {
	return &Site{
		config:  cfg,
		builder: builder.New(cfg.ContentDir, cfg.TemplateDir, cfg.StaticDir, cfg.OutputDir),
	}
}

func (s *Site) Build() error {
	return s.builder.Build()
}

func (s *Site) Serve(addr string) error {
	if err := s.Build(); err != nil {
		return err
	}

	w, err := watcher.New(s.config.ContentDir, s.config.TemplateDir, s.config.StaticDir, func() {
		log.Println("Changes detected, rebuilding...")
		if err := s.Build(); err != nil {
			log.Printf("Build error: %v\n", err)
		}
	})
	if err != nil {
		return err
	}
	defer w.Close()

	return server.Serve(addr, s.config.OutputDir)
}
