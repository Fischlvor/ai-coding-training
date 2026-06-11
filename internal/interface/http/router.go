package httpui

func (s *Server) routes() {
	s.mux.HandleFunc("/healthz", s.pageHandler.Healthz)
	s.mux.HandleFunc("/", s.pageHandler.Index)
	s.mux.HandleFunc("/admin", s.pageHandler.Admin)
	s.mux.HandleFunc("/admin/configs", s.pageHandler.AdminConfigs)
	s.mux.HandleFunc("/admin/cluster", s.pageHandler.AdminCluster)
}
