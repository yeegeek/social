package friend

import (
	"context"
	"fmt"
	"time"

	"github.com/yeegeek/uyou-go-api-starter/internal/errors"
)

// Service 好友服务接口
type Service interface {
	// Friend operations
	SendFriendRequest(ctx context.Context, userID, friendID uint) (*FriendshipResponse, error)
	AcceptFriendRequest(ctx context.Context, userID, friendshipID uint) (*FriendshipResponse, error)
	RejectFriendRequest(ctx context.Context, userID, friendshipID uint) error
	DeleteFriend(ctx context.Context, userID, friendID uint) error
	GetFriendsList(ctx context.Context, userID uint) ([]*FriendshipResponse, error)
	GetFriendRequests(ctx context.Context, userID uint) ([]*FriendshipResponse, error)
	UpdateFriendRemark(ctx context.Context, userID, friendID uint, remark string) error
	UpdateFriendGroup(ctx context.Context, userID, friendID uint, groupName string) error

	// Blacklist operations
	BlockUser(ctx context.Context, userID, blockedUserID uint, reason string) error
	UnblockUser(ctx context.Context, userID, blockedUserID uint) error
	GetBlockedUsers(ctx context.Context, userID uint) ([]*BlacklistResponse, error)
	IsBlocked(ctx context.Context, userID, blockedUserID uint) (bool, error)
}

// service 好友服务实现
type service struct {
	repo Repository
}

// NewService 创建好友服务实例
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// SendFriendRequest 发送好友请求
func (s *service) SendFriendRequest(ctx context.Context, userID, friendID uint) (*FriendshipResponse, error) {
	// 检查是否给自己发送请求
	if userID == friendID {
		return nil, errors.BadRequest("cannot send friend request to yourself")
	}

	// 检查是否已经是好友或已有请求
	existing, err := s.repo.GetFriendship(ctx, userID, friendID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing friendship: %w", err)
	}
	if existing != nil {
		if existing.Status == StatusAccepted {
			return nil, errors.Conflict("already friends")
		}
		if existing.Status == StatusPending {
			return nil, errors.Conflict("friend request already sent")
		}
	}

	// 检查是否被拉黑
	blocked, err := s.repo.IsBlocked(ctx, friendID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check blacklist: %w", err)
	}
	if blocked {
		return nil, errors.Forbidden("cannot send friend request to this user")
	}

	// 创建好友请求
	friendship := &Friendship{
		UserID:      userID,
		FriendID:    friendID,
		Status:      StatusPending,
		RequestedAt: time.Now(),
	}

	if err := s.repo.CreateFriendship(ctx, friendship); err != nil {
		return nil, fmt.Errorf("failed to create friendship: %w", err)
	}

	// 重新加载以获取关联数据
	friendship, err = s.repo.GetFriendshipByID(ctx, friendship.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get friendship: %w", err)
	}

	return ToFriendshipResponse(friendship), nil
}

// AcceptFriendRequest 接受好友请求
func (s *service) AcceptFriendRequest(ctx context.Context, userID, friendshipID uint) (*FriendshipResponse, error) {
	friendship, err := s.repo.GetFriendshipByID(ctx, friendshipID)
	if err != nil {
		return nil, fmt.Errorf("failed to get friendship: %w", err)
	}
	if friendship == nil {
		return nil, errors.NotFound("friendship not found")
	}

	// 验证是否为接收方
	if friendship.FriendID != userID {
		return nil, errors.Forbidden("not authorized to accept this request")
	}

	// 验证状态
	if friendship.Status != StatusPending {
		return nil, errors.BadRequest("friend request is not pending")
	}

	// 更新状态
	now := time.Now()
	friendship.Status = StatusAccepted
	friendship.AcceptedAt = &now

	if err := s.repo.UpdateFriendship(ctx, friendship); err != nil {
		return nil, fmt.Errorf("failed to update friendship: %w", err)
	}

	return ToFriendshipResponse(friendship), nil
}

// RejectFriendRequest 拒绝好友请求
func (s *service) RejectFriendRequest(ctx context.Context, userID, friendshipID uint) error {
	friendship, err := s.repo.GetFriendshipByID(ctx, friendshipID)
	if err != nil {
		return fmt.Errorf("failed to get friendship: %w", err)
	}
	if friendship == nil {
		return errors.NotFound("friendship not found")
	}

	// 验证是否为接收方
	if friendship.FriendID != userID {
		return errors.Forbidden("not authorized to reject this request")
	}

	// 验证状态
	if friendship.Status != StatusPending {
		return errors.BadRequest("friend request is not pending")
	}

	// 删除请求
	return s.repo.DeleteFriendship(ctx, friendshipID)
}

// DeleteFriend 删除好友
func (s *service) DeleteFriend(ctx context.Context, userID, friendID uint) error {
	friendship, err := s.repo.GetFriendship(ctx, userID, friendID)
	if err != nil {
		return fmt.Errorf("failed to get friendship: %w", err)
	}
	if friendship == nil {
		return errors.NotFound("friendship not found")
	}

	// 验证是否为关系中的一方
	if friendship.UserID != userID && friendship.FriendID != userID {
		return errors.Forbidden("not authorized to delete this friendship")
	}

	return s.repo.DeleteFriendship(ctx, friendship.ID)
}

