package v1

import (
	"KaldalisCMS/internal/core/entity"
	"KaldalisCMS/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostAPI struct {
	service *service.PostService
}

func NewPostAPI(service *service.PostService) *PostAPI {
	return &PostAPI{service: service}
}

func (api *PostAPI) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/posts", api.GetPosts)
	group.POST("/posts", api.CreatePost)
	group.GET("/posts/:id", api.GetPostByID)
	group.PUT("/posts/:id", api.UpdatePost)
	group.DELETE("/posts/:id", api.DeletePost)
}

func (api *PostAPI) GetPosts(c *gin.Context) {
	posts, err := api.service.GetAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (api *PostAPI) GetPostByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	post, err := api.service.GetPostByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}

func (api *PostAPI) CreatePost(c *gin.Context) {
	var newPost entity.Post
	if err := c.ShouldBindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := api.service.CreatePost(newPost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, nil)
}

func (api *PostAPI) UpdatePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	var updatedPost entity.Post
	if err := c.ShouldBindJSON(&updatedPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = api.service.UpdatePost(id, updatedPost)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (api *PostAPI) DeletePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	err = api.service.DeletePost(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
