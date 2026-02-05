package payment

import "time"

// RechargeRequest 充值请求
type RechargeRequest struct {
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	PaymentMethod string  `json:"payment_method" validate:"required,oneof=paypal wechat alipay stripe"`
}

// VIPUpgradeRequest VIP升级请求
type VIPUpgradeRequest struct {
	Duration      int    `json:"duration" validate:"required,oneof=30 90 180 365"` // 天数
	PaymentMethod string `json:"payment_method" validate:"required"`
}

// PurchaseRequest 购买请求
type PurchaseRequest struct {
	RelatedID   int    `json:"related_id" validate:"required"`
	RelatedType string `json:"related_type" validate:"required,oneof=post gift"`
}

// TransactionResponse 交易响应
type TransactionResponse struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	Description   string    `json:"description,omitempty"`
	Status        string    `json:"status"`
	PaymentMethod string    `json:"payment_method,omitempty"`
	PaymentID     string    `json:"payment_id,omitempty"`
	RelatedID     *int      `json:"related_id,omitempty"`
	RelatedType   string    `json:"related_type,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// TransactionListResponse 交易列表响应
type TransactionListResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Total        int64                 `json:"total"`
}

// BalanceResponse 余额响应
type BalanceResponse struct {
	Balance float64 `json:"balance"`
	Coins   int     `json:"coins"`
}
