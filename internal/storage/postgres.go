package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/NGerasimovvv/GraphQL/internal/config"
	"github.com/NGerasimovvv/GraphQL/internal/models"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	DB *sql.DB
}

func InitPostgresDatabase(cfg *config.Config) *PostgresStorage {
	const op = "postgres.InitPostgresDatabase"
	dbHost := cfg.Postgres.PostgresHost
	dbPort := cfg.Postgres.PostgresPort
	dbUser := cfg.Postgres.PostgresUser
	dbPasswd := cfg.Postgres.PostgresPassword
	dbName := cfg.Postgres.DatabaseName

	postgresUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPasswd, dbName)
	db, err := sql.Open("postgres", postgresUrl)
	if err != nil {
		log.Fatalf("%s: %v", op, err)
	}

	createPostTable, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS post (
		id UUID PRIMARY KEY,
		text TEXT NOT NULL,
		authorPost VARCHAR(50) NOT NULL,
		commentable BOOLEAN NOT NULL);
`)
	if err != nil {
		log.Fatalf("%s: %v", op, err)
	}
	_, err = createPostTable.Exec()
	if err != nil {
		log.Fatalf("%s: %v", op, err)
	}

	createCommentTable, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS comment (
		id UUID PRIMARY KEY,
		comment VARCHAR(2000),
		authorComment VARCHAR(50) NOT NULL,
		post_id UUID NOT NULL,
		parent_comment_id UUID,
		FOREIGN KEY (post_id) REFERENCES post(id),
		FOREIGN KEY (parent_comment_id) REFERENCES comment(id)
	);`)
	if err != nil {
		log.Fatalf("%s: %v", op, err)
	}
	_, err = createCommentTable.Exec()
	if err != nil {
		log.Fatalf("%s: %v", op, err)
	}

	return &PostgresStorage{DB: db}
}

func (s *PostgresStorage) ClosePostgres() error {
	return s.DB.Close()
}

