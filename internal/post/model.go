package post

import (
	"time"

	"github.com/lib/pq"
)

// Post 动态模型
type Post struct {
	ID            int            `gorm:"primaryKey" json:"id"`
	UserID        int            `gorm:"not null" json:"user_id"`
	Content       string         `gorm:"not null" json:"content"`
	MediaURLs     pq.StringArray `gorm:"type:text[]" json:"media_urls,omitempty"`
	Location      string         `json:"location,omitempty"`
	Latitude      *float64       `json:"latitude,omitempty"`
	Longitude     *float64       `json:"longitude,omitempty"`
	Visibility    string         `gorm:"default:public" json:"visibility"` // public, friends, private
	IsPaid        bool           `gorm:"default:false" json:"is_paid"`
	Price         float64        `gorm:"default:0" json:"price"`
	LikesCount    int            `gorm:"default:0" json:"likes_count"`
	CommentsCount int            `gorm:"default:0" json:"comments_count"`
	SharesCount   int            `gorm:"default:0" json:"shares_count"`
	ViewsCount    int            `gorm:"default:0" json:"views_count"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     *time.Time     `json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (Post) TableName() string {
	return "posts"
}

// Like 点赞模型
type Like struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	UserID    int       `gorm:"not null" json:"user_id"`
	PostID    int       `gorm:"not null" json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (Like) TableName() string {
	return "likes"
}

// Comment 评论模型
type Comment struct {
	ID         int        `gorm:"primaryKey" json:"id"`
	PostID     int        `gorm:"not null" json:"post_id"`
	UserID     int        `gorm:"not null" json:"user_id"`
	ParentID   *int       `json:"parent_id,omitempty"`
	Content    string     `gorm:"not null" json:"content"`
	LikesCount int        `gorm:"default:0" json:"likes_count"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}
