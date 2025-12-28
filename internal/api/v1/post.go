package v1

import (
	"KaldalisCMS/internal/api/middleware" // <-- New Import
	"KaldalisCMS/internal/api/v1/dto"
	"KaldalisCMS/internal/core"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostAPI struct {
	service core.PostService
}

func NewPostAPI(service core.PostService) *PostAPI {
	return &PostAPI{service: service}
}

//目前在router的protected中进行注册了
//func (api *PostAPI) RegisterRoutes(group *gin.RouterGroup) {
//	group.GET("/posts", api.GetPosts)
//	group.POST("/posts", api.CreatePost)
//	group.GET("/posts/:id", api.GetPostByID)
//	group.PUT("/posts/:id", api.UpdatePost)
//	group.DELETE("/posts/:id", api.DeletePost)
//}

func (api *PostAPI) GetPosts(c *gin.Context) {
	posts, err := api.service.GetAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.ToPostListResponse(posts))
}

func (api *PostAPI) GetPostByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	post, err := api.service.GetPostByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.ToPostResponse(&post))
}

func (api *PostAPI) CreatePost(c *gin.Context) {
	var createReq dto.CreatePostRequest

	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: user ID not found in context"})
		return
	}

	// 业务逻辑处理
	newPost := createReq.ToEntity(userID)

	if err := api.service.CreatePost(*newPost); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully"})
}

func (api *PostAPI) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}
	var req dto.UpdatePostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedPartialPost := req.ToEntity()

	err = api.service.UpdatePost(uint(id), updatedPartialPost)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (api *PostAPI) DeletePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	err = api.service.DeletePost(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
