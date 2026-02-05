package message

import (
	"errors"
	"time"
)

// Service 消息服务接口
type Service interface {
	// Conversation operations
	GetConversations(userID int, page, pageSize int) (*ConversationListResponse, error)
	GetConversation(conversationID, userID int) (*ConversationResponse, error)
	DeleteConversation(conversationID, userID int) error

	// Message operations
	SendMessage(senderID int, req *SendMessageRequest) (*MessageResponse, error)
	GetMessages(conversationID, userID int, page, pageSize int) (*MessageListResponse, error)
	MarkAsRead(messageID, userID int) error
	MarkConversationAsRead(conversationID, userID int) error
	RecallMessage(messageID, userID int) error

	// Draft operations
	SaveDraft(userID int, req *SaveDraftRequest) (*MessageResponse, error)
	GetDrafts(userID int) ([]MessageResponse, error)
	DeleteDraft(draftID, userID int) error

	// Translation
	TranslateMessage(messageID, userID int, targetLanguage string) (*MessageResponse, error)
}

type service struct {
	repo Repository
}

// NewService 创建消息服务
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// GetConversations 获取对话列表
func (s *service) GetConversations(userID int, page, pageSize int) (*ConversationListResponse, error) {
	offset := (page - 1) * pageSize
	conversations, total, err := s.repo.GetUserConversations(userID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	response := &ConversationListResponse{
		Conversations: make([]ConversationResponse, 0, len(conversations)),
		Total:         total,
	}

	for _, conv := range conversations {
		// 确定对方用户ID
		otherUserID := conv.User2ID
		unreadCount := conv.User2UnreadCount
		if conv.User2ID == userID {
			otherUserID = conv.User1ID
			unreadCount = conv.User1UnreadCount
		}

		convResp := ConversationResponse{
			ID:          conv.ID,
			OtherUserID: otherUserID,
			UnreadCount: unreadCount,
			UpdatedAt:   conv.UpdatedAt,
		}

		// 获取最后一条消息
		if conv.LastMessageID != nil {
			lastMsg, err := s.repo.GetMessageByID(*conv.LastMessageID)
			if err == nil {
				convResp.LastMessage = s.toMessageResponse(lastMsg)
			}
		}

		response.Conversations = append(response.Conversations, convResp)
	}

	return response, nil
}

// GetConversation 获取单个对话
func (s *service) GetConversation(conversationID, userID int) (*ConversationResponse, error) {
	conv, err := s.repo.GetConversationByID(conversationID)
	if err != nil {
		return nil, err
	}

	// 验证用户是否属于该对话
	if conv.User1ID != userID && conv.User2ID != userID {
		return nil, errors.New("unauthorized access to conversation")
	}

	otherUserID := conv.User2ID
	unreadCount := conv.User2UnreadCount
	if conv.User2ID == userID {
		otherUserID = conv.User1ID
		unreadCount = conv.User1UnreadCount
	}

	response := &ConversationResponse{
		ID:          conv.ID,
		OtherUserID: otherUserID,
		UnreadCount: unreadCount,
		UpdatedAt:   conv.UpdatedAt,
	}

	return response, nil
}

// DeleteConversation 删除对话
func (s *service) DeleteConversation(conversationID, userID int) error {
	return s.repo.DeleteConversation(conversationID, userID)
}

// SendMessage 发送消息
func (s *service) SendMessage(senderID int, req *SendMessageRequest) (*MessageResponse, error) {
	// 获取或创建对话
	conversation, err := s.repo.GetOrCreateConversation(senderID, req.ReceiverID)
	if err != nil {
		return nil, err
	}

	// 创建消息
	messageType := req.MessageType
	if messageType == "" {
		messageType = "text"
	}

	message := &Message{
		ConversationID: conversation.ID,
		SenderID:       senderID,
		ReceiverID:     req.ReceiverID,
		Content:        req.Content,
		MessageType:    messageType,
		MediaURL:       req.MediaURL,
		GiftID:         req.GiftID,
		IsDraft:        false,
	}

	err = s.repo.CreateMessage(message)
	if err != nil {
		return nil, err
	}

	// 更新对话信息
	conversation.LastMessageID = &message.ID
	now := time.Now()
	conversation.LastMessageAt = &now

	// 更新未读计数
	if conversation.User1ID == req.ReceiverID {
		conversation.User1UnreadCount++
	} else {
		conversation.User2UnreadCount++
	}

	err = s.repo.UpdateConversation(conversation)
	if err != nil {
		return nil, err
	}

	return s.toMessageResponse(message), nil
}

// GetMessages 获取消息列表
func (s *service) GetMessages(conversationID, userID int, page, pageSize int) (*MessageListResponse, error) {
	// 验证用户是否属于该对话
	conv, err := s.repo.GetConversationByID(conversationID)
	if err != nil {
		return nil, err
	}

	if conv.User1ID != userID && conv.User2ID != userID {
		return nil, errors.New("unauthorized access to conversation")
	}

	offset := (page - 1) * pageSize
	messages, total, err := s.repo.GetConversationMessages(conversationID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	response := &MessageListResponse{
		Messages: make([]MessageResponse, 0, len(messages)),
		Total:    total,
	}

	for _, msg := range messages {
		response.Messages = append(response.Messages, *s.toMessageResponse(&msg))
	}

	return response, nil
}

// MarkAsRead 标记消息为已读
func (s *service) MarkAsRead(messageID, userID int) error {
	return s.repo.MarkAsRead(messageID, userID)
}

// MarkConversationAsRead 标记对话为已读
func (s *service) MarkConversationAsRead(conversationID, userID int) error {
	// 标记所有消息为已读
	err := s.repo.MarkConversationAsRead(conversationID, userID)
	if err != nil {
		return err
	}

	// 重置未读计数
	conv, err := s.repo.GetConversationByID(conversationID)
	if err != nil {
		return err
	}

	if conv.User1ID == userID {
		conv.User1UnreadCount = 0
	} else if conv.User2ID == userID {
		conv.User2UnreadCount = 0
	}

	return s.repo.UpdateConversation(conv)
}

// RecallMessage 撤回消息
func (s *service) RecallMessage(messageID, userID int) error {
	message, err := s.repo.GetMessageByID(messageID)
	if err != nil {
		return err
	}

	// 验证是否是发送者
	if message.SenderID != userID {
		return errors.New("only sender can recall message")
	}

	// 检查是否在可撤回时间内（例如2分钟）
	if time.Since(message.CreatedAt) > 2*time.Minute {
		return errors.New("message can only be recalled within 2 minutes")
	}

	return s.repo.RecallMessage(messageID)
}

// SaveDraft 保存草稿
func (s *service) SaveDraft(userID int, req *SaveDraftRequest) (*MessageResponse, error) {
	// 获取或创建对话
	conversation, err := s.repo.GetOrCreateConversation(userID, req.ReceiverID)
	if err != nil {
		return nil, err
	}

	// 创建草稿
	draft := &Message{
		ConversationID: conversation.ID,
		SenderID:       userID,
		ReceiverID:     req.ReceiverID,
		Content:        req.Content,
		MessageType:    "text",
		IsDraft:        true,
	}

	err = s.repo.CreateMessage(draft)
	if err != nil {
		return nil, err
	}

	return s.toMessageResponse(draft), nil
}

// GetDrafts 获取草稿列表
func (s *service) GetDrafts(userID int) ([]MessageResponse, error) {
	drafts, err := s.repo.GetDrafts(userID)
	if err != nil {
		return nil, err
	}

	response := make([]MessageResponse, 0, len(drafts))
	for _, draft := range drafts {
		response = append(response, *s.toMessageResponse(&draft))
	}

	return response, nil
}

// DeleteDraft 删除草稿
func (s *service) DeleteDraft(draftID, userID int) error {
	draft, err := s.repo.GetMessageByID(draftID)
	if err != nil {
		return err
	}

	// 验证是否是草稿所有者
	if draft.SenderID != userID {
		return errors.New("unauthorized to delete draft")
	}

	if !draft.IsDraft {
		return errors.New("not a draft message")
	}

	return s.repo.DeleteDraft(draftID)
}

// TranslateMessage 翻译消息
func (s *service) TranslateMessage(messageID, userID int, targetLanguage string) (*MessageResponse, error) {
	message, err := s.repo.GetMessageByID(messageID)
	if err != nil {
		return nil, err
	}

	// 验证用户是否有权限访问该消息
	if message.SenderID != userID && message.ReceiverID != userID {
		return nil, errors.New("unauthorized access to message")
	}

	// TODO: 集成翻译服务（Google Translate API）
	// 这里暂时返回原文
	message.TranslatedContent = message.Content + " [translated to " + targetLanguage + "]"

	err = s.repo.UpdateMessage(message)
	if err != nil {
		return nil, err
	}

	return s.toMessageResponse(message), nil
}

// toMessageResponse 转换为消息响应
func (s *service) toMessageResponse(message *Message) *MessageResponse {
	return &MessageResponse{
		ID:                message.ID,
		ConversationID:    message.ConversationID,
		SenderID:          message.SenderID,
		ReceiverID:        message.ReceiverID,
		Content:           message.Content,
		MessageType:       message.MessageType,
		MediaURL:          message.MediaURL,
		TranslatedContent: message.TranslatedContent,
		IsRead:            message.IsRead,
		ReadAt:            message.ReadAt,
		IsRecalled:        message.IsRecalled,
		IsDraft:           message.IsDraft,
		CreatedAt:         message.CreatedAt,
	}
}
