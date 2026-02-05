package payment

import "gorm.io/gorm"

// Repository 支付仓库接口
type Repository interface {
	CreateTransaction(transaction *Transaction) error
	GetTransactionByID(id int) (*Transaction, error)
	GetUserTransactions(userID int, limit, offset int) ([]Transaction, int64, error)
	UpdateTransaction(transaction *Transaction) error
	GetUserBalance(userID int) (float64, error)
	UpdateUserBalance(userID int, amount float64) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository 创建支付仓库
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreateTransaction 创建交易记录
func (r *repository) CreateTransaction(transaction *Transaction) error {
	return r.db.Create(transaction).Error
}

// GetTransactionByID 根据ID获取交易记录
func (r *repository) GetTransactionByID(id int) (*Transaction, error) {
	var transaction Transaction
	err := r.db.First(&transaction, id).Error
	return &transaction, err
}

// GetUserTransactions 获取用户的交易记录列表
func (r *repository) GetUserTransactions(userID int, limit, offset int) ([]Transaction, int64, error) {
	var transactions []Transaction
	var total int64

	query := r.db.Model(&Transaction{}).Where("user_id = ?", userID).Order("created_at DESC")

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Limit(limit).Offset(offset).Find(&transactions).Error
	return transactions, total, err
}

// UpdateTransaction 更新交易记录
func (r *repository) UpdateTransaction(transaction *Transaction) error {
	return r.db.Save(transaction).Error
}

// GetUserBalance 获取用户余额（从 users 表的 coins 字段）
func (r *repository) GetUserBalance(userID int) (float64, error) {
	var balance float64
	err := r.db.Table("users").Select("coins").Where("id = ?", userID).Scan(&balance).Error
	return balance, err
}

// UpdateUserBalance 更新用户余额
func (r *repository) UpdateUserBalance(userID int, amount float64) error {
	return r.db.Table("users").Where("id = ?", userID).
		UpdateColumn("coins", gorm.Expr("coins + ?", amount)).Error
}
