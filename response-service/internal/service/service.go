package service

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/survey-app/response-service/internal/contextkeys"
	"github.com/survey-app/response-service/internal/models"
	"github.com/survey-app/response-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ResponseServiceInterface defines methods for response-related business logic
type ResponseServiceInterface interface {
	SubmitResponse(ctx context.Context, req *models.CreateResponseRequest) error
	GetSurveyResponses(ctx context.Context, surveyID int) ([]*models.Response, error)
	GetSurveyAnalytics(ctx context.Context, surveyID int) (*models.SurveyAnalyticsResponse, error)
	ExportSurveyResponsesCSV(ctx context.Context, surveyID int) (csvData string, filename string, err error)
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
	log.Printf("[SERVICE_INFO] getSurveyDetails: Calling Survey Service at URL: %s", surveyURL)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", surveyURL, nil)
	if err != nil {
		log.Printf("[SERVICE_ERROR] getSurveyDetails: Failed to create request to survey-service for SurveyID %d: %v", surveyID, err)
		return nil, fmt.Errorf("failed to create request to survey-service: %w", err)
	}

	// Forward Authorization header if present in context
	if authHeaderVal := ctx.Value(contextkeys.AuthorizationHeaderKey); authHeaderVal != nil {
		if authHeader, ok := authHeaderVal.(string); ok && authHeader != "" {
			httpReq.Header.Set("Authorization", authHeader)
			log.Printf("[SERVICE_INFO] getSurveyDetails: Forwarding Authorization header to survey-service.")
		}
	}

	// Propagate X-User-ID and X-User-Roles from context if they exist
	if userIDVal := ctx.Value(contextkeys.UserIDKey); userIDVal != nil {
		if userID, ok := userIDVal.(int); ok {
			httpReq.Header.Set("X-User-ID", strconv.Itoa(userID))
			log.Printf("[SERVICE_INFO] getSurveyDetails: Forwarding X-User-ID: %d to survey-service", userID)
		}
	}
	if userRolesVal := ctx.Value(contextkeys.UserRolesKey); userRolesVal != nil {
		if roles, ok := userRolesVal.([]string); ok {
			rolesStr := fmt.Sprintf("%v", roles) // Produces "[role1 role2 ...]"
			httpReq.Header.Set("X-User-Roles", rolesStr)
			log.Printf("[SERVICE_INFO] getSurveyDetails: Forwarding X-User-Roles: %s to survey-service", rolesStr)
		}
	}

	httpResp, err := s.httpClient.Do(httpReq)
	if err != nil {
		log.Printf("[SERVICE_ERROR] getSurveyDetails: Failed to call survey-service for SurveyID %d: %v", surveyID, err)
		return nil, fmt.Errorf("failed to call survey-service: %w", err)
	}
	defer httpResp.Body.Close()

	log.Printf("[SERVICE_INFO] getSurveyDetails: Response status from survey-service for SurveyID %d: %d", surveyID, httpResp.StatusCode)
	if httpResp.StatusCode == http.StatusNotFound {
		return nil, errors.New("survey not found in survey-service")
	}
	if httpResp.StatusCode != http.StatusOK {
		// Consider logging the body for more context on non-OK responses
		return nil, fmt.Errorf("survey-service returned status %d", httpResp.StatusCode)
	}

	var surveyDetails models.SurveyDetailsFromService
	if err := json.NewDecoder(httpResp.Body).Decode(&surveyDetails); err != nil {
		log.Printf("[SERVICE_ERROR] getSurveyDetails: Failed to decode response from survey-service for SurveyID %d: %v", surveyID, err)
		return nil, fmt.Errorf("failed to decode response from survey-service: %w", err)
	}
	return &surveyDetails, nil
}

