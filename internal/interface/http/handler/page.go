package handler

import (
	"html/template"
	"net/http"
)

type PageHandler struct {
	template *template.Template
}

func NewPageHandler(template *template.Template) *PageHandler {
	return &PageHandler{template: template}
}

func (h *PageHandler) Healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (h *PageHandler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	_ = h.template.ExecuteTemplate(w, "index.html", map[string]any{
		"Title": "Config Center",
	})
}

func (h *PageHandler) Admin(w http.ResponseWriter, _ *http.Request) {
	_ = h.template.ExecuteTemplate(w, "admin.html", map[string]any{
		"Title":        "Admin Console",
		"ActiveRoute":  "overview",
		"ServiceState": "Healthy",
	})
}

func (h *PageHandler) AdminConfigs(w http.ResponseWriter, _ *http.Request) {
	_ = h.template.ExecuteTemplate(w, "configs.html", map[string]any{
		"Title":       "Configurations",
		"ActiveRoute": "configs",
	})
}
