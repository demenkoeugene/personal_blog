package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"personal_blog/model"
)

func main() {
	http.HandleFunc("/", getArticleList)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func getArticleList(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("./articles")
	check(err)

	var articles []model.Article

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			data, err := os.ReadFile("./articles/" + file.Name())
			if err != nil {
				log.Println("Error reading file:", file.Name(), err)
				continue
			}
			var article model.Article
			if err := json.Unmarshal(data, &article); err != nil {
				log.Println("Error parsing JSON:", file.Name(), err)
				continue
			}
			articles = append(articles, article)
		}
	}

	tmpl, err := template.ParseFiles("./templates/home.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		log.Println("Error loading template:", err)
		return
	}

	if err := tmpl.Execute(w, articles); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Println("Error rendering template:", err)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