func (s *PostgresStorage) GetAllPosts(ctx context.Context, limit, offset *int) ([]*models.Post, error) {
	query := "SELECT id, text, authorPost, commentable FROM post"
	var rows *sql.Rows
	var err error

	if limit != nil && offset != nil {
		query += " LIMIT $1 OFFSET $2"
		rows, err = s.DB.QueryContext(ctx, query, *limit, *offset)
	} else {
		rows, err = s.DB.QueryContext(ctx, query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.TextPost, &post.AuthorPost, &post.Commentable); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func (s *PostgresStorage) GetPostByID(ctx context.Context, postID string) (*models.Post, error) {
	var post models.Post
	err := s.DB.QueryRowContext(ctx, "SELECT id, text, authorPost, commentable FROM post WHERE id=$1", postID).Scan(&post.ID, &post.TextPost, &post.AuthorPost, &post.Commentable)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *PostgresStorage) CreatePost(ctx context.Context, id, textPost string, commentable bool, authorPost string) (*models.Post, error) {
	_, err := s.DB.ExecContext(ctx, "INSERT INTO post (id, text, authorPost, commentable) VALUES ($1, $2, $3, $4)", id, textPost, authorPost, commentable)
	if err != nil {
		return nil, err
	}
	return &models.Post{ID: id, TextPost: textPost, Commentable: commentable, AuthorPost: authorPost}, nil
}

func (s *PostgresStorage) GetAllComments(ctx context.Context, limit, offset *int) ([]*models.CommentResponse, error) {
	query := "SELECT id, comment, authorComment, post_id, parent_comment_id FROM comment"
	var params []interface{}

	if limit != nil && offset != nil {
		query += " LIMIT $1 OFFSET $2"
		params = append(params, *limit, *offset)
	} else if limit != nil {
		query += " LIMIT $1"
		params = append(params, *limit)
	} else if offset != nil {
		query += " OFFSET $1"
		params = append(params, *offset)
	}

	rows, err := s.DB.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.CommentResponse
	for rows.Next() {
		var comment models.CommentResponse
		if err := rows.Scan(&comment.ID, &comment.TextComment, &comment.AuthorComment, &comment.PostID, &comment.ParentCommentID); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

func (s *PostgresStorage) GetCommentsByPostID(ctx context.Context, postID string, limit, offset *int) ([]*models.CommentResponse, error) {
	query := "SELECT id, comment, authorComment, post_id, parent_comment_id FROM comment WHERE post_id=$1"
	args := []interface{}{postID}

	if limit != nil {
		query += " LIMIT $2"
		args = append(args, *limit)
	}
	if offset != nil {
		query += " OFFSET $3"
		args = append(args, *offset)
	}

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.CommentResponse
	for rows.Next() {
		var comment models.CommentResponse
		err := rows.Scan(&comment.ID, &comment.TextComment, &comment.AuthorComment, &comment.PostID, &comment.ParentCommentID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

func (s *PostgresStorage) GetCommentsByParentID(ctx context.Context, parentID string, limit, offset *int) ([]*models.CommentResponse, error) {
	query := "SELECT id, comment, authorComment, post_id, parent_comment_id FROM comment WHERE parent_comment_id=$1"

	if limit != nil && offset != nil {
		query += " LIMIT $2 OFFSET $3"
		rows, err := s.DB.QueryContext(ctx, query, parentID, *limit, *offset)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return scanerComments(rows)
	}

	rows, err := s.DB.QueryContext(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanerComments(rows)
}

func scanerComments(rows *sql.Rows) ([]*models.CommentResponse, error) {
	var comments []*models.CommentResponse
	for rows.Next() {
		var comment models.CommentResponse
		err := rows.Scan(&comment.ID, &comment.TextComment, &comment.AuthorComment, &comment.PostID, &comment.ParentCommentID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *PostgresStorage) GetCommentByID(ctx context.Context, id string) (*models.CommentResponse, error) {
	var comment models.CommentResponse
	err := s.DB.QueryRowContext(ctx, "SELECT id, comment, authorComment, post_id, parent_comment_id FROM comment WHERE id=$1", id).Scan(&comment.ID, &comment.TextComment, &comment.AuthorComment, &comment.PostID, &comment.ParentCommentID)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (s *PostgresStorage) CreateComment(ctx context.Context, commentText, itemId, user string) (*models.CommentResponse, error) {
	var isReply bool
	var parentCommentID *string
	var postID string
	var commentAble bool

	err := s.DB.QueryRowContext(ctx, "SELECT commentable FROM post WHERE id=$1", itemId).Scan(&commentAble)
	if errors.Is(err, sql.ErrNoRows) {
		err = s.DB.QueryRowContext(ctx, "SELECT post_id FROM comment WHERE id=$1", itemId).Scan(&postID)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("item not found")
		} else if err != nil {
			return nil, err
		} else {
			parentCommentID = &itemId
			isReply = true
		}
	} else if err != nil {
		return nil, err
	} else if !commentAble {
		return nil, errors.New("author turned off comments under this post")
	} else {
		postID = itemId
		isReply = false
	}
	var query string
	id := uuid.New().String()
	if isReply {
		query = "INSERT INTO comment (id, comment, authorComment, post_id, parent_comment_id) VALUES ($1, $2, $3, $4, $5)"
		_, err := s.DB.ExecContext(ctx, query, id, commentText, user, postID, itemId)
		if err != nil {
			return nil, err
		}
		return &models.CommentResponse{ID: id, TextComment: commentText, AuthorComment: user, PostID: postID, ParentCommentID: parentCommentID}, nil
	} else {
		query = "INSERT INTO comment (id, comment, authorComment, post_id, parent_comment_id) VALUES ($1, $2, $3, $4, NULL)"
		_, err := s.DB.ExecContext(ctx, query, id, commentText, user, postID)
		if err != nil {
			return nil, err
		}
		return &models.CommentResponse{ID: id, TextComment: commentText, AuthorComment: user, PostID: postID}, nil
	}
}