// GetFriendsList 获取好友列表
func (s *service) GetFriendsList(ctx context.Context, userID uint) ([]*FriendshipResponse, error) {
	friendships, err := s.repo.ListFriends(ctx, userID, StatusAccepted)
	if err != nil {
		return nil, fmt.Errorf("failed to list friends: %w", err)
	}

	responses := make([]*FriendshipResponse, len(friendships))
	for i, f := range friendships {
		responses[i] = ToFriendshipResponse(f)
	}

	return responses, nil
}

// GetFriendRequests 获取好友请求列表
func (s *service) GetFriendRequests(ctx context.Context, userID uint) ([]*FriendshipResponse, error) {
	friendships, err := s.repo.ListFriendRequests(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list friend requests: %w", err)
	}

	responses := make([]*FriendshipResponse, len(friendships))
	for i, f := range friendships {
		responses[i] = ToFriendshipResponse(f)
	}

	return responses, nil
}

// UpdateFriendRemark 更新好友备注
func (s *service) UpdateFriendRemark(ctx context.Context, userID, friendID uint, remark string) error {
	friendship, err := s.repo.GetFriendship(ctx, userID, friendID)
	if err != nil {
		return fmt.Errorf("failed to get friendship: %w", err)
	}
	if friendship == nil {
		return errors.NotFound("friendship not found")
	}

	// 验证是否为好友
	if friendship.Status != StatusAccepted {
		return errors.BadRequest("not friends")
	}

	// 确保更新的是当前用户的备注
	if friendship.UserID != userID {
		// 如果当前用户是 FriendID，需要交换关系
		friendship.UserID = userID
		friendship.FriendID = friendID
	}

	friendship.Remark = remark
	return s.repo.UpdateFriendship(ctx, friendship)
}

// UpdateFriendGroup 更新好友分组
func (s *service) UpdateFriendGroup(ctx context.Context, userID, friendID uint, groupName string) error {
	friendship, err := s.repo.GetFriendship(ctx, userID, friendID)
	if err != nil {
		return fmt.Errorf("failed to get friendship: %w", err)
	}
	if friendship == nil {
		return errors.NotFound("friendship not found")
	}

	// 验证是否为好友
	if friendship.Status != StatusAccepted {
		return errors.BadRequest("not friends")
	}

	// 确保更新的是当前用户的分组
	if friendship.UserID != userID {
		friendship.UserID = userID
		friendship.FriendID = friendID
	}

	friendship.GroupName = groupName
	return s.repo.UpdateFriendship(ctx, friendship)
}

// BlockUser 拉黑用户
func (s *service) BlockUser(ctx context.Context, userID, blockedUserID uint, reason string) error {
	// 检查是否给自己拉黑
	if userID == blockedUserID {
		return errors.BadRequest("cannot block yourself")
	}

	// 检查是否已拉黑
	existing, err := s.repo.GetBlacklist(ctx, userID, blockedUserID)
	if err != nil {
		return fmt.Errorf("failed to check blacklist: %w", err)
	}
	if existing != nil {
		return errors.Conflict("user already blocked")
	}

	// 删除好友关系（如果存在）
	friendship, err := s.repo.GetFriendship(ctx, userID, blockedUserID)
	if err != nil {
		return fmt.Errorf("failed to check friendship: %w", err)
	}
	if friendship != nil {
		if err := s.repo.DeleteFriendship(ctx, friendship.ID); err != nil {
			return fmt.Errorf("failed to delete friendship: %w", err)
		}
	}

	// 创建黑名单记录
	blacklist := &Blacklist{
		UserID:        userID,
		BlockedUserID: blockedUserID,
		Reason:        reason,
	}

	return s.repo.CreateBlacklist(ctx, blacklist)
}

// UnblockUser 取消拉黑
func (s *service) UnblockUser(ctx context.Context, userID, blockedUserID uint) error {
	blacklist, err := s.repo.GetBlacklist(ctx, userID, blockedUserID)
	if err != nil {
		return fmt.Errorf("failed to get blacklist: %w", err)
	}
	if blacklist == nil {
		return errors.NotFound("blacklist entry not found")
	}

	return s.repo.DeleteBlacklist(ctx, blacklist.ID)
}

// GetBlockedUsers 获取黑名单列表
func (s *service) GetBlockedUsers(ctx context.Context, userID uint) ([]*BlacklistResponse, error) {
	blacklist, err := s.repo.ListBlacklist(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list blacklist: %w", err)
	}

	responses := make([]*BlacklistResponse, len(blacklist))
	for i, b := range blacklist {
		responses[i] = ToBlacklistResponse(b)
	}

	return responses, nil
}

// IsBlocked 检查是否被拉黑
func (s *service) IsBlocked(ctx context.Context, userID, blockedUserID uint) (bool, error) {
	return s.repo.IsBlocked(ctx, userID, blockedUserID)
}
