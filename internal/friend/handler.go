package friend

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yeegeek/uyou-go-api-starter/internal/contextutil"
	"github.com/yeegeek/uyou-go-api-starter/internal/errors"
)

// Handler 好友处理器
type Handler struct {
	service Service
}

// NewHandler 创建好友处理器实例
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// SendFriendRequest godoc
// @Summary 发送好友请求
// @Description 向指定用户发送好友请求
// @Tags friends
// @Accept json
// @Produce json
// @Param request body SendFriendRequestRequest true "好友请求"
// @Success 200 {object} FriendshipResponse
// @Failure 400 {object} errors.Response
// @Failure 401 {object} errors.Response
// @Security BearerAuth
// @Router /api/v1/friends/request [post]
func (h *Handler) SendFriendRequest(c *gin.Context) {
	userID := contextutil.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, errors.Unauthorized("user not authenticated"))
		return
	}

	var req SendFriendRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ValidationError(err.Error()))
		return
	}

	friendship, err := h.service.SendFriendRequest(c.Request.Context(), userID, req.FriendID)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.Status, apiErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.InternalServerError(err))
		}
		return
	}

	c.JSON(http.StatusOK, friendship)
}

// GetFriendsList godoc
// @Summary 获取好友列表
// @Description 获取当前用户的好友列表
// @Tags friends
// @Accept json
// @Produce json
// @Success 200 {array} FriendshipResponse
// @Failure 401 {object} errors.Response
// @Security BearerAuth
// @Router /api/v1/friends [get]
func (h *Handler) GetFriendsList(c *gin.Context) {
	userID := contextutil.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, errors.Unauthorized("user not authenticated"))
		return
	}

	friends, err := h.service.GetFriendsList(c.Request.Context(), userID)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.Status, apiErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.InternalServerError(err))
		}
		return
	}

	c.JSON(http.StatusOK, friends)
}

// GetFriendRequests godoc
// @Summary 获取好友请求列表
// @Description 获取当前用户收到的好友请求列表
// @Tags friends
// @Accept json
// @Produce json
// @Success 200 {array} FriendshipResponse
// @Failure 401 {object} errors.Response
// @Security BearerAuth
// @Router /api/v1/friends/requests [get]
func (h *Handler) GetFriendRequests(c *gin.Context) {
	userID := contextutil.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, errors.Unauthorized("user not authenticated"))
		return
	}

	requests, err := h.service.GetFriendRequests(c.Request.Context(), userID)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.Status, apiErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.InternalServerError(err))
		}
		return
	}

	c.JSON(http.StatusOK, requests)
}

// AcceptFriendRequest godoc
// @Summary 接受好友请求
// @Description 接受指定的好友请求
// @Tags friends
// @Accept json
// @Produce json
// @Param id path int true "好友请求ID"
// @Success 200 {object} FriendshipResponse
// @Failure 400 {object} errors.Response
// @Failure 401 {object} errors.Response
// @Security BearerAuth
// @Router /api/v1/friends/requests/{id}/accept [post]
func (h *Handler) AcceptFriendRequest(c *gin.Context) {
	userID := contextutil.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, errors.Unauthorized("user not authenticated"))
		return
	}

	friendshipID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest("invalid friendship ID"))
		return
	}

	friendship, err := h.service.AcceptFriendRequest(c.Request.Context(), userID, uint(friendshipID))
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.Status, apiErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.InternalServerError(err))
		}
		return
	}

	c.JSON(http.StatusOK, friendship)
}

