package db

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service/internal/config"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service/internal/model"
)

// AnalyticsRepository handles operations for analytics data
type AnalyticsRepository struct {
	client        *mongo.Client
	responsesColl *mongo.Collection
	surveysColl   *mongo.Collection
}

// NewAnalyticsRepository creates a new repository for analytics
func NewAnalyticsRepository(cfg config.MongoDB) (*AnalyticsRepository, error) {
	// Create MongoDB client
	clientOptions := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// Get collections
	responsesColl := client.Database(cfg.Database).Collection(cfg.ResponsesColl)
	surveysColl := client.Database(cfg.Database).Collection(cfg.SurveysColl)

	return &AnalyticsRepository{
		client:        client,
		responsesColl: responsesColl,
		surveysColl:   surveysColl,
	}, nil
}

// GetSurveyById retrieves a survey by ID
func (r *AnalyticsRepository) GetSurveyById(ctx context.Context, surveyID string) (*model.Survey, error) {
	id, err := primitive.ObjectIDFromHex(surveyID)
	if err != nil {
		return nil, fmt.Errorf("invalid survey ID: %w", err)
	}

	filter := bson.M{"_id": id}
	var survey model.Survey

	err = r.surveysColl.FindOne(ctx, filter).Decode(&survey)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("survey not found")
		}
		return nil, fmt.Errorf("failed to get survey: %w", err)
	}

	return &survey, nil
}

// GetSurveyResults retrieves aggregated results for a survey
func (r *AnalyticsRepository) GetSurveyResults(ctx context.Context, surveyID string) (*model.SurveyResults, error) {
	// Convert survey ID from string to ObjectID
	id, err := primitive.ObjectIDFromHex(surveyID)
	if err != nil {
		return nil, fmt.Errorf("invalid survey ID: %w", err)
	}

	// Get survey information
	survey, err := r.GetSurveyById(ctx, surveyID)
	if err != nil {
		return nil, err
	}

	// Get total number of responses
	filter := bson.M{"survey_id": id}
	totalResponses, err := r.responsesColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count responses: %w", err)
	}

	// Initialize survey results
	results := &model.SurveyResults{
		SurveyID:        surveyID,
		Title:           survey.Title,
		Description:     survey.Description,
		TotalResponses:  int(totalResponses),
		QuestionResults: make([]model.QuestionResult, 0, len(survey.Questions)),
	}

	// Process each question
	for _, question := range survey.Questions {
		questionResult, err := r.getQuestionResults(ctx, question, id, int(totalResponses))
		if err != nil {
			return nil, fmt.Errorf("failed to get results for question %s: %w", question.ID, err)
		}
		results.QuestionResults = append(results.QuestionResults, *questionResult)
	}

	return results, nil
}

