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
		"Nodes": []map[string]any{
			{"Name": "node-a", "Leader": true, "State": "Healthy"},
			{"Name": "node-b", "Leader": false, "State": "Healthy"},
		},
		"LatestRelease": map[string]any{
			"App":     "example-service",
			"Env":     "prod",
			"Group":   "default",
			"Version": "v12",
			"Status":  "Published",
		},
	})
}

func (h *PageHandler) AdminConfigs(w http.ResponseWriter, _ *http.Request) {
	_ = h.template.ExecuteTemplate(w, "configs.html", map[string]any{
		"Title":       "Configurations",
		"ActiveRoute": "configs",
		"Configs": []map[string]any{
			{"Name": "example-service", "Environment": "prod", "Group": "default", "Version": "v12", "Status": "Published"},
			{"Name": "example-service", "Environment": "staging", "Group": "default", "Version": "v13-draft", "Status": "Draft"},
		},
	})
}

func (h *PageHandler) AdminCluster(w http.ResponseWriter, _ *http.Request) {
	_ = h.template.ExecuteTemplate(w, "cluster.html", map[string]any{
		"Title":       "Cluster Status",
		"ActiveRoute": "cluster",
		"Nodes": []map[string]any{
			{"Name": "node-a", "Leader": true, "State": "Healthy", "Address": "10.0.0.11:9000"},
			{"Name": "node-b", "Leader": false, "State": "Healthy", "Address": "10.0.0.12:9000"},
		},
	})
}
