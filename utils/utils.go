package utils

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"personal_blog/model"
)

func FetchArticles(directory string) ([]model.Article, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var articles []model.Article
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		article, err := ReadArticle(filepath.Join(directory, file.Name()))
		if err != nil {
			log.Println("Error reading or parsing article:", file.Name(), err)
			continue
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func HandleError(w http.ResponseWriter, message string, statusCode int, err error) {
	log.Println(message, ":", err)
	http.Error(w, message, statusCode)
}
