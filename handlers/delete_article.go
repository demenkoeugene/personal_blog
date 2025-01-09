package handlers

import (
	"fmt"
	"net/http"
	"personal_blog/services"
	"personal_blog/utils"
)

func DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	if err := services.ValidateRequest(r, http.MethodGet); err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	id, err := services.ParseArticleID(r.URL.Path, "/delete/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filePath := fmt.Sprintf("articles/article%s.json", id)

	if err := utils.RemoveFile(filePath); err != nil {
		utils.HandleFileRemovalError(err, w)
		return
	}

	services.LogArticleOperation("deleted successfully", id)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
