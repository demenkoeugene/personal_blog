package utils

import (
	"encoding/json"
	"log"
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

func SaveArticle(filePath string, article model.Article) error {
	data, err := json.MarshalIndent(article, "", "  ")
	if err != nil {
		log.Printf("Error marshalling article %+v: %v", article, err)
		return err
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		log.Printf("Error writing to file %s: %v", filePath, err)
		return err
	}

	log.Println("Article saved successfully to", filePath)
	return nil
}
