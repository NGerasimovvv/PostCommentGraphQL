package gateway

import (
	"context"

	"github.com/NGerasimovvv/GraphQL/internal/models"
	"github.com/NGerasimovvv/GraphQL/internal/storage"
)

type CommentGateway interface {
	GetAllComments(ctx context.Context, limit, offset *int) ([]*models.CommentResponse, error)
	GetCommentByID(ctx context.Context, id string) (*models.CommentResponse, error)
	CreateComment(ctx context.Context, commentText, itemId, user string) (*models.CommentResponse, error)
	GetCommentsByPostID(ctx context.Context, postID string, limit, offset *int) ([]*models.CommentResponse, error)
	GetCommentsByParentID(ctx context.Context, parentID string, limit, offset *int) ([]*models.CommentResponse, error)
}

type commentGateway struct {
	storage storage.Storage
}

func NewCommentGateway(storage storage.Storage) CommentGateway {
	return &commentGateway{storage: storage}
}

func (s *commentGateway) GetAllComments(ctx context.Context, limit, offset *int) ([]*models.CommentResponse, error) {
	return s.storage.GetAllComments(ctx, limit, offset)
}

func (s *commentGateway) GetCommentByID(ctx context.Context, id string) (*models.CommentResponse, error) {
	return s.storage.GetCommentByID(ctx, id)
}

func (s *commentGateway) CreateComment(ctx context.Context, commentText, itemId, user string) (*models.CommentResponse, error) {
	return s.storage.CreateComment(ctx, commentText, itemId, user)
}

func (s *commentGateway) GetCommentsByPostID(ctx context.Context, postID string, limit, offset *int) ([]*models.CommentResponse, error) {
	return s.storage.GetCommentsByPostID(ctx, postID, limit, offset)
}

func (s *commentGateway) GetCommentsByParentID(ctx context.Context, parentID string, limit, offset *int) ([]*models.CommentResponse, error) {
	return s.storage.GetCommentsByParentID(ctx, parentID, limit, offset)
}

type PostGateway interface {
	CreatePost(ctx context.Context, id, text string, commentable bool, authorPost string) (*models.Post, error)
	GetPostByID(ctx context.Context, id string) (*models.Post, error)
	GetAllPosts(ctx context.Context, limit, offset *int) ([]*models.Post, error)
}

type postGateway struct {
	storage storage.Storage
}

func NewPostGateway(storage storage.Storage) PostGateway {
	return &postGateway{storage: storage}
}

func (s *postGateway) CreatePost(ctx context.Context, id, textPost string, commentable bool, authorPost string) (*models.Post, error) {
	return s.storage.CreatePost(ctx, id, textPost, commentable, authorPost)
}

func (s *postGateway) GetPostByID(ctx context.Context, id string) (*models.Post, error) {
	return s.storage.GetPostByID(ctx, id)
}

func (s *postGateway) GetAllPosts(ctx context.Context, limit, offset *int) ([]*models.Post, error) {
	return s.storage.GetAllPosts(ctx, limit, offset)
}
