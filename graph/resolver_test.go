package graph

import (
	"context"
	"errors"
	"testing"

	"github.com/NGerasimovvv/GraphQL/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockPostGateway struct {
	CreatePostFunc  func(ctx context.Context, id string, textPost string, commentable bool, authorPost string) (*models.Post, error)
	GetAllPostsFunc func(ctx context.Context, limit *int, offset *int) ([]*models.Post, error)
	GetPostByIDFunc func(ctx context.Context, id string) (*models.Post, error)
}

func (m *MockPostGateway) CreatePost(ctx context.Context, id string, textPost string, commentable bool, authorPost string) (*models.Post, error) {
	return m.CreatePostFunc(ctx, id, textPost, commentable, authorPost)
}

func (m *MockPostGateway) GetAllPosts(ctx context.Context, limit *int, offset *int) ([]*models.Post, error) {
	return m.GetAllPostsFunc(ctx, limit, offset)
}

func (m *MockPostGateway) GetPostByID(ctx context.Context, id string) (*models.Post, error) {
	return m.GetPostByIDFunc(ctx, id)
}

type MockCommentGateway struct {
	CreateCommentFunc         func(ctx context.Context, commentText string, itemID string, authorComment string) (*models.CommentResponse, error)
	GetCommentsByPostIDFunc   func(ctx context.Context, postID string, limit *int, offset *int) ([]*models.CommentResponse, error)
	GetCommentsByParentIDFunc func(ctx context.Context, parentID string, limit *int, offset *int) ([]*models.CommentResponse, error)
	GetCommentByIDFunc        func(ctx context.Context, id string) (*models.CommentResponse, error)
	GetAllCommentsFunc        func(ctx context.Context, limit *int, offset *int) ([]*models.CommentResponse, error)
}

func (m *MockCommentGateway) CreateComment(ctx context.Context, commentText string, itemID string, authorComment string) (*models.CommentResponse, error) {
	return m.CreateCommentFunc(ctx, commentText, itemID, authorComment)
}

func (r *Resolver) CreateComment(ctx context.Context, commentText string, itemID string, authorComment string) (*models.CommentResponse, error) {
	return r.CommentGateway.CreateComment(ctx, commentText, itemID, authorComment)
}

func (r *Resolver) CreatePost(ctx context.Context, textPost string, commentable bool, authorPost string) (*models.Post, error) {
	id := uuid.New().String() // Генерация нового ID для поста
	return r.PostGateway.CreatePost(ctx, id, textPost, commentable, authorPost)
}

func (r *Resolver) Posts(ctx context.Context, limit *int, offset *int) ([]*models.Post, error) {
	posts, err := r.PostGateway.GetAllPosts(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		post.Comments, err = r.CommentGateway.GetCommentsByPostID(ctx, post.ID, limit, offset)
		if err != nil {
			return nil, err
		}

		for _, comment := range post.Comments {
			comment.Replies, err = r.CommentGateway.GetCommentsByParentID(ctx, comment.ID, limit, offset)
			if err != nil {
				return nil, err
			}
		}
	}

	return posts, nil
}

func (r *Resolver) Post(ctx context.Context, id string, limit *int, offset *int) (*models.Post, error) {
	post, err := r.PostGateway.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}

	post.Comments, err = r.CommentGateway.GetCommentsByPostID(ctx, post.ID, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, comment := range post.Comments {
		comment.Replies, err = r.CommentGateway.GetCommentsByParentID(ctx, comment.ID, limit, offset)
		if err != nil {
			return nil, err
		}
	}

	return post, nil
}

