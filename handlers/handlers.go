package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"personal_blog/model"
	"personal_blog/services"
	"personal_blog/utils"
	"strconv"
	"strings"
	"time"
)

func GetDashboard(w http.ResponseWriter, r *http.Request) {
	handleArticleList(w, r, "./templates/dashboard.html")
}

func GetArticleList(w http.ResponseWriter, r *http.Request) {
	handleArticleList(w, r, "./templates/home.html")
}

func handleArticleList(w http.ResponseWriter, r *http.Request, templatePath string) {
	articles, err := services.GetAllArticles()
	if err != nil {
		utils.HandleError(w, "Failed to load articles", http.StatusInternalServerError, err)
		return
	}

	if err := utils.RenderTemplate(w, templatePath, articles); err != nil {
		utils.HandleError(w, "Failed to render template", http.StatusInternalServerError, err)
	}
}

func GetArticle(w http.ResponseWriter, r *http.Request) {
	articles, err := services.GetAllArticles()
	if err != nil {
		http.Error(w, "Failed to load articles", http.StatusInternalServerError)
		log.Println("Error fetching articles:", err)
		return
	}

	idStr := r.URL.Path[len("/article/"):]
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

func UpdateArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		path := r.URL.Path
		id := strings.TrimPrefix(path, "/edit/")

		if id == "" {
			http.Error(w, "Article ID missing", http.StatusBadRequest)
			return
		}

		filePath := fmt.Sprintf("articles/article%s.json", id)

		data, err := os.ReadFile(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				http.Error(w, "Article not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to load article", http.StatusInternalServerError)
			}
			return
		}

		// Розпарсити JSON у структуру Article
		var article model.Article
		err = json.Unmarshal(data, &article)
		if err != nil {
			http.Error(w, "Failed to parse article data", http.StatusInternalServerError)
			return
		}

		// Завантаження шаблону
		tmpl, err := template.ParseFiles("templates/updateArticle.html")
		if err != nil {
			http.Error(w, "Failed to load template", http.StatusInternalServerError)
			return
		}

		// Рендеринг шаблону з даними статті
		err = tmpl.Execute(w, article)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodPost {
		// Отримання ID зі шляху
		path := r.URL.Path
		id := strings.TrimPrefix(path, "/edit/")

		if id == "" {
			http.Error(w, "Article ID missing", http.StatusBadRequest)
			return
		}

		articleID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid article ID format", http.StatusBadRequest)
			return
		}

		// Отримання даних із форми
		err = r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		date := r.FormValue("date")

		// Валідація даних
		if title == "" || content == "" || date == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		// Оновлення статті
		updatedArticle := model.Article{
			ID:      articleID,
			Title:   title,
			Date:    date,
			Content: content,
		}

		// Формування шляху до JSON-файлу
		filePath := fmt.Sprintf("articles/article%s.json", id)

		// Запис оновленої статті у JSON
		data, err := json.MarshalIndent(updatedArticle, "", "  ")
		if err != nil {
			http.Error(w, "Failed to serialize article", http.StatusInternalServerError)
			return
		}

		err = os.WriteFile(filePath, data, 0644)
		if err != nil {
			http.Error(w, "Failed to save article", http.StatusInternalServerError)
			return
		}

		log.Printf("Article with ID %s updated successfully", id)

		// Перенаправлення після успішного збереження
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	id := strings.TrimPrefix(path, "/delete/") // Видаляємо "/delete/" з шляху

	if id == "" {
		http.Error(w, "Article ID missing", http.StatusBadRequest)
		return
	}

	// Формуємо шлях до файлу статті
	filePath := fmt.Sprintf("articles/article%s.json", id)

	// Спроба видалити файл
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Article not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete article", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Article with ID %s deleted successfully", id)

	// Перенаправлення на сторінку адміністрування після успішного видалення
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func AddNewArticle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Віддаємо HTML форму для створення нової статті
		http.ServeFile(w, r, "templates/newArticle.html")

	case http.MethodPost:
		// Парсимо дані з форми
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			log.Println("Error parsing form data:", err)
			return
		}

		// Отримуємо значення полів форми
		title := r.FormValue("title")
		content := r.FormValue("content")
		date := r.FormValue("date")

		// Валідація форми
		if title == "" || len(title) > 100 {
			http.Error(w, "Title is required and must be under 100 characters", http.StatusBadRequest)
			log.Println("Invalid title:", title)
			return
		}
		if content == "" {
			http.Error(w, "Content is required", http.StatusBadRequest)
			log.Println("Content is missing")
			return
		}
		_, err = time.Parse("2006-01-02", date)
		if err != nil {
			http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
			log.Println("Invalid date format:", date)
			return
		}

		// Завантажуємо існуючі статті
		articles, err := utils.FetchArticles("./articles")
		if err != nil {
			http.Error(w, "Failed to load articles", http.StatusInternalServerError)
			log.Println("Error fetching articles:", err)
			return
		}

		// Пошук максимального ID
		maxID := 0
		for _, article := range articles {
			if article.ID > maxID {
				maxID = article.ID
			}
		}

		// Створення нової статті
		newArticle := model.Article{
			ID:      maxID + 1,
			Title:   title,
			Content: content,
			Date:    date,
		}

		// Шлях для збереження статті
		filePath := "./articles/article" + strconv.Itoa(newArticle.ID) + ".json"

		// Збереження статті у файл
		err = utils.SaveArticle(filePath, newArticle)
		if err != nil {
			http.Error(w, "Failed to save article", http.StatusInternalServerError)
			log.Println("Error saving article:", err)
			return
		}

		// Перенаправлення користувача після успішного збереження
		http.Redirect(w, r, "/admin", http.StatusSeeOther)

	default:
		// Якщо метод не підтримується
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
