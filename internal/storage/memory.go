package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/NGerasimovvv/GraphQL/internal/models"
	"github.com/google/uuid"
)

type InMemoryStorage struct {
	postCounter    int
	commentCounter int
	posts          map[string]*models.Post
	comments       map[string]*models.CommentResponse
	mu             sync.RWMutex
}

func NewMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		postCounter:    0,
		commentCounter: 0,
		posts:          make(map[string]*models.Post),
		comments:       make(map[string]*models.CommentResponse),
	}
}

func InitMemoryStorage() *InMemoryStorage {
	storage := NewMemoryStorage()
	storage.mu.Lock()
	defer storage.mu.Unlock()
	return storage
}

func (s *InMemoryStorage) GetAllPosts(ctx context.Context, limit, offset *int) ([]*models.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var posts []*models.Post
	for _, post := range s.posts {
		posts = append(posts, post)
	}

	if limit != nil && offset != nil {
		start := *offset
		end := *offset + *limit
		if start > len(posts) {
			return []*models.Post{}, nil
		}
		if end > len(posts) {
			end = len(posts)
		}
		posts = posts[start:end]
	}

	return posts, nil
}

func pagination(page, pageSize *int) (offset, limit int) {
	if page != nil && *page <= 0 {
		page = nil
	}
	if pageSize != nil && *pageSize < 0 {
		pageSize = nil
	}
	if page == nil || pageSize == nil {
		limit = -1
		offset = 0
	} else {
		offset = (*page - 1) * *pageSize
		limit = *pageSize
	}
	return
}

func (s *InMemoryStorage) GetPostByID(ctx context.Context, postID string) (*models.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	post, exists := s.posts[postID]
	if !exists {
		return nil, fmt.Errorf("post not found")
	}
	return post, nil
}

func (s *InMemoryStorage) CreatePost(ctx context.Context, id, text string, commentable bool, authorPost string) (*models.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	post := &models.Post{ID: id, TextPost: text, Commentable: commentable, AuthorPost: authorPost}
	s.posts[id] = post
	return post, nil
}

func (s *InMemoryStorage) GetAllComments(ctx context.Context, limit, offset *int) ([]*models.CommentResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	comments := make([]*models.CommentResponse, 0, len(s.comments))
	for _, comment := range s.comments {
		comments = append(comments, comment)
	}
	start := 0
	if offset != nil {
		start = *offset
	}
	end := len(comments)
	if limit != nil && start+*limit < end {
		end = start + *limit
	}

	return comments[start:end], nil
}

func (s *InMemoryStorage) GetCommentsByPostID(ctx context.Context, postID string, limit, offset *int) ([]*models.CommentResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var comments []*models.CommentResponse
	for _, comment := range s.comments {
		if comment.PostID == postID {
			comments = append(comments, comment)
		}
	}

	start := 0
	if offset != nil {
		start = *offset
	}
	end := len(comments)
	if limit != nil {
		end = start + *limit
		if end > len(comments) {
			end = len(comments)
		}
	}
	if start > len(comments) {
		return []*models.CommentResponse{}, nil
	}

	return comments[start:end], nil
}

func (s *InMemoryStorage) GetCommentsByParentID(ctx context.Context, parentID string, limit, offset *int) ([]*models.CommentResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var comments []*models.CommentResponse
	for _, comment := range s.comments {
		if comment.ParentCommentID != nil && *comment.ParentCommentID == parentID {
			comments = append(comments, comment)
		}
	}

	start := 0
	if offset != nil {
		start = *offset
	}
	end := len(comments)
	if limit != nil {
		end = start + *limit
		if end > len(comments) {
			end = len(comments)
		}
	}
	if start > len(comments) {
		return []*models.CommentResponse{}, nil
	}

	return comments[start:end], nil
}

func (s *InMemoryStorage) GetCommentByID(ctx context.Context, id string) (*models.CommentResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	comment, exists := s.comments[id]
	if !exists {
		return nil, fmt.Errorf("comment not found")
	}
	return comment, nil
}

func (s *InMemoryStorage) CreateComment(ctx context.Context, commentText, itemId, user string) (*models.CommentResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var isReply bool
	var parentCommentID *string
	var postID string
	var commentAble bool

	if post, exists := s.posts[itemId]; exists {
		postID = itemId
		isReply = false
		commentAble = post.Commentable
		if !commentAble {
			return nil, errors.New("author turned off comments under this post")
		}
	} else if comment, exists := s.comments[itemId]; exists {
		postID = comment.PostID
		parentCommentID = &itemId
		isReply = true
	} else {
		return nil, errors.New("item not found")
	}

	var newComment *models.CommentResponse
	id := uuid.New().String()
	if isReply {
		newComment = &models.CommentResponse{ID: id, TextComment: commentText, AuthorComment: user, PostID: postID, ParentCommentID: parentCommentID}
	} else {
		newComment = &models.CommentResponse{ID: id, TextComment: commentText, AuthorComment: user, PostID: postID}
	}
	s.comments[id] = newComment

	return newComment, nil
}
