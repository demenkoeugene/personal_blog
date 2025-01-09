package handlers

import (
	"net/http"
	"personal_blog/services"
	"personal_blog/utils"
	"strconv"
)

func UpdateArticleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := services.ParseArticleID(r.URL.Path, "/edit/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGetUpdateArticle(w, r, id)
	case http.MethodPost:
		handlePostUpdateArticle(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetUpdateArticle(w http.ResponseWriter, r *http.Request, id string) {
	article, err := utils.LoadArticleFromFile(id)
	if err != nil {
		utils.HandleError(w, "Failed to load article", http.StatusInternalServerError, err)
		return
	}

	if err := utils.RenderTemplate(w, "templates/updateArticle.html", article); err != nil {
		utils.HandleError(w, "Failed to render template", http.StatusInternalServerError, err)
	}
}

func handlePostUpdateArticle(w http.ResponseWriter, r *http.Request, id string) {
	articleID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid article ID format", http.StatusBadRequest)
		return
	}

	article, err := utils.ParseArticleForm(r, articleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := services.SaveArticle(*article); err != nil {
		utils.HandleError(w, "Failed to save article", http.StatusInternalServerError, err)
		return
	}

	services.LogArticleOperation("updated successfully", id)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
