package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/survey-app/response-service/internal/models"
	"github.com/survey-app/response-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ResponseServiceInterface defines methods for response-related business logic
type ResponseServiceInterface interface {
	SubmitResponse(ctx context.Context, req *models.CreateResponseRequest) error
	GetSurveyResponses(ctx context.Context, surveyID int) ([]*models.Response, error)
	GetSurveyAnalytics(ctx context.Context, surveyID int) (*models.SurveyAnalyticsResponse, error)
}

// ResponseService implements ResponseServiceInterface
type ResponseService struct {
	repo             repository.ResponseRepositoryInterface
	surveyServiceURL string
	httpClient       *http.Client
}

// NewResponseService creates a new ResponseService
func NewResponseService(repo repository.ResponseRepositoryInterface, surveyServiceURL string) *ResponseService {
	return &ResponseService{
		repo:             repo,
		surveyServiceURL: surveyServiceURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second, // Add a timeout for HTTP requests
		},
	}
}

// getSurveyDetails fetches full survey details from the survey-service
func (s *ResponseService) getSurveyDetails(ctx context.Context, surveyID int) (*models.SurveyDetailsFromService, error) {
	surveyURL := fmt.Sprintf("%s/api/v1/surveys/%d", s.surveyServiceURL, surveyID)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", surveyURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to survey-service: %w", err)
	}

	httpResp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call survey-service: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == http.StatusNotFound {
		return nil, errors.New("survey not found in survey-service")
	}
	if httpResp.StatusCode != http.StatusOK {
		// Consider logging the body for more context on non-OK responses
		return nil, fmt.Errorf("survey-service returned status %d", httpResp.StatusCode)
	}

	var surveyDetails models.SurveyDetailsFromService
	if err := json.NewDecoder(httpResp.Body).Decode(&surveyDetails); err != nil {
		return nil, fmt.Errorf("failed to decode response from survey-service: %w", err)
	}
	return &surveyDetails, nil
}

// SubmitResponse handles the business logic for submitting a new survey response
func (s *ResponseService) SubmitResponse(ctx context.Context, req *models.CreateResponseRequest) error {
	surveyDetails, err := s.getSurveyDetails(ctx, req.SurveyID)
	if err != nil {
		return err // Handles not found, service errors, decode errors
	}

	if !surveyDetails.IsActive {
		return errors.New("survey is not active and cannot accept new responses")
	}

	// TODO: Validate req.Answers against surveyDetails.Questions
	// - Check if question IDs in answers are valid for the survey.
	// - Check if answer values are consistent with question types (e.g., selected option ID is valid for single_choice).

	response := &models.Response{
		SurveyID: req.SurveyID,
		UserID:   req.UserID, // Assuming UserID is passed in CreateResponseRequest or obtained from context
		Answers:  req.Answers,
		// SubmittedAt will be set by the repository
	}

	return s.repo.CreateResponse(ctx, response)
}

// GetSurveyResponses retrieves all responses for a specific survey
func (s *ResponseService) GetSurveyResponses(ctx context.Context, surveyID int) ([]*models.Response, error) {
	// TODO: Add any transformation or additional logic if needed
	return s.repo.GetResponsesBySurveyID(ctx, surveyID)
}