// SubmitResponse handles the business logic for submitting a new survey response
func (s *ResponseService) SubmitResponse(ctx context.Context, req *models.CreateResponseRequest) error {
	log.Printf("[SERVICE_INFO] SubmitResponse: Attempting to submit response for SurveyID: %d, UserID: %v", req.SurveyID, req.UserID)
	log.Printf("[SERVICE_INFO] SubmitResponse: Request Answers: %+v", req.Answers)

	surveyDetails, err := s.getSurveyDetails(ctx, req.SurveyID)
	if err != nil {
		log.Printf("[SERVICE_ERROR] SubmitResponse: Failed to get survey details for SurveyID %d: %v", req.SurveyID, err)
		return fmt.Errorf("failed to retrieve survey details (ID: %d): %w", req.SurveyID, err)
	}
	log.Printf("[SERVICE_INFO] SubmitResponse: Successfully fetched survey details for SurveyID %d: %+v", req.SurveyID, surveyDetails)

	if !surveyDetails.IsActive {
		log.Printf("[SERVICE_WARN] SubmitResponse: SurveyID %d is not active. Aborting submission.", req.SurveyID)
		return errors.New("survey is not active and cannot accept new responses")
	}

	// TODO: Validate req.Answers against surveyDetails.Questions
	// - Check if question IDs in answers are valid for the survey.
	// - Check if answer values are consistent with question types (e.g., selected option ID is valid for single_choice).
	log.Printf("[SERVICE_INFO] SubmitResponse: SurveyID %d is active. Proceeding with response creation.", req.SurveyID)

	response := &models.Response{
		SurveyID: req.SurveyID,
		UserID:   req.UserID, // UserID is now reliably set by the handler from X-User-ID or original req
		Answers:  req.Answers,
		// SubmittedAt will be set by the repository
	}

	log.Printf("[SERVICE_INFO] SubmitResponse: Attempting to create response in repository for SurveyID: %d", req.SurveyID)
	err = s.repo.CreateResponse(ctx, response)
	if err != nil {
		log.Printf("[SERVICE_ERROR] SubmitResponse: Failed to create response in repository for SurveyID %d: %v", req.SurveyID, err)
		return fmt.Errorf("failed to save response to database (SurveyID: %d): %w", req.SurveyID, err)
	}

	log.Printf("[SERVICE_INFO] SubmitResponse: Successfully created response for SurveyID: %d", req.SurveyID)
	return nil
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

// ExportSurveyResponsesCSV generates a CSV string of all responses for a given survey.
// It fetches survey details for question text (headers) and all responses.
func (s *ResponseService) ExportSurveyResponsesCSV(ctx context.Context, surveyID int) (csvDataStr string, filename string, err error) {
	log.Printf("[SERVICE_INFO] ExportSurveyResponsesCSV: Starting export for SurveyID %d", surveyID)

	// 1. Fetch Survey Details (for question texts as headers)
	surveyDetails, err := s.getSurveyDetails(ctx, surveyID)
	if err != nil {
		log.Printf("[SERVICE_ERROR] ExportSurveyResponsesCSV: Failed to get survey details for SurveyID %d: %v", surveyID, err)
		return "", "", fmt.Errorf("failed to retrieve survey details (ID: %d): %w", surveyID, err)
	}
	log.Printf("[SERVICE_INFO] ExportSurveyResponsesCSV: Successfully fetched survey details for SurveyID %d. Number of questions: %d", surveyID, len(surveyDetails.Questions))

	// 2. Fetch All Responses for this Survey
	responses, err := s.GetSurveyResponses(ctx, surveyID) // This already calls s.repo.GetResponsesBySurveyID
	if err != nil {
		log.Printf("[SERVICE_ERROR] ExportSurveyResponsesCSV: Failed to get responses for SurveyID %d: %v", surveyID, err)
		return "", "", fmt.Errorf("failed to retrieve responses (ID: %d): %w", surveyID, err)
	}
	log.Printf("[SERVICE_INFO] ExportSurveyResponsesCSV: Successfully fetched %d responses for SurveyID %d", len(responses), surveyID)

	// 3. Prepare CSV Data
	var csvBuffer strings.Builder
	csvWriter := csv.NewWriter(&csvBuffer)

	// 3a. Write Headers
	headers := []string{"ResponseID", "SubmittedAt", "UserID"}
	questionIDToHeaderIndex := make(map[int]int) // Map question ID to its column index in the CSV
	questionIDToType := make(map[int]string)     // Map question ID to its type for answer formatting

	for _, q := range surveyDetails.Questions {
		headers = append(headers, q.Text) // Use question text as header
		questionIDToHeaderIndex[q.ID] = len(headers) - 1
		questionIDToType[q.ID] = q.Type
	}

	if err := csvWriter.Write(headers); err != nil {
		log.Printf("[SERVICE_ERROR] ExportSurveyResponsesCSV: Failed to write CSV headers for SurveyID %d: %v", surveyID, err)
		return "", "", fmt.Errorf("failed to write CSV headers: %w", err)
	}

	// 3b. Write Response Rows
	for _, resp := range responses {
		row := make([]string, len(headers)) // Initialize row with empty strings for all columns

		// Standard columns
		row[0] = resp.ID.Hex() // Convert ObjectID to string
		row[1] = resp.SubmittedAt.Format(time.RFC3339)
		if resp.UserID != nil {
			row[2] = strconv.Itoa(*resp.UserID)
		} else {
			row[2] = "Anonymous"
		}

		// Answer columns - map answers to the correct question column
		for _, ans := range resp.Answers {
			if headerIndex, ok := questionIDToHeaderIndex[ans.QuestionID]; ok {
				// Format answer value based on its type and question type
				var formattedValue string
				if ans.Value == nil {
					formattedValue = ""
				} else {
					qType := questionIDToType[ans.QuestionID]
					switch qType {
					case "checkbox": // Checkbox answers are []interface{} or []string
						if valSlice, ok := ans.Value.([]interface{}); ok {
							var strVals []string
							for _, item := range valSlice {
								strVals = append(strVals, fmt.Sprintf("%v", item))
							}
							formattedValue = strings.Join(strVals, "; ") // Join multiple checkbox values
						} else if valStrSlice, ok := ans.Value.([]string); ok {
							formattedValue = strings.Join(valStrSlice, "; ")
						} else {
							formattedValue = fmt.Sprintf("%v", ans.Value) // Fallback
						}
					default:
						formattedValue = fmt.Sprintf("%v", ans.Value)
					}
				}
				row[headerIndex] = formattedValue
			} else {
				log.Printf("[SERVICE_WARN] ExportSurveyResponsesCSV: Answer for QuestionID %d found in response %s, but this QuestionID is not in the survey's question list.", ans.QuestionID, resp.ID.Hex())
			}
		}
		if err := csvWriter.Write(row); err != nil {
			log.Printf("[SERVICE_ERROR] ExportSurveyResponsesCSV: Error writing CSV row for response %s, SurveyID %d: %v", resp.ID.Hex(), surveyID, err)
			// Potentially continue and try other rows, or return an error
		}
	}
	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		log.Printf("[SERVICE_ERROR] ExportSurveyResponsesCSV: CSV writer error for SurveyID %d: %v", surveyID, err)
		return "", "", fmt.Errorf("error during CSV generation: %w", err)
	}

	generatedFilename := fmt.Sprintf("survey_%d_responses_%s.csv", surveyID, time.Now().Format("20060102_150405"))
	log.Printf("[SERVICE_INFO] ExportSurveyResponsesCSV: Successfully generated CSV data for SurveyID %d. Filename: %s", surveyID, generatedFilename)

	return csvBuffer.String(), generatedFilename, nil
}
