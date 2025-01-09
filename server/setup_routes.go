package server

import "personal_blog/handlers"

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/", handlers.GetArticleList)
	s.mux.HandleFunc("/article/", handlers.GetArticle)

	s.mux.HandleFunc("/new", s.auth.BasicAuth(handlers.AddNewArticle))
	s.mux.HandleFunc("/edit/", s.auth.BasicAuth(handlers.UpdateArticleHandler))
	s.mux.HandleFunc("/delete/", s.auth.BasicAuth(handlers.DeleteArticleHandler))
	s.mux.HandleFunc("/admin", s.auth.BasicAuth(handlers.GetDashboard))
}