func (r *Resolver) Comments(ctx context.Context, limit *int, offset *int) ([]*models.CommentResponse, error) {
	comments, err := r.CommentGateway.GetAllComments(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, comment := range comments {
		comment.Replies, err = r.CommentGateway.GetCommentsByParentID(ctx, comment.ID, limit, offset)
		if err != nil {
			return nil, err
		}
	}

	return comments, nil
}

func (r *Resolver) Comment(ctx context.Context, id string, limit *int, offset *int) (*models.CommentResponse, error) {
	comment, err := r.CommentGateway.GetCommentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	comment.Replies, err = r.CommentGateway.GetCommentsByParentID(ctx, comment.ID, limit, offset)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (m *MockCommentGateway) GetCommentsByPostID(ctx context.Context, postID string, limit *int, offset *int) ([]*models.CommentResponse, error) {
	return m.GetCommentsByPostIDFunc(ctx, postID, limit, offset)
}

func (m *MockCommentGateway) GetCommentsByParentID(ctx context.Context, parentID string, limit *int, offset *int) ([]*models.CommentResponse, error) {
	return m.GetCommentsByParentIDFunc(ctx, parentID, limit, offset)
}

func (m *MockCommentGateway) GetCommentByID(ctx context.Context, id string) (*models.CommentResponse, error) {
	return m.GetCommentByIDFunc(ctx, id)
}

func (m *MockCommentGateway) GetAllComments(ctx context.Context, limit *int, offset *int) ([]*models.CommentResponse, error) {
	return m.GetAllCommentsFunc(ctx, limit, offset)
}

func TestCreatePost(t *testing.T) {
	mockPostGateway := &MockPostGateway{
		CreatePostFunc: func(ctx context.Context, id string, textPost string, commentable bool, authorPost string) (*models.Post, error) {
			return &models.Post{ID: id, TextPost: textPost, Commentable: commentable, AuthorPost: authorPost}, nil
		},
	}

	resolver := &Resolver{PostGateway: mockPostGateway}

	post, err := resolver.CreatePost(context.Background(), "Тестовый пост", true, "author1")

	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, "Тестовый пост", post.TextPost)
	assert.Equal(t, true, post.Commentable)
	assert.Equal(t, "author1", post.AuthorPost)
}

func TestCreateComment(t *testing.T) {
	mockCommentGateway := &MockCommentGateway{
		CreateCommentFunc: func(ctx context.Context, commentText string, itemID string, authorComment string) (*models.CommentResponse, error) {
			return &models.CommentResponse{ID: uuid.New().String(), TextComment: commentText, AuthorComment: authorComment}, nil
		},
	}

	resolver := &Resolver{CommentGateway: mockCommentGateway}

	comment, err := resolver.CreateComment(context.Background(), "Тестовый комментарий", "postID", "author2")

	assert.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, "Тестовый комментарий", comment.TextComment)
	assert.Equal(t, "author2", comment.AuthorComment)
}

func TestPosts(t *testing.T) {
	mockPostGateway := &MockPostGateway{
		GetAllPostsFunc: func(ctx context.Context, limit *int, offset *int) ([]*models.Post, error) {
			return []*models.Post{
				{ID: "1", TextPost: "Post 1", Commentable: true, AuthorPost: "author1"},
				{ID: "2", TextPost: "Post 2", Commentable: false, AuthorPost: "author2"},
			}, nil
		},
	}

	mockCommentGateway := &MockCommentGateway{
		GetCommentsByPostIDFunc: func(ctx context.Context, postID string, limit *int, offset *int) ([]*models.CommentResponse, error) {
			return []*models.CommentResponse{
				{ID: uuid.New().String(), TextComment: "Comment 1", AuthorComment: "author3"},
				{ID: uuid.New().String(), TextComment: "Comment 2", AuthorComment: "author4"},
			}, nil
		},
		GetCommentsByParentIDFunc: func(ctx context.Context, parentID string, limit *int, offset *int) ([]*models.CommentResponse, error) {
			return nil, nil
		},
	}

	resolver := &Resolver{PostGateway: mockPostGateway, CommentGateway: mockCommentGateway}

	posts, err := resolver.Posts(context.Background(), nil, nil)

	assert.NoError(t, err)
	assert.Len(t, posts, 2)
	assert.Equal(t, "1", posts[0].ID)
	assert.Equal(t, "Post 1", posts[0].TextPost)
	assert.Len(t, posts[0].Comments, 2)
	assert.Equal(t, "2", posts[1].ID)
	assert.Equal(t, "Post 2", posts[1].TextPost)
	assert.Len(t, posts[1].Comments, 2)
}

