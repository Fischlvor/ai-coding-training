package httpui

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"ai-coding-training/internal/interface/http/handler"
)

//go:embed templates/*.html
var templateFS embed.FS

type Server struct {
	addr        string
	template    *template.Template
	mux         *http.ServeMux
	pageHandler *handler.PageHandler
}

func NewServer(port int) (*Server, error) {
	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}

	s := &Server{
		addr:        ":" + strconv.Itoa(port),
		template:    tmpl,
		mux:         http.NewServeMux(),
		pageHandler: handler.NewPageHandler(tmpl),
	}
	s.routes()
	return s, nil
}

func (s *Server) Handler() http.Handler { return s.mux }

func (s *Server) Addr() string { return s.addr }

func (s *Server) Start(ctx context.Context) error {
	srv := &http.Server{Addr: s.addr, Handler: s.mux}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
