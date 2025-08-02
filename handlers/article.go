package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	db "github.com/hexlet-components/go-gin-example/db/generated"
)

type ArticleParams struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

type ArticleHandler struct {
	Queries *db.Queries
}

func NewArticleHandler(queries *db.Queries) *ArticleHandler {
	return &ArticleHandler{Queries: queries}
}

func (h *ArticleHandler) Register(rg *gin.RouterGroup) {
	rg.POST("", h.Create)
	rg.GET("/:id", h.Get)
	rg.GET("", h.List)
	rg.PUT("/:id", h.Update)
	rg.DELETE("/:id", h.Delete)
}

func (h *ArticleHandler) Create(c *gin.Context) {
	input, ok := h.parseAndValidateParams(c)
	if !ok {
		return
	}

	article, err := h.Queries.CreateArticle(c, input.Name)
	if err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(http.StatusCreated, article)
}

func (h *ArticleHandler) Get(c *gin.Context) {
	id, err := h.parseID(c)
	if err != nil {
		badRequest(c, err)
		return
	}

	article, err := h.Queries.GetArticle(c, id)
	if err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) List(c *gin.Context) {
	articles, err := h.Queries.ListArticles(c)
	if err != nil {
		handleDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, articles)
}

func (h *ArticleHandler) Update(c *gin.Context) {
	id, err := h.parseID(c)
	if err != nil {
		badRequest(c, err)
		return
	}

	input, ok := h.parseAndValidateParams(c)
	if !ok {
		return
	}

	updateParams := db.UpdateArticleParams{
		ID:   id,
		Name: input.Name,
	}

	article, err := h.Queries.UpdateArticle(c, updateParams)
	if err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) Delete(c *gin.Context) {
	id, err := h.parseID(c)
	if err != nil {
		badRequest(c, err)
		return
	}

	err = h.Queries.DeleteArticle(c, id)
	if err != nil {
		handleDBError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ArticleHandler) parseID(c *gin.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, ErrorInvalidID
	}
	return id, nil
}

func (h *ArticleHandler) parseAndValidateParams(c *gin.Context) (ArticleParams, bool) {
	var params ArticleParams

	if err := c.ShouldBindJSON(&params); err != nil {
		badRequest(c, err)
		return ArticleParams{}, false
	}

	params.Name = strings.TrimSpace(params.Name)
	if len(params.Name) == 0 {
		badRequest(c, ErrorNameEmpty)
		return ArticleParams{}, false
	}

	return params, true
}
