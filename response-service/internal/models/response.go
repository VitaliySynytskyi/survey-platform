package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Answer represents a single answer to a question within a survey response
type Answer struct {
	QuestionID int         `bson:"questionId" json:"questionId"`
	Value      interface{} `bson:"value" json:"value"` // Can be string or []string for checkboxes
}

// Response represents a set of answers submitted for a survey
type Response struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SurveyID    int                `bson:"surveyId" json:"surveyId"`
	UserID      *int               `bson:"userId,omitempty" json:"userId,omitempty"` // Pointer for optionality
	SubmittedAt time.Time          `bson:"submittedAt" json:"submittedAt"`
	Answers     []Answer           `bson:"answers" json:"answers"`
}

// CreateResponseRequest defines the structure for submitting a new response
type CreateResponseRequest struct {
	SurveyID int      `json:"surveyId" binding:"required"`
	UserID   *int     `json:"userId,omitempty"` // Optional, for logged-in users
	Answers  []Answer `json:"answers" binding:"required,dive"`
}
