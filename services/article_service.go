package services

import (
	"fmt"
	"net/http"
	"personal_blog/model"
	"personal_blog/utils"
	"strings"
)

func GetAllArticles() ([]model.Article, error) {
	return utils.FetchArticles("./articles")
}

func CreateNewArticle(articles []model.Article, formData model.Article) model.Article {
	maxID := findMaxArticleID(articles)
	formData.ID = maxID + 1
	return formData
}

func SaveArticle(article model.Article) error {
	filePath := fmt.Sprintf("articles/article%d.json", article.ID)
	return utils.WriteJSONToFile(filePath, article)
}

func ParseArticleID(path string, prefix string) (string, error) {
	id := strings.TrimPrefix(path, prefix)
	if id == "" {
		return "", fmt.Errorf("article ID missing")
	}
	return id, nil
}

func ValidateRequest(r *http.Request, method string) error {
	if r.Method != method {
		return fmt.Errorf("invalid request method")
	}
	return nil
}

func HandleArticleList(w http.ResponseWriter, r *http.Request, templatePath string) {
	articles, err := GetAllArticles()
	if err != nil {
		utils.HandleError(w, "Failed to load articles", http.StatusInternalServerError, err)
		return
	}

	if err := utils.RenderTemplate(w, templatePath, articles); err != nil {
		utils.HandleError(w, "Failed to render template", http.StatusInternalServerError, err)
	}
}

func findMaxArticleID(articles []model.Article) int {
	maxID := 0
	for _, article := range articles {
		if article.ID > maxID {
			maxID = article.ID
		}
	}
	return maxID
}

func FindArticleByID(articles []model.Article, id int) *model.Article {
	for _, article := range articles {
		if article.ID == id {
			return &article
		}
	}
	return nil
}
