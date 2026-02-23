package v1

import (
	"KaldalisCMS/internal/api/middleware" // <-- New Import
	"KaldalisCMS/internal/api/v1/dto"
	"KaldalisCMS/internal/core"
	"context" // Added
	"errors"  // Added
	"net/http"
	"strconv"
	"time" // 引入 time 包

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
	// 读取操作通常较快，设置 5 秒超时
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	posts, err := api.service.GetAllPosts(ctx)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Get posts timed out"})
			return
		}
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

	// 读取详情操作，设置 5 秒超时
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	post, err := api.service.GetPostByID(ctx, uint(id))
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Get post timed out"})
			return
		}
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

	// 设置超时时间：例如 10 秒
	// 注意：媒体同步可能涉及文件操作或复杂的解析，给足够的时间
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if err := api.service.CreatePost(ctx, *newPost); err != nil {
		// 检查是否是超时错误
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Create post timed out"})
			return
		}
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

	// 同样为更新操作设置超时
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	err = api.service.UpdatePost(ctx, uint(id), updatedPartialPost)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Update post timed out"})
			return
		}
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

	// 删除操作可能涉及级联，设置 10 秒超时
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	err = api.service.DeletePost(ctx, uint(id))
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Delete post timed out"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
