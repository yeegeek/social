package post

import "time"

// CreatePostRequest 创建动态请求
type CreatePostRequest struct {
	Content    string   `json:"content" validate:"required"`
	MediaURLs  []string `json:"media_urls,omitempty"`
	Location   string   `json:"location,omitempty"`
	Latitude   *float64 `json:"latitude,omitempty"`
	Longitude  *float64 `json:"longitude,omitempty"`
	Visibility string   `json:"visibility" validate:"omitempty,oneof=public friends private"`
	IsPaid     bool     `json:"is_paid,omitempty"`
	Price      float64  `json:"price,omitempty"`
}

// UpdatePostRequest 更新动态请求
type UpdatePostRequest struct {
	Content    string   `json:"content,omitempty"`
	MediaURLs  []string `json:"media_urls,omitempty"`
	Visibility string   `json:"visibility,omitempty"`
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	PostID   int    `json:"post_id" validate:"required"`
	ParentID *int   `json:"parent_id,omitempty"`
	Content  string `json:"content" validate:"required"`
}

// PostResponse 动态响应
type PostResponse struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	Username      string    `json:"username,omitempty"`
	UserAvatar    string    `json:"user_avatar,omitempty"`
	Content       string    `json:"content"`
	MediaURLs     []string  `json:"media_urls,omitempty"`
	Location      string    `json:"location,omitempty"`
	Latitude      *float64  `json:"latitude,omitempty"`
	Longitude     *float64  `json:"longitude,omitempty"`
	Visibility    string    `json:"visibility"`
	IsPaid        bool      `json:"is_paid"`
	Price         float64   `json:"price,omitempty"`
	LikesCount    int       `json:"likes_count"`
	CommentsCount int       `json:"comments_count"`
	SharesCount   int       `json:"shares_count"`
	ViewsCount    int       `json:"views_count"`
	IsLiked       bool      `json:"is_liked"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CommentResponse 评论响应
type CommentResponse struct {
	ID         int       `json:"id"`
	PostID     int       `json:"post_id"`
	UserID     int       `json:"user_id"`
	Username   string    `json:"username,omitempty"`
	UserAvatar string    `json:"user_avatar,omitempty"`
	ParentID   *int      `json:"parent_id,omitempty"`
	Content    string    `json:"content"`
	LikesCount int       `json:"likes_count"`
	CreatedAt  time.Time `json:"created_at"`
}

// PostListResponse 动态列表响应
type PostListResponse struct {
	Posts []PostResponse `json:"posts"`
	Total int64          `json:"total"`
}

// CommentListResponse 评论列表响应
type CommentListResponse struct {
	Comments []CommentResponse `json:"comments"`
	Total    int64             `json:"total"`
}
