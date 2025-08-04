package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	db "github.com/hexlet-components/go-gin-example/db/generated"
)

func SetupRouter(database *sql.DB) *gin.Engine {
	queries := db.New(database)
	h := NewArticleHandler(queries)

	r := gin.Default()
	articles := r.Group("/articles")
	h.Register(articles)

	return r
}
