package graph

import (
	"context"
	"github.com/NGerasimovvv/GraphQL/internal/gateway"
	"github.com/NGerasimovvv/GraphQL/internal/models"
	"github.com/google/uuid"
)

func (r *mutationResolver) CreatePost(ctx context.Context, textPost string, commentable bool, authorPost string) (*models.Post, error) {
	id := uuid.New().String()
	post, err := r.PostGateway.CreatePost(ctx, id, textPost, commentable, authorPost)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, commentText string, itemID string, authorComment string) (*models.CommentResponse, error) {
	comment, err := r.CommentGateway.CreateComment(ctx, commentText, itemID, authorComment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (r *queryResolver) Posts(ctx context.Context, limit *int, offset *int) ([]*models.Post, error) {
	var posts []*models.Post
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

func (r *queryResolver) Post(ctx context.Context, id string, limit *int, offset *int) (*models.Post, error) {
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

func (r *queryResolver) Comments(ctx context.Context, limit *int, offset *int) ([]*models.CommentResponse, error) {
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

func (r *queryResolver) Comment(ctx context.Context, id string, limit *int, offset *int) (*models.CommentResponse, error) {
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

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

type Resolver struct {
	CommentGateway gateway.CommentGateway
	PostGateway    gateway.PostGateway
}
