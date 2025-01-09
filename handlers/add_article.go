package handlers

import (
	"log"
	"net/http"
	"personal_blog/services"
	"personal_blog/utils"
)

func AddNewArticle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetNewArticle(w, r)
	case http.MethodPost:
		handlePostNewArticle(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetNewArticle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/newArticle.html")
}

func handlePostNewArticle(w http.ResponseWriter, r *http.Request) {
	var formData, err = utils.ParseAndValidateForm(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Form validation error:", err)
		return
	}

	articles, err := utils.FetchArticles("./articles")
	if err != nil {
		http.Error(w, "Failed to load articles", http.StatusInternalServerError)
		log.Println("Error fetching articles:", err)
		return
	}

	newArticle := services.CreateNewArticle(articles, formData)

	if err := services.SaveArticle(newArticle); err != nil {
		http.Error(w, "Failed to save article", http.StatusInternalServerError)
		log.Println("Error saving article:", err)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
