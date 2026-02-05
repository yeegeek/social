package post

import (
	"time"

	"gorm.io/gorm"
)

// Repository 动态仓库接口
type Repository interface {
	// Post operations
	CreatePost(post *Post) error
	GetPostByID(id int) (*Post, error)
	GetPosts(userID *int, visibility string, limit, offset int) ([]Post, int64, error)
	GetUserPosts(userID int, limit, offset int) ([]Post, int64, error)
	UpdatePost(post *Post) error
	DeletePost(postID int) error
	IncrementViewCount(postID int) error

	// Like operations
	LikePost(userID, postID int) error
	UnlikePost(userID, postID int) error
	IsPostLiked(userID, postID int) (bool, error)
	GetPostLikes(postID int, limit, offset int) ([]Like, int64, error)

	// Comment operations
	CreateComment(comment *Comment) error
	GetCommentByID(id int) (*Comment, error)
	GetPostComments(postID int, limit, offset int) ([]Comment, int64, error)
	UpdateComment(comment *Comment) error
	DeleteComment(commentID int) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository 创建动态仓库
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreatePost 创建动态
func (r *repository) CreatePost(post *Post) error {
	return r.db.Create(post).Error
}

// GetPostByID 根据ID获取动态
func (r *repository) GetPostByID(id int) (*Post, error) {
	var post Post
	err := r.db.Where("deleted_at IS NULL").First(&post, id).Error
	return &post, err
}

// GetPosts 获取动态列表
func (r *repository) GetPosts(userID *int, visibility string, limit, offset int) ([]Post, int64, error) {
	var posts []Post
	var total int64

	query := r.db.Model(&Post{}).Where("deleted_at IS NULL")

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if visibility != "" {
		query = query.Where("visibility = ?", visibility)
	} else {
		query = query.Where("visibility = ?", "public")
	}

	query = query.Order("created_at DESC")

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Limit(limit).Offset(offset).Find(&posts).Error
	return posts, total, err
}

// GetUserPosts 获取用户的动态列表
func (r *repository) GetUserPosts(userID int, limit, offset int) ([]Post, int64, error) {
	var posts []Post
	var total int64

	query := r.db.Model(&Post{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC")

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Limit(limit).Offset(offset).Find(&posts).Error
	return posts, total, err
}

// UpdatePost 更新动态
func (r *repository) UpdatePost(post *Post) error {
	return r.db.Save(post).Error
}

// DeletePost 删除动态（软删除）
func (r *repository) DeletePost(postID int) error {
	now := time.Now()
	return r.db.Model(&Post{}).Where("id = ?", postID).Update("deleted_at", now).Error
}

// IncrementViewCount 增加浏览次数
func (r *repository) IncrementViewCount(postID int) error {
	return r.db.Model(&Post{}).Where("id = ?", postID).
		UpdateColumn("views_count", gorm.Expr("views_count + ?", 1)).Error
}

// LikePost 点赞动态
func (r *repository) LikePost(userID, postID int) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建点赞记录
	like := &Like{
		UserID: userID,
		PostID: postID,
	}
	if err := tx.Create(like).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 增加点赞数
	if err := tx.Model(&Post{}).Where("id = ?", postID).
		UpdateColumn("likes_count", gorm.Expr("likes_count + ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// UnlikePost 取消点赞
func (r *repository) UnlikePost(userID, postID int) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除点赞记录
	if err := tx.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&Like{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 减少点赞数
	if err := tx.Model(&Post{}).Where("id = ?", postID).
		UpdateColumn("likes_count", gorm.Expr("likes_count - ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// IsPostLiked 检查是否已点赞
func (r *repository) IsPostLiked(userID, postID int) (bool, error) {
	var count int64
	err := r.db.Model(&Like{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error
	return count > 0, err
}

// GetPostLikes 获取动态的点赞列表
func (r *repository) GetPostLikes(postID int, limit, offset int) ([]Like, int64, error) {
	var likes []Like
	var total int64

	query := r.db.Model(&Like{}).Where("post_id = ?", postID).Order("created_at DESC")

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Limit(limit).Offset(offset).Find(&likes).Error
	return likes, total, err
}

// CreateComment 创建评论
func (r *repository) CreateComment(comment *Comment) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建评论
	if err := tx.Create(comment).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 增加评论数
	if err := tx.Model(&Post{}).Where("id = ?", comment.PostID).
		UpdateColumn("comments_count", gorm.Expr("comments_count + ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetCommentByID 根据ID获取评论
func (r *repository) GetCommentByID(id int) (*Comment, error) {
	var comment Comment
	err := r.db.Where("deleted_at IS NULL").First(&comment, id).Error
	return &comment, err
}

// GetPostComments 获取动态的评论列表
func (r *repository) GetPostComments(postID int, limit, offset int) ([]Comment, int64, error) {
	var comments []Comment
	var total int64

	query := r.db.Model(&Comment{}).
		Where("post_id = ? AND deleted_at IS NULL", postID).
		Order("created_at DESC")

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Limit(limit).Offset(offset).Find(&comments).Error
	return comments, total, err
}

// UpdateComment 更新评论
func (r *repository) UpdateComment(comment *Comment) error {
	return r.db.Save(comment).Error
}

// DeleteComment 删除评论（软删除）
func (r *repository) DeleteComment(commentID int) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取评论信息
	var comment Comment
	if err := tx.First(&comment, commentID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 软删除评论
	now := time.Now()
	if err := tx.Model(&Comment{}).Where("id = ?", commentID).Update("deleted_at", now).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 减少评论数
	if err := tx.Model(&Post{}).Where("id = ?", comment.PostID).
		UpdateColumn("comments_count", gorm.Expr("comments_count - ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
