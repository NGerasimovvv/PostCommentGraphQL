package storage

import (
	"context"

	"github.com/NGerasimovvv/GraphQL/internal/config"
	"github.com/NGerasimovvv/GraphQL/internal/models"
)

type Storage interface {
	GetAllPosts(ctx context.Context, limit, offset *int) ([]*models.Post, error)
	GetPostByID(ctx context.Context, postID string) (*models.Post, error)
	CreatePost(ctx context.Context, id, textPost string, commentable bool, authorPost string) (*models.Post, error)

	GetAllComments(ctx context.Context, limit, offset *int) ([]*models.CommentResponse, error)
	GetCommentsByPostID(ctx context.Context, postID string, limit, offset *int) ([]*models.CommentResponse, error) // Обновлено
	GetCommentsByParentID(ctx context.Context, parentID string, limit, offset *int) ([]*models.CommentResponse, error)
	GetCommentByID(ctx context.Context, id string) (*models.CommentResponse, error)
	CreateComment(ctx context.Context, textComment, itemId, user string) (*models.CommentResponse, error)
}

func StorageType(cfg *config.Config) Storage {
	storageType := cfg.Storage.StorageType
	var storage Storage
	if storageType == "memory" {
		storage = InitMemoryStorage()
	} else {
		storage = InitPostgresDatabase(cfg)
	}
	return storage
}
