package v1

import (
	"KaldalisCMS/internal/api/middleware" // <-- New Import
	"KaldalisCMS/internal/service"
	"net/http"
	"strconv"
	"KaldalisCMS/internal/api/v1/dto"
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
	//对CreatePost进行DTO转换
	var createReq dto.CreatePostRequest
	
	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取 userID
	rawUserID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: user ID not found in context"})
		return
	}
	// Correctly handle the type conversion from context
	userIDInt, ok := rawUserID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error: user ID in context is not an integer"})
		return
	}
	userID := uint(userIDInt)

	newPost	:= createReq.ToEntity(userID) // <-- Pass userID to ToEntity()
	err := api.service.CreatePost(*newPost) // Note: ToEntity returns *entity.Post, so dereference it if CreatePost expects entity.Post
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, nil)
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
