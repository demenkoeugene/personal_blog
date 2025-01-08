package server

import (
	"fmt"
	"log"
	"net/http"
	"personal_blog/config"
	"personal_blog/handlers"
	"personal_blog/middleware"
)

type Server struct {
	mux    *http.ServeMux
	auth   *middleware.Authenticator
	config *config.Config
}

func NewServer() (*Server, error) {
	cfg, err := config.Load(config.GetConfigPath())
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	auth := middleware.NewAuthenticator(cfg)

	server := &Server{
		mux:    http.NewServeMux(),
		auth:   auth,
		config: cfg,
	}

	server.setupRoutes()
	return server, nil
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.GetPort())
	log.Printf("Starting server on %s", addr)

	server := &http.Server{
		Addr:    addr,
		Handler: s.mux,
	}

	return server.ListenAndServe()
}

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/", handlers.GetArticleList)
	s.mux.HandleFunc("/article/", handlers.GetArticle)

	s.mux.HandleFunc("/new", s.auth.BasicAuth(handlers.AddNewArticle))
	s.mux.HandleFunc("/edit/", s.auth.BasicAuth(handlers.UpdateArticleHandler))
	s.mux.HandleFunc("/delete/", s.auth.BasicAuth(handlers.DeleteArticleHandler))
	s.mux.HandleFunc("/admin", s.auth.BasicAuth(handlers.GetDashboard))
}
