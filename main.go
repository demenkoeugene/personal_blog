package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"personal_blog/model"
	"personal_blog/utils"
	"strconv"
)

func initConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	} else {
		log.Println("Environment variables loaded successfully")
	}
}
func main() {
	initConfig()
	http.HandleFunc("/", getArticleList)
	http.HandleFunc("/article/{id}", getArticle)
	http.HandleFunc("/admin", basicAuth(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is the protected handler"))
	},
	))
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

func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usernameENV := os.Getenv("ADMIN_USERNAME")
		passwordENV := os.Getenv("ADMIN_PASSWORD")
		if usernameENV == "" || passwordENV == "" {
			log.Fatal("Environment variables ADMIN_USERNAME and ADMIN_PASSWORD must be set")
		}

		// Extract the username and password from the Authorization header.
		username, password, ok := r.BasicAuth()
		if ok {
			// Hash both provided and expected credentials using SHA-256.
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(usernameENV))
			expectedPasswordHash := sha256.Sum256([]byte(passwordENV))

			// Use ConstantTimeCompare to compare the hashes.
			usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

			// If both the username and password match, proceed to the next handler.
			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		// If authentication fails, return a 401 Unauthorized response with the
		// WWW-Authenticate header to prompt for basic authentication.
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}
