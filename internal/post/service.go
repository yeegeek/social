package post

import (
	"errors"

	"github.com/lib/pq"
)

// Service 动态服务接口
type Service interface {
	// Post operations
	CreatePost(userID int, req *CreatePostRequest) (*PostResponse, error)
	GetPost(postID, userID int) (*PostResponse, error)
	GetPosts(userID *int, page, pageSize int) (*PostListResponse, error)
	GetUserPosts(targetUserID, currentUserID int, page, pageSize int) (*PostListResponse, error)
	UpdatePost(postID, userID int, req *UpdatePostRequest) (*PostResponse, error)
	DeletePost(postID, userID int) error

	// Like operations
	LikePost(userID, postID int) error
	UnlikePost(userID, postID int) error

	// Comment operations
	CreateComment(userID int, req *CreateCommentRequest) (*CommentResponse, error)
	GetComments(postID int, page, pageSize int) (*CommentListResponse, error)
	DeleteComment(commentID, userID int) error
}

type service struct {
	repo Repository
}

// NewService 创建动态服务
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreatePost 创建动态
func (s *service) CreatePost(userID int, req *CreatePostRequest) (*PostResponse, error) {
	visibility := req.Visibility
	if visibility == "" {
		visibility = "public"
	}

	post := &Post{
		UserID:     userID,
		Content:    req.Content,
		MediaURLs:  pq.StringArray(req.MediaURLs),
		Location:   req.Location,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		Visibility: visibility,
		IsPaid:     req.IsPaid,
		Price:      req.Price,
	}

	err := s.repo.CreatePost(post)
	if err != nil {
		return nil, err
	}

	return s.toPostResponse(post, userID)
}

// GetPost 获取动态详情
func (s *service) GetPost(postID, userID int) (*PostResponse, error) {
	post, err := s.repo.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	// 增加浏览次数
	_ = s.repo.IncrementViewCount(postID)

	return s.toPostResponse(post, userID)
}

// GetPosts 获取动态列表
func (s *service) GetPosts(userID *int, page, pageSize int) (*PostListResponse, error) {
	offset := (page - 1) * pageSize
	posts, total, err := s.repo.GetPosts(userID, "public", pageSize, offset)
	if err != nil {
		return nil, err
	}

	response := &PostListResponse{
		Posts: make([]PostResponse, 0, len(posts)),
		Total: total,
	}

	currentUserID := 0
	if userID != nil {
		currentUserID = *userID
	}

	for _, post := range posts {
		postResp, err := s.toPostResponse(&post, currentUserID)
		if err == nil {
			response.Posts = append(response.Posts, *postResp)
		}
	}

	return response, nil
}

// GetUserPosts 获取用户的动态列表
func (s *service) GetUserPosts(targetUserID, currentUserID int, page, pageSize int) (*PostListResponse, error) {
	offset := (page - 1) * pageSize
	posts, total, err := s.repo.GetUserPosts(targetUserID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	response := &PostListResponse{
		Posts: make([]PostResponse, 0, len(posts)),
		Total: total,
	}

	for _, post := range posts {
		// 检查可见性
		if post.Visibility == "private" && post.UserID != currentUserID {
			continue
		}

		postResp, err := s.toPostResponse(&post, currentUserID)
		if err == nil {
			response.Posts = append(response.Posts, *postResp)
		}
	}

	return response, nil
}

// UpdatePost 更新动态
func (s *service) UpdatePost(postID, userID int, req *UpdatePostRequest) (*PostResponse, error) {
	post, err := s.repo.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	// 验证权限
	if post.UserID != userID {
		return nil, errors.New("unauthorized to update post")
	}

	// 更新字段
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.MediaURLs != nil {
		post.MediaURLs = pq.StringArray(req.MediaURLs)
	}
	if req.Visibility != "" {
		post.Visibility = req.Visibility
	}

	err = s.repo.UpdatePost(post)
	if err != nil {
		return nil, err
	}

	return s.toPostResponse(post, userID)
}

// DeletePost 删除动态
func (s *service) DeletePost(postID, userID int) error {
	post, err := s.repo.GetPostByID(postID)
	if err != nil {
		return err
	}

	// 验证权限
	if post.UserID != userID {
		return errors.New("unauthorized to delete post")
	}

	return s.repo.DeletePost(postID)
}

// LikePost 点赞动态
func (s *service) LikePost(userID, postID int) error {
	// 检查是否已点赞
	isLiked, err := s.repo.IsPostLiked(userID, postID)
	if err != nil {
		return err
	}

	if isLiked {
		return errors.New("already liked")
	}

	return s.repo.LikePost(userID, postID)
}

// UnlikePost 取消点赞
func (s *service) UnlikePost(userID, postID int) error {
	// 检查是否已点赞
	isLiked, err := s.repo.IsPostLiked(userID, postID)
	if err != nil {
		return err
	}

	if !isLiked {
		return errors.New("not liked yet")
	}

	return s.repo.UnlikePost(userID, postID)
}

// CreateComment 创建评论
func (s *service) CreateComment(userID int, req *CreateCommentRequest) (*CommentResponse, error) {
	comment := &Comment{
		PostID:   req.PostID,
		UserID:   userID,
		ParentID: req.ParentID,
		Content:  req.Content,
	}

	err := s.repo.CreateComment(comment)
	if err != nil {
		return nil, err
	}

	return s.toCommentResponse(comment), nil
}

// GetComments 获取评论列表
func (s *service) GetComments(postID int, page, pageSize int) (*CommentListResponse, error) {
	offset := (page - 1) * pageSize
	comments, total, err := s.repo.GetPostComments(postID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	response := &CommentListResponse{
		Comments: make([]CommentResponse, 0, len(comments)),
		Total:    total,
	}

	for _, comment := range comments {
		response.Comments = append(response.Comments, *s.toCommentResponse(&comment))
	}

	return response, nil
}

// DeleteComment 删除评论
func (s *service) DeleteComment(commentID, userID int) error {
	comment, err := s.repo.GetCommentByID(commentID)
	if err != nil {
		return err
	}

	// 验证权限
	if comment.UserID != userID {
		return errors.New("unauthorized to delete comment")
	}

	return s.repo.DeleteComment(commentID)
}

// toPostResponse 转换为动态响应
func (s *service) toPostResponse(post *Post, userID int) (*PostResponse, error) {
	isLiked := false
	if userID > 0 {
		isLiked, _ = s.repo.IsPostLiked(userID, post.ID)
	}

	return &PostResponse{
		ID:            post.ID,
		UserID:        post.UserID,
		Content:       post.Content,
		MediaURLs:     post.MediaURLs,
		Location:      post.Location,
		Latitude:      post.Latitude,
		Longitude:     post.Longitude,
		Visibility:    post.Visibility,
		IsPaid:        post.IsPaid,
		Price:         post.Price,
		LikesCount:    post.LikesCount,
		CommentsCount: post.CommentsCount,
		SharesCount:   post.SharesCount,
		ViewsCount:    post.ViewsCount,
		IsLiked:       isLiked,
		CreatedAt:     post.CreatedAt,
		UpdatedAt:     post.UpdatedAt,
	}, nil
}

// toCommentResponse 转换为评论响应
func (s *service) toCommentResponse(comment *Comment) *CommentResponse {
	return &CommentResponse{
		ID:         comment.ID,
		PostID:     comment.PostID,
		UserID:     comment.UserID,
		ParentID:   comment.ParentID,
		Content:    comment.Content,
		LikesCount: comment.LikesCount,
		CreatedAt:  comment.CreatedAt,
	}
}
