package utils

import (
	"encoding/json"
	"os"
	"personal_blog/model"
)

func ReadArticle(filePath string) (model.Article, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return model.Article{}, err
	}

	var article model.Article
	if err := json.Unmarshal(data, &article); err != nil {
		return model.Article{}, err
	}

	return article, nil
}