// GetSurveyAnalytics retrieves and processes survey responses to generate analytics
func (s *ResponseService) GetSurveyAnalytics(ctx context.Context, surveyID int) (*models.SurveyAnalyticsResponse, error) {
	// 1. Fetch survey details
	surveyDetails, err := s.getSurveyDetails(ctx, surveyID)
	if err != nil {
		log.Printf("Error fetching survey details for analytics (surveyID: %d): %v", surveyID, err) // Keep error log
		return nil, fmt.Errorf("failed to get survey details for analytics: %w", err)
	}

	// 2. Fetch all responses
	responses, err := s.repo.GetResponsesBySurveyID(ctx, surveyID)
	if err != nil {
		log.Printf("Error fetching responses for analytics (surveyID: %d): %v", surveyID, err) // Keep error log
		return nil, fmt.Errorf("failed to get survey responses for analytics: %w", err)
	}

	totalResponses := len(responses)
	analyticsResp := &models.SurveyAnalyticsResponse{
		SurveyID:          surveyDetails.ID,
		SurveyTitle:       surveyDetails.Title,
		TotalResponses:    totalResponses,
		QuestionAnalytics: make([]models.QuestionAnalytics, 0, len(surveyDetails.Questions)),
	}

	optionTextToIDMap := make(map[int]map[string]int)
	for _, q := range surveyDetails.Questions {
		if q.Type == "single_choice" || q.Type == "multiple_choice" || q.Type == "dropdown" || q.Type == "checkbox" {
			optionTextToIDMap[q.ID] = make(map[string]int)
			for _, opt := range q.Options {
				optionTextToIDMap[q.ID][opt.Text] = opt.ID
			}
		}
	}

	if totalResponses == 0 {
		for _, q := range surveyDetails.Questions {
			qa := models.QuestionAnalytics{
				QuestionID:   q.ID,
				QuestionText: q.Text,
				QuestionType: q.Type,
			}
			if q.Type == "single_choice" || q.Type == "multiple_choice" || q.Type == "dropdown" || q.Type == "checkbox" || q.Type == "linear_scale" {
				qa.OptionsSummary = make([]models.OptionSummary, 0)
				if q.Type == "linear_scale" {
					for i := 1; i <= 5; i++ { // Assuming 1-5 scale
						val := i
						qa.OptionsSummary = append(qa.OptionsSummary, models.OptionSummary{OptionID: &val, OptionText: fmt.Sprintf("%d", val), Count: 0, Percentage: 0})
					}
				} else {
					for _, opt := range q.Options {
						optID := opt.ID
						qa.OptionsSummary = append(qa.OptionsSummary, models.OptionSummary{OptionID: &optID, OptionText: opt.Text, Count: 0, Percentage: 0})
					}
				}
			}
			analyticsResp.QuestionAnalytics = append(analyticsResp.QuestionAnalytics, qa)
		}
		return analyticsResp, nil
	}

	for _, q := range surveyDetails.Questions {
		qa := models.QuestionAnalytics{
			QuestionID:   q.ID,
			QuestionText: q.Text,
			QuestionType: q.Type,
		}
		actualRespondersToThisQuestion := 0

		switch q.Type {
		case "single_choice", "multiple_choice", "dropdown":
			optionCounts := make(map[int]int)
			for _, opt := range q.Options {
				optionCounts[opt.ID] = 0
			}
			for _, resp := range responses {
				foundAnswerToThisQuestionInResp := false
				for _, ans := range resp.Answers {
					if ans.QuestionID == q.ID {
						if selectedOptionText, ok := ans.Value.(string); ok {
							if optID, found := optionTextToIDMap[q.ID][selectedOptionText]; found {
								optionCounts[optID]++
								if !foundAnswerToThisQuestionInResp {
									actualRespondersToThisQuestion++
									foundAnswerToThisQuestionInResp = true
								}
							}
						}
						break
					}
				}
			}
			qa.OptionsSummary = make([]models.OptionSummary, 0, len(q.Options))
			for _, opt := range q.Options {
				count := optionCounts[opt.ID]
				percentage := 0.0
				if actualRespondersToThisQuestion > 0 {
					percentage = (float64(count) / float64(actualRespondersToThisQuestion)) * 100
				}
				optID := opt.ID
				qa.OptionsSummary = append(qa.OptionsSummary, models.OptionSummary{OptionID: &optID, OptionText: opt.Text, Count: count, Percentage: percentage})
			}

		case "checkbox":
			optionCounts := make(map[int]int)
			for _, opt := range q.Options {
				optionCounts[opt.ID] = 0
			}
			for _, resp := range responses {
				foundAnswerToThisQuestionInResp := false
				for _, ans := range resp.Answers {
					if ans.QuestionID == q.ID {
						if selectedOptionValues, ok := ans.Value.(primitive.A); ok {
							if len(selectedOptionValues) > 0 && !foundAnswerToThisQuestionInResp {
								actualRespondersToThisQuestion++
								foundAnswerToThisQuestionInResp = true
							}
							for _, valInterface := range selectedOptionValues {
								if optionValueStr, isString := valInterface.(string); isString {
									selectedOptionID, errAtoi := strconv.Atoi(optionValueStr)
									if errAtoi == nil {
										if _, knownOption := optionCounts[selectedOptionID]; knownOption {
											optionCounts[selectedOptionID]++
										}
									}
								}
							}
						}
						break
					}
				}
			}
			qa.OptionsSummary = make([]models.OptionSummary, 0, len(q.Options))
			for _, opt := range q.Options {
				count := optionCounts[opt.ID]
				percentage := 0.0
				if actualRespondersToThisQuestion > 0 {
					percentage = (float64(count) / float64(actualRespondersToThisQuestion)) * 100
				}
				optID := opt.ID
				qa.OptionsSummary = append(qa.OptionsSummary, models.OptionSummary{OptionID: &optID, OptionText: opt.Text, Count: count, Percentage: percentage})
			}

		case "linear_scale":
			valueCounts := make(map[int]int)
			minScale, maxScale := 1, 5 // Assuming 1-5 scale
			for i := minScale; i <= maxScale; i++ {
				valueCounts[i] = 0
			}
			for _, resp := range responses {
				foundAnswerToThisQuestionInResp := false
				for _, ans := range resp.Answers {
					if ans.QuestionID == q.ID {
						if selectedValueFloat, ok := ans.Value.(float64); ok {
							selectedValueInt := int(selectedValueFloat)
							if selectedValueInt >= minScale && selectedValueInt <= maxScale {
								valueCounts[selectedValueInt]++
								if !foundAnswerToThisQuestionInResp {
									actualRespondersToThisQuestion++
									foundAnswerToThisQuestionInResp = true
								}
							}
						}
						break
					}
				}
			}
			qa.OptionsSummary = make([]models.OptionSummary, 0, maxScale-minScale+1)
			for i := minScale; i <= maxScale; i++ {
				count := valueCounts[i]
				percentage := 0.0
				if actualRespondersToThisQuestion > 0 {
					percentage = (float64(count) / float64(actualRespondersToThisQuestion)) * 100
				}
				val := i
				qa.OptionsSummary = append(qa.OptionsSummary, models.OptionSummary{OptionID: &val, OptionText: fmt.Sprintf("%d", i), Count: count, Percentage: percentage})
			}

		case "text", "paragraph", "short_answer", "date":
			qa.TextResponses = make([]models.TextResponseData, 0)
			for _, resp := range responses {
				for _, ans := range resp.Answers {
					if ans.QuestionID == q.ID {
						if textValue, ok := ans.Value.(string); ok && textValue != "" {
							qa.TextResponses = append(qa.TextResponses, models.TextResponseData{Response: textValue})
						}
						break
					}
				}
			}

		default:
			// Handle unknown or non-analyzable question type
		}
		analyticsResp.QuestionAnalytics = append(analyticsResp.QuestionAnalytics, qa)
	}

	return analyticsResp, nil
}
