package post

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler 动态处理器
type Handler struct {
	service Service
}

// NewHandler 创建动态处理器
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// CreatePost 创建动态
// @Summary 创建动态
// @Tags Posts
// @Accept json
// @Produce json
// @Param request body CreatePostRequest true "创建动态请求"
// @Success 200 {object} PostResponse
// @Router /posts [post]
func (h *Handler) CreatePost(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.CreatePost(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// GetPost 获取动态详情
// @Summary 获取动态详情
// @Tags Posts
// @Accept json
// @Produce json
// @Param post_id path int true "动态ID"
// @Success 200 {object} PostResponse
// @Router /posts/{post_id} [get]
func (h *Handler) GetPost(c *gin.Context) {
	userID := c.GetInt("user_id")
	postID, _ := strconv.Atoi(c.Param("post_id"))

	response, err := h.service.GetPost(postID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// GetPosts 获取动态列表
// @Summary 获取动态列表
// @Tags Posts
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} PostListResponse
// @Router /posts [get]
func (h *Handler) GetPosts(c *gin.Context) {
	userID := c.GetInt("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var uid *int
	if userID > 0 {
		uid = &userID
	}

	response, err := h.service.GetPosts(uid, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// GetUserPosts 获取用户的动态列表
// @Summary 获取用户的动态列表
// @Tags Posts
// @Accept json
// @Produce json
// @Param user_id path int true "用户ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} PostListResponse
// @Router /posts/user/{user_id} [get]
func (h *Handler) GetUserPosts(c *gin.Context) {
	currentUserID := c.GetInt("user_id")
	targetUserID, _ := strconv.Atoi(c.Param("user_id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	response, err := h.service.GetUserPosts(targetUserID, currentUserID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// UpdatePost 更新动态
// @Summary 更新动态
// @Tags Posts
// @Accept json
// @Produce json
// @Param post_id path int true "动态ID"
// @Param request body UpdatePostRequest true "更新动态请求"
// @Success 200 {object} PostResponse
// @Router /posts/{post_id} [put]
func (h *Handler) UpdatePost(c *gin.Context) {
	userID := c.GetInt("user_id")
	postID, _ := strconv.Atoi(c.Param("post_id"))

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.UpdatePost(postID, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// DeletePost 删除动态
// @Summary 删除动态
// @Tags Posts
// @Accept json
// @Produce json
// @Param post_id path int true "动态ID"
// @Success 200
// @Router /posts/{post_id} [delete]
func (h *Handler) DeletePost(c *gin.Context) {
	userID := c.GetInt("user_id")
	postID, _ := strconv.Atoi(c.Param("post_id"))

	err := h.service.DeletePost(postID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Post deleted"})
}

// LikePost 点赞动态
// @Summary 点赞动态
// @Tags Posts
// @Accept json
// @Produce json
// @Param post_id path int true "动态ID"
// @Success 200
// @Router /posts/{post_id}/like [post]
func (h *Handler) LikePost(c *gin.Context) {
	userID := c.GetInt("user_id")
	postID, _ := strconv.Atoi(c.Param("post_id"))

	err := h.service.LikePost(userID, postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Post liked"})
}

// UnlikePost 取消点赞
// @Summary 取消点赞
// @Tags Posts
// @Accept json
// @Produce json
// @Param post_id path int true "动态ID"
// @Success 200
// @Router /posts/{post_id}/unlike [post]
func (h *Handler) UnlikePost(c *gin.Context) {
	userID := c.GetInt("user_id")
	postID, _ := strconv.Atoi(c.Param("post_id"))

	err := h.service.UnlikePost(userID, postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Post unliked"})
}

// CreateComment 创建评论
// @Summary 创建评论
// @Tags Posts
// @Accept json
// @Produce json
// @Param request body CreateCommentRequest true "创建评论请求"
// @Success 200 {object} CommentResponse
// @Router /posts/comments [post]
func (h *Handler) CreateComment(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.CreateComment(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// GetComments 获取评论列表
// @Summary 获取评论列表
// @Tags Posts
// @Accept json
// @Produce json
// @Param post_id path int true "动态ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(50)
// @Success 200 {object} CommentListResponse
// @Router /posts/{post_id}/comments [get]
func (h *Handler) GetComments(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("post_id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	response, err := h.service.GetComments(postID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// DeleteComment 删除评论
// @Summary 删除评论
// @Tags Posts
// @Accept json
// @Produce json
// @Param comment_id path int true "评论ID"
// @Success 200
// @Router /posts/comments/{comment_id} [delete]
func (h *Handler) DeleteComment(c *gin.Context) {
	userID := c.GetInt("user_id")
	commentID, _ := strconv.Atoi(c.Param("comment_id"))

	err := h.service.DeleteComment(commentID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Comment deleted"})
}
