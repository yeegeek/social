// Package friend 定义好友关系数据模型
package friend

import (
	"time"

	"github.com/yeegeek/uyou-go-api-starter/internal/user"
)

// FriendshipStatus 好友关系状态
type FriendshipStatus string

const (
	StatusPending  FriendshipStatus = "pending"  // 待接受
	StatusAccepted FriendshipStatus = "accepted" // 已接受
	StatusRejected FriendshipStatus = "rejected" // 已拒绝
	StatusBlocked  FriendshipStatus = "blocked"  // 已拉黑
)

// Friendship 好友关系模型
type Friendship struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	UserID      uint             `gorm:"not null;index:idx_friendship_user" json:"user_id"`
	FriendID    uint             `gorm:"not null;index:idx_friendship_friend" json:"friend_id"`
	Status      FriendshipStatus `gorm:"not null;index:idx_friendship_status" json:"status"`
	GroupName   string           `json:"group_name,omitempty"`
	Remark      string           `json:"remark,omitempty"`
	RequestedAt time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"requested_at"`
	AcceptedAt  *time.Time       `json:"accepted_at,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`

	// Relations
	User   user.User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Friend user.User `gorm:"foreignKey:FriendID" json:"friend,omitempty"`
}

// TableName 指定表名
func (Friendship) TableName() string {
	return "friendships"
}

// Blacklist 黑名单模型
type Blacklist struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"not null;index:idx_blacklist_user" json:"user_id"`
	BlockedUserID uint      `gorm:"not null;index:idx_blacklist_blocked" json:"blocked_user_id"`
	Reason        string    `json:"reason,omitempty"`
	CreatedAt     time.Time `json:"created_at"`

	// Relations
	User        user.User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BlockedUser user.User `gorm:"foreignKey:BlockedUserID" json:"blocked_user,omitempty"`
}

// TableName 指定表名
func (Blacklist) TableName() string {
	return "blacklist"
}