func TestPost(t *testing.T) {
	mockPostGateway := &MockPostGateway{
		GetPostByIDFunc: func(ctx context.Context, id string) (*models.Post, error) {
			return &models.Post{ID: id, TextPost: "Тестовый пост", Commentable: true, AuthorPost: "author1"}, nil
		},
	}

	mockCommentGateway := &MockCommentGateway{
		GetCommentsByPostIDFunc: func(ctx context.Context, postID string, limit *int, offset *int) ([]*models.CommentResponse, error) {
			return []*models.CommentResponse{
				{ID: uuid.New().String(), TextComment: "Комментарий 1", AuthorComment: "author2"},
				{ID: uuid.New().String(), TextComment: "Комментарий 2", AuthorComment: "author3"},
			}, nil
		},
		GetCommentsByParentIDFunc: func(ctx context.Context, parentID string, limit *int, offset *int) ([]*models.CommentResponse, error) {
			return []*models.CommentResponse{
				{ID: uuid.New().String(), TextComment: "Ответ 1", AuthorComment: "author4"},
				{ID: uuid.New().String(), TextComment: "Ответ 2", AuthorComment: "author5"},
			}, nil
		},
	}

	resolver := &Resolver{PostGateway: mockPostGateway, CommentGateway: mockCommentGateway}

	post, err := resolver.Post(context.Background(), "postID", nil, nil)

	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, "Тестовый пост", post.TextPost)
	assert.Len(t, post.Comments, 2)
	assert.Len(t, post.Comments[0].Replies, 2)
	assert.Len(t, post.Comments[1].Replies, 2)
}

func TestComments(t *testing.T) {
	mockCommentGateway := &MockCommentGateway{
		GetAllCommentsFunc: func(ctx context.Context, limit *int, offset *int) ([]*models.CommentResponse, error) {
			return []*models.CommentResponse{
				{ID: uuid.New().String(), TextComment: "Комментарий 1", AuthorComment: "author1"},
				{ID: uuid.New().String(), TextComment: "Комментарий 2", AuthorComment: "author2"},
			}, nil
		},
		GetCommentsByParentIDFunc: func(ctx context.Context, parentID string, limit *int, offset *int) ([]*models.CommentResponse, error) {
			return []*models.CommentResponse{
				{ID: uuid.New().String(), TextComment: "Ответ 1", AuthorComment: "author3"},
				{ID: uuid.New().String(), TextComment: "Ответ 2", AuthorComment: "author4"},
			}, nil
		},
	}

	resolver := &Resolver{CommentGateway: mockCommentGateway}

	comments, err := resolver.Comments(context.Background(), nil, nil)

	assert.NoError(t, err)
	assert.Len(t, comments, 2)
	assert.Len(t, comments[0].Replies, 2)
	assert.Len(t, comments[1].Replies, 2)
}

func TestComment(t *testing.T) {
	mockCommentGateway := &MockCommentGateway{
		GetCommentByIDFunc: func(ctx context.Context, id string) (*models.CommentResponse, error) {
			return &models.CommentResponse{ID: id, TextComment: "Тестовый комментарий", AuthorComment: "author1"}, nil
		},
		GetCommentsByParentIDFunc: func(ctx context.Context, parentID string, limit *int, offset *int) ([]*models.CommentResponse, error) {
			return []*models.CommentResponse{
				{ID: uuid.New().String(), TextComment: "Ответ 1", AuthorComment: "author2"},
				{ID: uuid.New().String(), TextComment: "Ответ 2", AuthorComment: "author3"},
			}, nil
		},
	}

	resolver := &Resolver{CommentGateway: mockCommentGateway}

	comment, err := resolver.Comment(context.Background(), "commentID", nil, nil)

	assert.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, "Тестовый комментарий", comment.TextComment)
	assert.Len(t, comment.Replies, 2)
}

func TestCreatePost_Error(t *testing.T) {
	mockPostGateway := &MockPostGateway{
		CreatePostFunc: func(ctx context.Context, id string, textPost string, commentable bool, authorPost string) (*models.Post, error) {
			return nil, errors.New("ошибка создания поста")
		},
	}

	resolver := &Resolver{PostGateway: mockPostGateway}

	post, err := resolver.CreatePost(context.Background(), "Тестовый пост", true, "author1")

	assert.Error(t, err)
	assert.Nil(t, post)
}
