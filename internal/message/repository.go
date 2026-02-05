package message

import (
	"time"

	"gorm.io/gorm"
)

// Repository 消息仓库接口
type Repository interface {
	// Conversation operations
	GetOrCreateConversation(user1ID, user2ID int) (*Conversation, error)
	GetConversationByID(id int) (*Conversation, error)
	GetUserConversations(userID int, limit, offset int) ([]Conversation, int64, error)
	UpdateConversation(conversation *Conversation) error
	DeleteConversation(conversationID, userID int) error
	
	// Message operations
	CreateMessage(message *Message) error
	GetMessageByID(id int) (*Message, error)
	GetConversationMessages(conversationID int, limit, offset int) ([]Message, int64, error)
	UpdateMessage(message *Message) error
	MarkAsRead(messageID, userID int) error
	MarkConversationAsRead(conversationID, userID int) error
	RecallMessage(messageID int) error
	GetDrafts(userID int) ([]Message, error)
	DeleteDraft(messageID int) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository 创建消息仓库
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// GetOrCreateConversation 获取或创建对话
func (r *repository) GetOrCreateConversation(user1ID, user2ID int) (*Conversation, error) {
	// 确保 user1ID < user2ID
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	var conversation Conversation
	err := r.db.Where("user1_id = ? AND user2_id = ?", user1ID, user2ID).First(&conversation).Error
	if err == gorm.ErrRecordNotFound {
		conversation = Conversation{
			User1ID: user1ID,
			User2ID: user2ID,
		}
		err = r.db.Create(&conversation).Error
	}
	return &conversation, err
}

// GetConversationByID 根据ID获取对话
func (r *repository) GetConversationByID(id int) (*Conversation, error) {
	var conversation Conversation
	err := r.db.First(&conversation, id).Error
	return &conversation, err
}

// GetUserConversations 获取用户的对话列表
func (r *repository) GetUserConversations(userID int, limit, offset int) ([]Conversation, int64, error) {
	var conversations []Conversation
	var total int64

	query := r.db.Model(&Conversation{}).
		Where("(user1_id = ? AND user1_deleted = false) OR (user2_id = ? AND user2_deleted = false)", userID, userID).
		Order("last_message_at DESC NULLS LAST")

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Limit(limit).Offset(offset).Find(&conversations).Error
	return conversations, total, err
}

// UpdateConversation 更新对话
func (r *repository) UpdateConversation(conversation *Conversation) error {
	return r.db.Save(conversation).Error
}

// DeleteConversation 删除对话（软删除）
func (r *repository) DeleteConversation(conversationID, userID int) error {
	conversation, err := r.GetConversationByID(conversationID)
	if err != nil {
		return err
	}

	if conversation.User1ID == userID {
		conversation.User1Deleted = true
	} else if conversation.User2ID == userID {
		conversation.User2Deleted = true
	}

	return r.UpdateConversation(conversation)
}

// CreateMessage 创建消息
func (r *repository) CreateMessage(message *Message) error {
	return r.db.Create(message).Error
}

// GetMessageByID 根据ID获取消息
func (r *repository) GetMessageByID(id int) (*Message, error) {
	var message Message
	err := r.db.First(&message, id).Error
	return &message, err
}

// GetConversationMessages 获取对话的消息列表
func (r *repository) GetConversationMessages(conversationID int, limit, offset int) ([]Message, int64, error) {
	var messages []Message
	var total int64

	query := r.db.Model(&Message{}).
		Where("conversation_id = ? AND is_draft = false", conversationID).
		Order("created_at DESC")

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Limit(limit).Offset(offset).Find(&messages).Error
	return messages, total, err
}

// UpdateMessage 更新消息
func (r *repository) UpdateMessage(message *Message) error {
	return r.db.Save(message).Error
}

// MarkAsRead 标记消息为已读
func (r *repository) MarkAsRead(messageID, userID int) error {
	now := time.Now()
	return r.db.Model(&Message{}).
		Where("id = ? AND receiver_id = ?", messageID, userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

// MarkConversationAsRead 标记对话中所有消息为已读
func (r *repository) MarkConversationAsRead(conversationID, userID int) error {
	now := time.Now()
	return r.db.Model(&Message{}).
		Where("conversation_id = ? AND receiver_id = ? AND is_read = false", conversationID, userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

// RecallMessage 撤回消息
func (r *repository) RecallMessage(messageID int) error {
	now := time.Now()
	return r.db.Model(&Message{}).
		Where("id = ?", messageID).
		Updates(map[string]interface{}{
			"is_recalled": true,
			"recalled_at": now,
		}).Error
}

// GetDrafts 获取用户的草稿
func (r *repository) GetDrafts(userID int) ([]Message, error) {
	var drafts []Message
	err := r.db.Where("sender_id = ? AND is_draft = true", userID).
		Order("created_at DESC").
		Find(&drafts).Error
	return drafts, err
}

// DeleteDraft 删除草稿
func (r *repository) DeleteDraft(messageID int) error {
	return r.db.Delete(&Message{}, messageID).Error
}
