package main

import (
	"log"
	"net/http"
	"personal_blog/model"
	"personal_blog/utils"
	"strconv"
)

func main() {
	http.HandleFunc("/", getArticleList)
	http.HandleFunc("/article/{id}", getArticle)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func getArticleList(w http.ResponseWriter, r *http.Request) {
	articles, err := utils.FetchArticles("./articles")
	if err != nil {
		http.Error(w, "Failed to load articles", http.StatusInternalServerError)
		log.Println("Error fetching articles:", err)
		return
	}

	if err := utils.RenderTemplate(w, "./templates/home.html", articles); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Println("Error rendering template:", err)
	}
}

func getArticle(w http.ResponseWriter, r *http.Request) {
	articles, err := utils.FetchArticles("./articles")
	if err != nil {
		http.Error(w, "Failed to load articles", http.StatusInternalServerError)
		log.Println("Error fetching articles:", err)
		return
	}

	idStr := r.URL.Path[len("/article/"):] // Витягування частини шляху після "/article/"
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusNotFound)
		log.Println("Error parsing article ID:", err)
		return
	}

	var foundArticle *model.Article
	for _, article := range articles {
		if article.ID == id {
			foundArticle = &article
			break
		}
	}

	if foundArticle == nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		log.Println("Article not found with ID:", id)
		return
	}

	if err := utils.RenderArticleTemplate(w, "./templates/articlepage.html", foundArticle); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Println("Error rendering template:", err)
	}
}
