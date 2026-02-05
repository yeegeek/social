package payment

import (
	"errors"
	"fmt"
	"time"
)

// Service 支付服务接口
type Service interface {
	Recharge(userID int, req *RechargeRequest) (*TransactionResponse, error)
	VIPUpgrade(userID int, req *VIPUpgradeRequest) (*TransactionResponse, error)
	Purchase(userID int, req *PurchaseRequest) (*TransactionResponse, error)
	GetTransactions(userID int, page, pageSize int) (*TransactionListResponse, error)
	GetBalance(userID int) (*BalanceResponse, error)
}

type service struct {
	repo Repository
}

// NewService 创建支付服务
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Recharge 充值
func (s *service) Recharge(userID int, req *RechargeRequest) (*TransactionResponse, error) {
	// 获取当前余额
	balanceBefore, err := s.repo.GetUserBalance(userID)
	if err != nil {
		return nil, err
	}

	// 创建交易记录
	transaction := &Transaction{
		UserID:        userID,
		Type:          "recharge",
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceBefore + req.Amount,
		Description:   fmt.Sprintf("Recharge %.2f via %s", req.Amount, req.PaymentMethod),
		Status:        "pending",
		PaymentMethod: req.PaymentMethod,
	}

	err = s.repo.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	// TODO: 集成真实支付网关（PayPal, Stripe等）
	// 这里模拟支付成功
	transaction.Status = "completed"
	transaction.PaymentID = fmt.Sprintf("PAY_%d_%d", userID, time.Now().Unix())
	err = s.repo.UpdateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	// 更新用户余额
	err = s.repo.UpdateUserBalance(userID, req.Amount)
	if err != nil {
		return nil, err
	}

	return s.toTransactionResponse(transaction), nil
}

// VIPUpgrade VIP升级
func (s *service) VIPUpgrade(userID int, req *VIPUpgradeRequest) (*TransactionResponse, error) {
	// VIP价格表
	priceMap := map[int]float64{
		30:  9.99,
		90:  24.99,
		180: 44.99,
		365: 79.99,
	}

	price, ok := priceMap[req.Duration]
	if !ok {
		return nil, errors.New("invalid VIP duration")
	}

	// 获取当前余额
	balanceBefore, err := s.repo.GetUserBalance(userID)
	if err != nil {
		return nil, err
	}

	// 检查余额是否足够
	if balanceBefore < price {
		return nil, errors.New("insufficient balance")
	}

	// 创建交易记录
	transaction := &Transaction{
		UserID:        userID,
		Type:          "vip_upgrade",
		Amount:        -price,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceBefore - price,
		Description:   fmt.Sprintf("VIP upgrade for %d days", req.Duration),
		Status:        "completed",
		RelatedType:   "vip",
	}

	err = s.repo.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	// 更新用户余额
	err = s.repo.UpdateUserBalance(userID, -price)
	if err != nil {
		return nil, err
	}

	// TODO: 更新用户VIP状态和过期时间

	return s.toTransactionResponse(transaction), nil
}

// Purchase 购买（付费内容）
func (s *service) Purchase(userID int, req *PurchaseRequest) (*TransactionResponse, error) {
	// 获取当前余额
	balanceBefore, err := s.repo.GetUserBalance(userID)
	if err != nil {
		return nil, err
	}

	// TODO: 根据 RelatedType 和 RelatedID 获取实际价格
	price := 1.99 // 示例价格

	// 检查余额是否足够
	if balanceBefore < price {
		return nil, errors.New("insufficient balance")
	}

	// 创建交易记录
	transaction := &Transaction{
		UserID:        userID,
		Type:          "purchase",
		Amount:        -price,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceBefore - price,
		Description:   fmt.Sprintf("Purchase %s #%d", req.RelatedType, req.RelatedID),
		Status:        "completed",
		RelatedID:     &req.RelatedID,
		RelatedType:   req.RelatedType,
	}

	err = s.repo.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	// 更新用户余额
	err = s.repo.UpdateUserBalance(userID, -price)
	if err != nil {
		return nil, err
	}

	return s.toTransactionResponse(transaction), nil
}

// GetTransactions 获取交易记录列表
func (s *service) GetTransactions(userID int, page, pageSize int) (*TransactionListResponse, error) {
	offset := (page - 1) * pageSize
	transactions, total, err := s.repo.GetUserTransactions(userID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	response := &TransactionListResponse{
		Transactions: make([]TransactionResponse, 0, len(transactions)),
		Total:        total,
	}

	for _, trans := range transactions {
		response.Transactions = append(response.Transactions, *s.toTransactionResponse(&trans))
	}

	return response, nil
}

// GetBalance 获取用户余额
func (s *service) GetBalance(userID int) (*BalanceResponse, error) {
	balance, err := s.repo.GetUserBalance(userID)
	if err != nil {
		return nil, err
	}

	return &BalanceResponse{
		Balance: balance,
		Coins:   int(balance),
	}, nil
}

// toTransactionResponse 转换为交易响应
func (s *service) toTransactionResponse(transaction *Transaction) *TransactionResponse {
	return &TransactionResponse{
		ID:            transaction.ID,
		UserID:        transaction.UserID,
		Type:          transaction.Type,
		Amount:        transaction.Amount,
		BalanceBefore: transaction.BalanceBefore,
		BalanceAfter:  transaction.BalanceAfter,
		Description:   transaction.Description,
		Status:        transaction.Status,
		PaymentMethod: transaction.PaymentMethod,
		PaymentID:     transaction.PaymentID,
		RelatedID:     transaction.RelatedID,
		RelatedType:   transaction.RelatedType,
		CreatedAt:     transaction.CreatedAt,
	}
}
