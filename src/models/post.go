package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// VoteQuestion represents a question in a vote post.
type VoteQuestion struct {
	Question string   `bson:"question,omitempty" json:"question,omitempty"`
	Options  []string `bson:"options,omitempty" json:"options,omitempty"`
}

type FormQuestion struct {
	Question  string   `bson:"question" json:"question"`
	InputType string   `bson:"inputType" json:"inputType"`
	MaxSel    string   `bson:"maxSel,omitempty" json:"maxSel,omitempty"`
	Options   []string `bson:"options" json:"options"`
}

// Post represents a general post.
type Post struct {
	PostID        string              `bson:"postID,omitempty" json:"postID,omitempty"` //use get postID in string format
	ID            primitive.ObjectID  `bson:"_id" json:"_id"`
	Kind          string              `bson:"kind" json:"kind"`
	AssignTo      []string            `bson:"assignTo" json:"assignTo"`
	Public        bool                `bson:"public" json:"public"`
	Title         string              `bson:"title" json:"title"`
	Description   string              `bson:"description" json:"description"`
	PostDate      primitive.DateTime  `bson:"postDate" json:"postDate"`
	EndDate       *primitive.DateTime `bson:"endDate" json:"endDate,omitempty"` // Nullable
	Author        string              `bson:"author" json:"author"`
	Markdown      string              `bson:"markdown,omitempty" json:"markdown,omitempty"`
	FormQuestions []FormQuestion      `bson:"formQuestions,omitempty" json:"formQuestions,omitempty"` // For form posts
	VoteQuestions VoteQuestion        `bson:"voteQuestions,omitempty" json:"voteQuestions,omitempty"` // For vote posts
}

// PPost extends Post for regular posts.
type PPost struct {
	Post
	TimeUp bool `bson:"timeUp" json:"timeUp"`
}

// PVote extends Post for vote posts.
type PVote struct {
	Post
	Questions VoteQuestion `bson:"voteQuestions,omitempty" json:"voteQuestions,omitempty"` // Include questions for vote posts
	TimeUp    bool         `bson:"timeUp" json:"timeUp"`
}

type PForm struct {
	Post
	Questions []FormQuestion `bson:"formQuestions,omitempty" json:"formQuestions,omitempty"`
	TimeUp    bool           `bson:"timeUp" json:"timeUp"`
}

type CreatePostRequest struct {
	EventID     primitive.ObjectID `bson:"eventID" json:"eventID"`
	UpdatedPost Post               `bson:"updatedPost" json:"updatedPost"`
}
