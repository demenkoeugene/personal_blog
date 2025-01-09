package utils

import (
	"log"
	"net/http"
)

func HandleError(w http.ResponseWriter, message string, statusCode int, err error) {
	log.Println(message, ":", err)
	http.Error(w, message, statusCode)
}
