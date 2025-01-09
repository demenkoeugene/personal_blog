package utils

import (
	"fmt"
	"net/http"
	"personal_blog/model"
	"time"
)

func ParseAndValidateForm(r *http.Request) (model.Article, error) {
	if err := r.ParseForm(); err != nil {
		return model.Article{}, fmt.Errorf("failed to parse form data")
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	date := r.FormValue("date")

	if title == "" || len(title) > 100 {
		return model.Article{}, fmt.Errorf("title is required and must be under 100 characters")
	}

	if content == "" {
		return model.Article{}, fmt.Errorf("content is required")
	}

	if _, err := time.Parse("2006-01-02", date); err != nil {
		return model.Article{}, fmt.Errorf("invalid date format. Use YYYY-MM-DD")
	}

	return model.Article{
		Title:   title,
		Content: content,
		Date:    date,
	}, nil
}
