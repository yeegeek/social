package payment

import "time"

// Transaction 交易记录模型
type Transaction struct {
	ID            int       `gorm:"primaryKey" json:"id"`
	UserID        int       `gorm:"not null" json:"user_id"`
	Type          string    `gorm:"not null" json:"type"` // recharge, vip_upgrade, purchase, withdraw, gift
	Amount        float64   `gorm:"not null" json:"amount"`
	BalanceBefore float64   `gorm:"not null" json:"balance_before"`
	BalanceAfter  float64   `gorm:"not null" json:"balance_after"`
	Description   string    `json:"description,omitempty"`
	Status        string    `gorm:"default:pending" json:"status"` // pending, completed, failed, cancelled
	PaymentMethod string    `json:"payment_method,omitempty"`      // paypal, wechat, alipay, stripe
	PaymentID     string    `json:"payment_id,omitempty"`
	RelatedID     *int      `json:"related_id,omitempty"`
	RelatedType   string    `json:"related_type,omitempty"` // post, vip, gift
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Transaction) TableName() string {
	return "transactions"
}