// getQuestionResults processes results for a specific question
func (r *AnalyticsRepository) getQuestionResults(ctx context.Context, question model.SurveyQuestion, surveyID primitive.ObjectID, totalResponses int) (*model.QuestionResult, error) {
	// Create a pipeline for aggregating data for this question
	pipeline := r.buildAggregationPipeline(question, surveyID)

	// Execute the aggregation pipeline
	cursor, err := r.responsesColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregation failed: %w", err)
	}
	defer cursor.Close(ctx)

	// Process the results based on question type
	result := &model.QuestionResult{
		QuestionID:   question.ID,
		Title:        question.Title,
		Type:         question.Type,
		TotalAnswers: 0,
		Analytics:    model.QuestionAnalytics{},
	}

	// Process based on question type
	switch question.Type {
	case "single-choice", "multiple-choice":
		result, err = r.processSingleMultipleChoiceResults(ctx, cursor, question, result, totalResponses)
	case "scale":
		result, err = r.processScaleResults(ctx, cursor, result)
	case "open-text":
		result, err = r.processOpenTextResults(ctx, cursor, result)
	case "matrix":
		result, err = r.processMatrixResults(ctx, cursor, question, result)
	default:
		// For unknown question types, just return basic stats
		result, err = r.processDefaultResults(ctx, cursor, result)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

// buildAggregationPipeline creates a MongoDB aggregation pipeline for a question
func (r *AnalyticsRepository) buildAggregationPipeline(question model.SurveyQuestion, surveyID primitive.ObjectID) mongo.Pipeline {
	// Match documents for this survey
	matchStage := bson.D{
		{Key: "$match", Value: bson.M{"survey_id": surveyID}},
	}

	// Unwind the answers array
	unwindStage := bson.D{
		{Key: "$unwind", Value: bson.M{"path": "$answers"}},
	}

	// Match only answers for this question
	matchQuestionStage := bson.D{
		{Key: "$match", Value: bson.M{"answers.question_id": question.ID}},
	}

	pipeline := mongo.Pipeline{matchStage, unwindStage, matchQuestionStage}

	// Add additional stages based on question type
	switch question.Type {
	case "single-choice", "multiple-choice":
		// For multiple-choice, we need to unwind the values if they're arrays
		if question.Type == "multiple-choice" {
			unwindValueStage := bson.D{
				{Key: "$addFields", Value: bson.M{
					"value_is_array": bson.M{"$isArray": "$answers.value"},
				}},
			}
			conditionalUnwindStage := bson.D{
				{Key: "$project", Value: bson.M{
					"value": bson.M{
						"$cond": bson.M{
							"if":   "$value_is_array",
							"then": bson.M{"$arrayElemAt": []interface{}{"$answers.value", 0}},
							"else": "$answers.value",
						},
					},
					"submitted_at": 1,
				}},
			}
			pipeline = append(pipeline, unwindValueStage, conditionalUnwindStage)
		} else {
			// For single-choice, just project the value
			projectStage := bson.D{
				{Key: "$project", Value: bson.M{
					"value":        "$answers.value",
					"submitted_at": 1,
				}},
			}
			pipeline = append(pipeline, projectStage)
		}

		// Group by value and count occurrences
		groupStage := bson.D{
			{Key: "$group", Value: bson.M{
				"_id":   "$value",
				"count": bson.M{"$sum": 1},
			}},
		}
		pipeline = append(pipeline, groupStage)

	case "scale":
		// Project the value and submitted_at
		projectStage := bson.D{
			{Key: "$project", Value: bson.M{
				"value":        bson.M{"$toDouble": "$answers.value"},
				"submitted_at": 1,
			}},
		}
		// Group all values to calculate statistics
		groupStage := bson.D{
			{Key: "$group", Value: bson.M{
				"_id":    nil,
				"values": bson.M{"$push": "$value"},
				"count":  bson.M{"$sum": 1},
				"avg":    bson.M{"$avg": "$value"},
				"min":    bson.M{"$min": "$value"},
				"max":    bson.M{"$max": "$value"},
			}},
		}
		pipeline = append(pipeline, projectStage, groupStage)

	case "open-text":
		// Project the value and submitted_at
		projectStage := bson.D{
			{Key: "$project", Value: bson.M{
				"value":        "$answers.value",
				"submitted_at": 1,
			}},
		}
		// Group to collect all text responses
		groupStage := bson.D{
			{Key: "$group", Value: bson.M{
				"_id":    nil,
				"values": bson.M{"$push": "$value"},
				"count":  bson.M{"$sum": 1},
			}},
		}
		pipeline = append(pipeline, projectStage, groupStage)

	default:
		// Default projection and grouping
		projectStage := bson.D{
			{Key: "$project", Value: bson.M{
				"value":        "$answers.value",
				"submitted_at": 1,
			}},
		}
		groupStage := bson.D{
			{Key: "$group", Value: bson.M{
				"_id":    nil,
				"values": bson.M{"$push": "$value"},
				"count":  bson.M{"$sum": 1},
			}},
		}
		pipeline = append(pipeline, projectStage, groupStage)
	}

	return pipeline
}

// processSingleMultipleChoiceResults processes results for single-choice or multiple-choice questions
func (r *AnalyticsRepository) processSingleMultipleChoiceResults(ctx context.Context, cursor *mongo.Cursor, question model.SurveyQuestion, result *model.QuestionResult, totalResponses int) (*model.QuestionResult, error) {
	optionCounts := make(map[string]model.OptionCount)

	// Initialize option counts with zero for all options
	for _, option := range question.Options {
		optionCounts[option.ID] = model.OptionCount{
			OptionID: option.ID,
			Text:     option.Text,
			Count:    0,
			Percent:  0,
		}
	}

	// Count total answers for this question
	totalAnswers := 0

	// Process aggregation results
	for cursor.Next(ctx) {
		var document struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}

		if err := cursor.Decode(&document); err != nil {
			return nil, fmt.Errorf("failed to decode document: %w", err)
		}

		// Update option count if it exists
		if option, exists := optionCounts[document.ID]; exists {
			option.Count = document.Count
			optionCounts[document.ID] = option
			totalAnswers += document.Count
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	// Calculate percentages
	if totalAnswers > 0 {
		for id, count := range optionCounts {
			count.Percent = float64(count.Count) / float64(totalAnswers) * 100
			optionCounts[id] = count
		}
	}

	// Update result
	result.TotalAnswers = totalAnswers
	result.Analytics.OptionCounts = optionCounts
	result.ResponseData = optionCounts

	return result, nil
}

// processScaleResults processes results for scale questions
func (r *AnalyticsRepository) processScaleResults(ctx context.Context, cursor *mongo.Cursor, result *model.QuestionResult) (*model.QuestionResult, error) {
	if cursor.Next(ctx) {
		var document struct {
			Values []float64 `bson:"values"`
			Count  int       `bson:"count"`
			Avg    float64   `bson:"avg"`
			Min    float64   `bson:"min"`
			Max    float64   `bson:"max"`
		}

		if err := cursor.Decode(&document); err != nil {
			return nil, fmt.Errorf("failed to decode document: %w", err)
		}

		// Calculate median
		median := calculateMedian(document.Values)

		// Update result
		result.TotalAnswers = document.Count
		result.Analytics.Average = &document.Avg
		result.Analytics.Median = &median
		result.Analytics.Min = &document.Min
		result.Analytics.Max = &document.Max
		result.ResponseData = map[string]interface{}{
			"average": document.Avg,
			"median":  median,
			"min":     document.Min,
			"max":     document.Max,
			"values":  document.Values,
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return result, nil
}

// processOpenTextResults processes results for open-text questions
func (r *AnalyticsRepository) processOpenTextResults(ctx context.Context, cursor *mongo.Cursor, result *model.QuestionResult) (*model.QuestionResult, error) {
	if cursor.Next(ctx) {
		var document struct {
			Values []string `bson:"values"`
			Count  int      `bson:"count"`
		}

		if err := cursor.Decode(&document); err != nil {
			return nil, fmt.Errorf("failed to decode document: %w", err)
		}

		// Calculate word frequency (optional)
		wordFrequency := calculateWordFrequency(document.Values)

		// Update result
		result.TotalAnswers = document.Count
		result.Analytics.WordFrequency = wordFrequency
		result.ResponseData = document.Values
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return result, nil
}

// processMatrixResults processes results for matrix questions
func (r *AnalyticsRepository) processMatrixResults(ctx context.Context, cursor *mongo.Cursor, question model.SurveyQuestion, result *model.QuestionResult) (*model.QuestionResult, error) {
	// Matrix questions are complex and would need a custom implementation
	// For now, we'll just collect the raw values
	var allValues []interface{}
	var count int

	if cursor.Next(ctx) {
		var document struct {
			Values []interface{} `bson:"values"`
			Count  int           `bson:"count"`
		}

		if err := cursor.Decode(&document); err != nil {
			return nil, fmt.Errorf("failed to decode document: %w", err)
		}

		allValues = document.Values
		count = document.Count
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	// Update result
	result.TotalAnswers = count
	result.ResponseData = allValues

	return result, nil
}

// processDefaultResults processes results for unknown question types
func (r *AnalyticsRepository) processDefaultResults(ctx context.Context, cursor *mongo.Cursor, result *model.QuestionResult) (*model.QuestionResult, error) {
	if cursor.Next(ctx) {
		var document struct {
			Values []interface{} `bson:"values"`
			Count  int           `bson:"count"`
		}

		if err := cursor.Decode(&document); err != nil {
			return nil, fmt.Errorf("failed to decode document: %w", err)
		}

		// Update result
		result.TotalAnswers = document.Count
		result.ResponseData = document.Values
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return result, nil
}

// GetIndividualResponses retrieves individual responses for a survey with pagination
func (r *AnalyticsRepository) GetIndividualResponses(ctx context.Context, filter model.IndividualResponsesFilter) (*model.IndividualResponsesResult, error) {
	// Convert survey ID from string to ObjectID
	surveyID, err := primitive.ObjectIDFromHex(filter.SurveyID)
	if err != nil {
		return nil, fmt.Errorf("invalid survey ID: %w", err)
	}

	// Create filter for MongoDB
	mongoFilter := bson.M{"survey_id": surveyID}

	// Add date filters if provided
	if filter.StartDate != nil {
		mongoFilter["submitted_at"] = bson.M{"$gte": filter.StartDate}
	}
	if filter.EndDate != nil {
		if _, hasSubmittedAt := mongoFilter["submitted_at"]; hasSubmittedAt {
			mongoFilter["submitted_at"].(bson.M)["$lte"] = filter.EndDate
		} else {
			mongoFilter["submitted_at"] = bson.M{"$lte": filter.EndDate}
		}
	}

	// Count total documents
	totalCount, err := r.responsesColl.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to count responses: %w", err)
	}

	// Calculate pagination
	if filter.Limit <= 0 {
		filter.Limit = 10 // Default limit
	}
	if filter.Page <= 0 {
		filter.Page = 1 // Default page
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(filter.Limit)))
	skip := (filter.Page - 1) * filter.Limit

	// Fetch paginated results
	findOptions := options.Find().
		SetSort(bson.D{{Key: "submitted_at", Value: -1}}). // Sort by submission date, newest first
		SetSkip(int64(skip)).
		SetLimit(int64(filter.Limit))

	cursor, err := r.responsesColl.Find(ctx, mongoFilter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch responses: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode responses
	var responses []model.SurveyResponse
	if err := cursor.All(ctx, &responses); err != nil {
		return nil, fmt.Errorf("failed to decode responses: %w", err)
	}

	// Create result
	result := &model.IndividualResponsesResult{
		SurveyID:     filter.SurveyID,
		TotalCount:   int(totalCount),
		Responses:    responses,
		CurrentPage:  filter.Page,
		TotalPages:   totalPages,
		ItemsPerPage: filter.Limit,
	}

	return result, nil
}

// Close closes the MongoDB connection
func (r *AnalyticsRepository) Close(ctx context.Context) error {
	return r.client.Disconnect(ctx)
}

// Helper functions

// calculateMedian calculates the median value from a slice of numbers
func calculateMedian(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sort.Float64s(values)

	if len(values)%2 == 0 {
		return (values[len(values)/2-1] + values[len(values)/2]) / 2
	}

	return values[len(values)/2]
}

// calculateWordFrequency calculates word frequency from text responses
func calculateWordFrequency(texts []string) map[string]int {
	frequency := make(map[string]int)

	for _, text := range texts {
		words := strings.Fields(strings.ToLower(text))
		for _, word := range words {
			// Clean word (remove punctuation, etc.) - implement as needed
			cleanWord := strings.Trim(word, ".,!?:;\"'()[]{}=+*&^%$#@~`<>/\\|")
			if len(cleanWord) > 2 { // Ignore very short words
				frequency[cleanWord]++
			}
		}
	}

	return frequency
}
