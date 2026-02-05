package notification

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler 通知处理器
type Handler struct {
	service Service
}

// NewHandler 创建通知处理器
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetNotifications 获取通知列表
// @Summary 获取通知列表
// @Tags Notifications
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} NotificationListResponse
// @Router /notifications [get]
func (h *Handler) GetNotifications(c *gin.Context) {
	userID := c.GetInt("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	response, err := h.service.GetNotifications(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// GetUnreadCount 获取未读通知数量
// @Summary 获取未读通知数量
// @Tags Notifications
// @Accept json
// @Produce json
// @Success 200 {object} map[string]int64
// @Router /notifications/unread-count [get]
func (h *Handler) GetUnreadCount(c *gin.Context) {
	userID := c.GetInt("user_id")

	count, err := h.service.GetUnreadCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"unread_count": count}})
}

// MarkAsRead 标记通知为已读
// @Summary 标记通知为已读
// @Tags Notifications
// @Accept json
// @Produce json
// @Param notification_id path int true "通知ID"
// @Success 200
// @Router /notifications/{notification_id}/read [post]
func (h *Handler) MarkAsRead(c *gin.Context) {
	userID := c.GetInt("user_id")
	notificationID, _ := strconv.Atoi(c.Param("notification_id"))

	err := h.service.MarkAsRead(notificationID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Notification marked as read"})
}

// MarkAllAsRead 标记所有通知为已读
// @Summary 标记所有通知为已读
// @Tags Notifications
// @Accept json
// @Produce json
// @Success 200
// @Router /notifications/read-all [post]
func (h *Handler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetInt("user_id")

	err := h.service.MarkAllAsRead(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "All notifications marked as read"})
}

// DeleteNotification 删除通知
// @Summary 删除通知
// @Tags Notifications
// @Accept json
// @Produce json
// @Param notification_id path int true "通知ID"
// @Success 200
// @Router /notifications/{notification_id} [delete]
func (h *Handler) DeleteNotification(c *gin.Context) {
	userID := c.GetInt("user_id")
	notificationID, _ := strconv.Atoi(c.Param("notification_id"))

	err := h.service.DeleteNotification(notificationID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Notification deleted"})
}
