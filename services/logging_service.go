package services

import "log"

func LogArticleOperation(operation string, id string) {
	log.Printf("Article with ID %s %s", id, operation)
}
