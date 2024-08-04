// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

type Comment struct {
	ID            string `json:"id"`
	TextComment   string `json:"textComment"`
	PostID        string `json:"postId"`
	AuthorComment string `json:"authorComment"`
}

type CommentResponse struct {
	ID              string             `json:"id"`
	TextComment     string             `json:"textComment"`
	PostID          string             `json:"postId"`
	ParentCommentID *string            `json:"parentCommentID,omitempty"`
	AuthorComment   string             `json:"authorComment"`
	Replies         []*CommentResponse `json:"replies"`
}

type Mutation struct {
}

type Post struct {
	ID          string             `json:"id"`
	TextPost    string             `json:"textPost"`
	AuthorPost  string             `json:"authorPost"`
	Comments    []*CommentResponse `json:"comments"`
	Commentable bool               `json:"commentable"`
}

type Query struct {
}
