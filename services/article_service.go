package services

import (
	"personal_blog/model"
	"personal_blog/utils"
)

func GetAllArticles() ([]model.Article, error) {
	return utils.FetchArticles("./articles")
}
