package handlers

import (
	"net/http"
	"personal_blog/services"
)

func GetDashboard(w http.ResponseWriter, r *http.Request) {
	services.HandleArticleList(w, r, "./templates/dashboard.html")
}
