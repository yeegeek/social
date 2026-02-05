package message

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler 消息处理器
type Handler struct {
	service Service
}

// NewHandler 创建消息处理器
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetConversations 获取对话列表
// @Summary 获取对话列表
// @Tags Messages
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} ConversationListResponse
// @Router /messages/conversations [get]
func (h *Handler) GetConversations(c *gin.Context) {
	userID := c.GetInt("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	response, err := h.service.GetConversations(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// GetMessages 获取消息列表
// @Summary 获取消息列表
// @Tags Messages
// @Accept json
// @Produce json
// @Param conversation_id path int true "对话ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(50)
// @Success 200 {object} MessageListResponse
// @Router /messages/conversations/{conversation_id}/messages [get]
func (h *Handler) GetMessages(c *gin.Context) {
	userID := c.GetInt("user_id")
	conversationID, _ := strconv.Atoi(c.Param("conversation_id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	response, err := h.service.GetMessages(conversationID, userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// SendMessage 发送消息
// @Summary 发送消息
// @Tags Messages
// @Accept json
// @Produce json
// @Param request body SendMessageRequest true "发送消息请求"
// @Success 200 {object} MessageResponse
// @Router /messages/send [post]
func (h *Handler) SendMessage(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.SendMessage(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// MarkAsRead 标记消息为已读
// @Summary 标记消息为已读
// @Tags Messages
// @Accept json
// @Produce json
// @Param message_id path int true "消息ID"
// @Success 200
// @Router /messages/{message_id}/read [post]
func (h *Handler) MarkAsRead(c *gin.Context) {
	userID := c.GetInt("user_id")
	messageID, _ := strconv.Atoi(c.Param("message_id"))

	err := h.service.MarkAsRead(messageID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Message marked as read"})
}

// MarkConversationAsRead 标记对话为已读
// @Summary 标记对话为已读
// @Tags Messages
// @Accept json
// @Produce json
// @Param conversation_id path int true "对话ID"
// @Success 200
// @Router /messages/conversations/{conversation_id}/read [post]
func (h *Handler) MarkConversationAsRead(c *gin.Context) {
	userID := c.GetInt("user_id")
	conversationID, _ := strconv.Atoi(c.Param("conversation_id"))

	err := h.service.MarkConversationAsRead(conversationID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Conversation marked as read"})
}

// RecallMessage 撤回消息
// @Summary 撤回消息
// @Tags Messages
// @Accept json
// @Produce json
// @Param message_id path int true "消息ID"
// @Success 200
// @Router /messages/{message_id}/recall [post]
func (h *Handler) RecallMessage(c *gin.Context) {
	userID := c.GetInt("user_id")
	messageID, _ := strconv.Atoi(c.Param("message_id"))

	err := h.service.RecallMessage(messageID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Message recalled"})
}

// SaveDraft 保存草稿
// @Summary 保存草稿
// @Tags Messages
// @Accept json
// @Produce json
// @Param request body SaveDraftRequest true "保存草稿请求"
// @Success 200 {object} MessageResponse
// @Router /messages/drafts [post]
func (h *Handler) SaveDraft(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req SaveDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.SaveDraft(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// GetDrafts 获取草稿列表
// @Summary 获取草稿列表
// @Tags Messages
// @Accept json
// @Produce json
// @Success 200 {array} MessageResponse
// @Router /messages/drafts [get]
func (h *Handler) GetDrafts(c *gin.Context) {
	userID := c.GetInt("user_id")

	response, err := h.service.GetDrafts(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// DeleteDraft 删除草稿
// @Summary 删除草稿
// @Tags Messages
// @Accept json
// @Produce json
// @Param draft_id path int true "草稿ID"
// @Success 200
// @Router /messages/drafts/{draft_id} [delete]
func (h *Handler) DeleteDraft(c *gin.Context) {
	userID := c.GetInt("user_id")
	draftID, _ := strconv.Atoi(c.Param("draft_id"))

	err := h.service.DeleteDraft(draftID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Draft deleted"})
}

// DeleteConversation 删除对话
// @Summary 删除对话
// @Tags Messages
// @Accept json
// @Produce json
// @Param conversation_id path int true "对话ID"
// @Success 200
// @Router /messages/conversations/{conversation_id} [delete]
func (h *Handler) DeleteConversation(c *gin.Context) {
	userID := c.GetInt("user_id")
	conversationID, _ := strconv.Atoi(c.Param("conversation_id"))

	err := h.service.DeleteConversation(conversationID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Conversation deleted"})
}
