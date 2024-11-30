package utils

import (
	"html/template"
	"log"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, templatePath string, data interface{}) error {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Println("Error loading template:", err)
		return err
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println("Error executing template:", err)
		return err
	}

	return nil
}

func RenderArticleTemplate(w http.ResponseWriter, templatePath string, article interface{}) error {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Println("Error loading template:", err)
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return err
	}

	if err := tmpl.Execute(w, article); err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return err
	}

	return nil
}
