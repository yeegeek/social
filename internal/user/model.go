// Package user 定义用户数据模型
package user

import (
	"time"

	"gorm.io/gorm"
)

// User 表示系统中的用户实体
type User struct {
	ID             uint           `gorm:"primaryKey" json:"id"`                      // 用户唯一标识
	Name           string         `gorm:"not null" json:"name"`                      // 用户姓名
	Username       string         `gorm:"uniqueIndex" json:"username,omitempty"`     // 用户名（唯一）
	Email          string         `gorm:"uniqueIndex;not null" json:"email"`         // 用户邮箱（唯一）
	Phone          string         `gorm:"index" json:"phone,omitempty"`              // 手机号
	PasswordHash   string         `gorm:"not null" json:"-"`                         // 密码哈希（不返回给客户端）
	AvatarURL      string         `json:"avatar_url,omitempty"`                       // 头像URL
	Gender         string         `json:"gender,omitempty"`                           // 性别
	Birthday       *time.Time     `json:"birthday,omitempty"`                         // 生日
	Country        string         `json:"country,omitempty"`                          // 国家
	City           string         `json:"city,omitempty"`                             // 城市
	Bio            string         `json:"bio,omitempty"`                              // 个人简介
	Language       string         `gorm:"default:en" json:"language"`                // 语言偏好
	IsVIP          bool           `gorm:"column:is_vip;default:false" json:"is_vip"`               // 是否VIP
	VIPExpiresAt   *time.Time     `gorm:"column:vip_expires_at" json:"vip_expires_at,omitempty"`                  // VIP过期时间
	IsOnline       bool           `gorm:"column:is_online;default:false" json:"is_online"`            // 是否在线
	LastActiveAt   *time.Time     `gorm:"column:last_active_at" json:"last_active_at,omitempty"`                  // 最后活跃时间
	Status         string         `gorm:"default:active" json:"status"`              // 用户状态
	Coins          int            `gorm:"default:0" json:"coins"`                    // 虚拟货币余额
	Fingerprint    string         `json:"-"`                       // 设备指纹
	Roles          []Role         `gorm:"many2many:user_roles;" json:"-"`            // 用户角色列表（多对多关系）
	CreatedAt      time.Time      `json:"created_at"`                                // 创建时间
	UpdatedAt      time.Time      `json:"updated_at"`                                // 更新时间
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`                            // 软删除时间
}

// TableName 指定用户模型对应的数据库表名
func (User) TableName() string {
	return "users"
}

// HasRole 检查用户是否拥有指定角色
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

// IsAdmin 检查用户是否为管理员
func (u *User) IsAdmin() bool {
	return u.HasRole(RoleAdmin)
}

// GetRoleNames 返回用户的所有角色名称列表
func (u *User) GetRoleNames() []string {
	roleNames := make([]string, len(u.Roles))
	for i, role := range u.Roles {
		roleNames[i] = role.Name
	}
	return roleNames
}
