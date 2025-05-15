package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RabbitMQMessage represents a message received from RabbitMQ
type RabbitMQMessage struct {
	SurveyID     string   `json:"survey_id"`
	RespondentID string   `json:"respondent_id,omitempty"`
	AnonymousID  string   `json:"anonymous_id,omitempty"`
	Answers      []Answer `json:"answers"`
	SubmittedAt  string   `json:"submitted_at"`
}

// Answer represents an answer to a survey question
type Answer struct {
	QuestionID string      `json:"question_id" bson:"question_id"`
	Value      interface{} `json:"value" bson:"value"` // can be string, []string, number, etc. depending on question type
}

// SurveyResponse represents a survey response in MongoDB
type SurveyResponse struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	SurveyID     primitive.ObjectID `bson:"survey_id"`
	RespondentID string             `bson:"respondent_id,omitempty"`
	AnonymousID  string             `bson:"anonymous_id,omitempty"`
	Answers      []Answer           `bson:"answers"`
	SubmittedAt  time.Time          `bson:"submitted_at"`
}
