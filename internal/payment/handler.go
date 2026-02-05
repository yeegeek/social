package payment

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler 支付处理器
type Handler struct {
	service Service
}

// NewHandler 创建支付处理器
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Recharge 充值
// @Summary 充值
// @Tags Payment
// @Accept json
// @Produce json
// @Param request body RechargeRequest true "充值请求"
// @Success 200 {object} TransactionResponse
// @Router /payment/recharge [post]
func (h *Handler) Recharge(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req RechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Recharge(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// VIPUpgrade VIP升级
// @Summary VIP升级
// @Tags Payment
// @Accept json
// @Produce json
// @Param request body VIPUpgradeRequest true "VIP升级请求"
// @Success 200 {object} TransactionResponse
// @Router /payment/vip-upgrade [post]
func (h *Handler) VIPUpgrade(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req VIPUpgradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.VIPUpgrade(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// Purchase 购买
// @Summary 购买付费内容
// @Tags Payment
// @Accept json
// @Produce json
// @Param request body PurchaseRequest true "购买请求"
// @Success 200 {object} TransactionResponse
// @Router /payment/purchase [post]
func (h *Handler) Purchase(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req PurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Purchase(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// GetTransactions 获取交易记录列表
// @Summary 获取交易记录列表
// @Tags Payment
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} TransactionListResponse
// @Router /payment/transactions [get]
func (h *Handler) GetTransactions(c *gin.Context) {
	userID := c.GetInt("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	response, err := h.service.GetTransactions(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// GetBalance 获取余额
// @Summary 获取用户余额
// @Tags Payment
// @Accept json
// @Produce json
// @Success 200 {object} BalanceResponse
// @Router /payment/balance [get]
func (h *Handler) GetBalance(c *gin.Context) {
	userID := c.GetInt("user_id")

	response, err := h.service.GetBalance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}