// RejectFriendRequest godoc
// @Summary 拒绝好友请求
// @Description 拒绝指定的好友请求
// @Tags friends
// @Accept json
// @Produce json
// @Param id path int true "好友请求ID"
// @Success 204
// @Failure 400 {object} errors.Response
// @Failure 401 {object} errors.Response
// @Security BearerAuth
// @Router /api/v1/friends/requests/{id}/reject [post]
func (h *Handler) RejectFriendRequest(c *gin.Context) {
	userID := contextutil.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, errors.Unauthorized("user not authenticated"))
		return
	}

	friendshipID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest("invalid friendship ID"))
		return
	}

	if err := h.service.RejectFriendRequest(c.Request.Context(), userID, uint(friendshipID)); err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.Status, apiErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.InternalServerError(err))
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteFriend godoc
// @Summary 删除好友
// @Description 删除指定的好友关系
// @Tags friends
// @Accept json
// @Produce json
// @Param id path int true "好友ID"
// @Success 204
// @Failure 400 {object} errors.Response
// @Failure 401 {object} errors.Response
// @Security BearerAuth
// @Router /api/v1/friends/{id} [delete]
func (h *Handler) DeleteFriend(c *gin.Context) {
	userID := contextutil.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, errors.Unauthorized("user not authenticated"))
		return
	}

	friendID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest("invalid friend ID"))
		return
	}

	if err := h.service.DeleteFriend(c.Request.Context(), userID, uint(friendID)); err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.Status, apiErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.InternalServerError(err))
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateFriendRemark godoc
// @Summary 更新好友备注
// @Description 更新指定好友的备注信息
// @Tags friends
// @Accept json
// @Produce json
// @Param id path int true "好友ID"
// @Param request body UpdateFriendRemarkRequest true "备注信息"
// @Success 204
// @Failure 400 {object} errors.Response
// @Failure 401 {object} errors.Response
// @Security BearerAuth
// @Router /api/v1/friends/{id}/remark [put]
func (h *Handler) UpdateFriendRemark(c *gin.Context) {
	userID := contextutil.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, errors.Unauthorized("user not authenticated"))
		return
	}

	friendID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest("invalid friend ID"))
		return
	}

	var req UpdateFriendRemarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ValidationError(err.Error()))
		return
	}

	if err := h.service.UpdateFriendRemark(c.Request.Context(), userID, uint(friendID), req.Remark); err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.Status, apiErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.InternalServerError(err))
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateFriendGroup godoc
// @Summary 更新好友分组
// @Description 更新指定好友的分组信息
// @Tags friends
// @Accept json
// @Produce json
// @Param id path int true "好友ID"
// @Param request body UpdateFriendGroupRequest true "分组信息"
// @Success 204
// @Failure 400 {object} errors.Response
// @Failure 401 {object} errors.Response
// @Security BearerAuth
// @Router /api/v1/friends/{id}/group [put]
func (h *Handler) UpdateFriendGroup(c *gin.Context) {
	userID := contextutil.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, errors.Unauthorized("user not authenticated"))
		return
	}

	friendID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest("invalid friend ID"))
		return
	}

	var req UpdateFriendGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ValidationError(err.Error()))
		return
	}

	if err := h.service.UpdateFriendGroup(c.Request.Context(), userID, uint(friendID), req.GroupName); err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.Status, apiErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.InternalServerError(err))
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// BlockUser godoc
// @Summary 拉黑用户
// @Description 将指定用户加入黑名单
// @Tags friends
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param request body BlockUserRequest true "拉黑信息"
// @Success 204
// @Failure 400 {object} errors.Response
// @Failure 401 {object} errors.Response
// @Security BearerAuth
// @Router /api/v1/friends/{id}/block [post]
func (h *Handler) BlockUser(c *gin.Context) {
	userID := contextutil.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, errors.Unauthorized("user not authenticated"))
		return
	}

	blockedUserID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest("invalid user ID"))
		return
	}

	var req BlockUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ValidationError(err.Error()))
		return
	}

	if err := h.service.BlockUser(c.Request.Context(), userID, uint(blockedUserID), req.Reason); err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.Status, apiErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.InternalServerError(err))
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// UnblockUser godoc
// @Summary 取消拉黑
// @Description 将指定用户从黑名单移除
// @Tags friends
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 204
// @Failure 400 {object} errors.Response
// @Failure 401 {object} errors.Response
// @Security BearerAuth
// @Router /api/v1/friends/{id}/unblock [delete]
func (h *Handler) UnblockUser(c *gin.Context) {
	userID := contextutil.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, errors.Unauthorized("user not authenticated"))
		return
	}

	blockedUserID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest("invalid user ID"))
		return
	}

	if err := h.service.UnblockUser(c.Request.Context(), userID, uint(blockedUserID)); err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.Status, apiErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.InternalServerError(err))
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// GetBlockedUsers godoc
// @Summary 获取黑名单列表
// @Description 获取当前用户的黑名单列表
// @Tags friends
// @Accept json
// @Produce json
// @Success 200 {array} BlacklistResponse
// @Failure 401 {object} errors.Response
// @Security BearerAuth
// @Router /api/v1/friends/blocked [get]
func (h *Handler) GetBlockedUsers(c *gin.Context) {
	userID := contextutil.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, errors.Unauthorized("user not authenticated"))
		return
	}

	blocked, err := h.service.GetBlockedUsers(c.Request.Context(), userID)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.Status, apiErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.InternalServerError(err))
		}
		return
	}

	c.JSON(http.StatusOK, blocked)
}
