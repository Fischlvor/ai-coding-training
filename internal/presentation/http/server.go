package httpui

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

//go:embed templates/*.html
var templateFS embed.FS

type Server struct {
	addr     string
	template *template.Template
	mux      *http.ServeMux
}

func NewServer(port int) (*Server, error) {
	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}

	s := &Server{
		addr:     ":" + strconv.Itoa(port),
		template: tmpl,
		mux:      http.NewServeMux(),
	}
	s.routes()
	return s, nil
}

func (s *Server) routes() {
	s.mux.HandleFunc("/healthz", s.healthz)
	s.mux.HandleFunc("/", s.index)
	s.mux.HandleFunc("/admin", s.admin)
	s.mux.HandleFunc("/admin/configs", s.adminConfigs)
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

func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	_ = s.template.ExecuteTemplate(w, "index.html", map[string]any{
		"Title": "Config Center",
	})
}

func (s *Server) admin(w http.ResponseWriter, r *http.Request) {
	_ = s.template.ExecuteTemplate(w, "admin.html", map[string]any{
		"Title":        "Admin Console",
		"ActiveRoute":  "overview",
		"ServiceState": "Healthy",
	})
}

func (s *Server) adminConfigs(w http.ResponseWriter, r *http.Request) {
	_ = s.template.ExecuteTemplate(w, "configs.html", map[string]any{
		"Title":       "Configurations",
		"ActiveRoute": "configs",
	})
}
