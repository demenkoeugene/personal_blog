package handlers

import (
	"net/http"
	"personal_blog/services"
	"personal_blog/utils"
	"strconv"
)

func GetArticleList(w http.ResponseWriter, r *http.Request) {
	services.HandleArticleList(w, r, "./templates/home.html")
}

func GetArticle(w http.ResponseWriter, r *http.Request) {
	id, err := services.ParseArticleID(r.URL.Path, "/article/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	articles, err := services.GetAllArticles()
	if err != nil {
		utils.HandleError(w, "Failed to load articles", http.StatusInternalServerError, err)
		return
	}

	articleID, err := strconv.Atoi(id)
	if err != nil {
		utils.HandleError(w, "Invalid article ID", http.StatusBadRequest, err)
		return
	}

	foundArticle := services.FindArticleByID(articles, articleID)
	if foundArticle == nil {
		services.LogArticleOperation("not found", id)
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	if err := utils.RenderArticleTemplate(w, "./templates/articlepage.html", foundArticle); err != nil {
		services.LogArticleOperation("failed to render", id)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}
