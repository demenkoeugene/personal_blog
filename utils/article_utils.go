package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"personal_blog/model"
)

func LoadArticleFromFile(id string) (*model.Article, error) {
	filePath := fmt.Sprintf("articles/article%s.json", id)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("article not found")
		}
		return nil, err
	}

	var article model.Article
	if err := json.Unmarshal(data, &article); err != nil {
		return nil, fmt.Errorf("failed to parse article data: %w", err)
	}

	return &article, nil
}

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

func ParseArticleForm(r *http.Request, articleID int) (*model.Article, error) {
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("failed to parse form data: %w", err)
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	date := r.FormValue("date")

	if title == "" || content == "" || date == "" {
		return nil, fmt.Errorf("all fields are required")
	}

	return &model.Article{
		ID:      articleID,
		Title:   title,
		Date:    date,
		Content: content,
	}, nil
}
