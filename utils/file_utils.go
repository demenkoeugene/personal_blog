package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func WriteJSONToFile(filePath string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func RemoveFile(filePath string) error {
	return os.Remove(filePath)
}

func HandleFileRemovalError(err error, w http.ResponseWriter) {
	if os.IsNotExist(err) {
		http.Error(w, "Article not found", http.StatusNotFound)
	} else {
		http.Error(w, "Failed to delete article", http.StatusInternalServerError)
	}
}
