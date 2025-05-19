package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/survey-app/survey-service/internal/models"
	"github.com/survey-app/survey-service/internal/repository"
)

// ContextKey is the type for context keys to avoid collisions.
// Exported to be used by handlers package.
type ContextKey string

const (
	// UserIDKey is the context key for the user's ID.
	UserIDKey ContextKey = "userID"
	// UserRolesKey is the context key for the user's roles.
	UserRolesKey ContextKey = "userRoles"
)

// Custom error types
var ErrForbidden = errors.New("forbidden")
var ErrNotFound = errors.New("not found")

// SurveyServiceInterface defines the interface for survey operations
type SurveyServiceInterface interface {
	// Survey operations
	CreateSurvey(ctx context.Context, survey *models.Survey, questions []models.QuestionUpdateRequest) (int, error)
	GetSurveys(ctx context.Context) ([]*models.Survey, error)
	GetSurvey(ctx context.Context, id int) (*models.Survey, error)
	UpdateSurvey(ctx context.Context, survey *models.Survey) error
	UpdateSurveyWithQuestions(ctx context.Context, survey *models.Survey, questions []models.QuestionUpdateRequest) error
	DeleteSurvey(ctx context.Context, id int) error
	UpdateSurveyStatus(ctx context.Context, id int, isActive bool) error

	// Question operations
	AddQuestion(ctx context.Context, req *models.CreateQuestionRequest) (int, error)
	UpdateQuestion(ctx context.Context, question *models.Question, options []*models.QuestionOption) error
	DeleteQuestion(ctx context.Context, id int) error
}

// SurveyService handles survey operations
type SurveyService struct {
	repo repository.SurveyRepositoryInterface
}

// Repository defines the interface for database operations
type Repository interface {
	// Survey operations
	CreateSurvey(ctx context.Context, survey *models.Survey) (int, error)
	GetSurvey(ctx context.Context, id int) (*models.Survey, error)
	GetSurveys(ctx context.Context, creatorID int) ([]*models.Survey, error)
	UpdateSurvey(ctx context.Context, survey *models.Survey) error
	DeleteSurvey(ctx context.Context, id int) error

	// Question operations
	CreateQuestion(ctx context.Context, question *models.Question) (int, error)
	GetQuestionsBySurveyID(ctx context.Context, surveyID int) ([]*models.Question, error)
	UpdateQuestion(ctx context.Context, question *models.Question) error
	DeleteQuestion(ctx context.Context, id int) error

	// Question option operations
	CreateQuestionOption(ctx context.Context, option *models.QuestionOption) (int, error)
	GetQuestionOptionsByQuestionID(ctx context.Context, questionID int) ([]*models.QuestionOption, error)
	DeleteQuestionOptions(ctx context.Context, questionID int) error
}

// NewSurveyService creates a new SurveyService
func NewSurveyService(repo repository.SurveyRepositoryInterface) *SurveyService {
	return &SurveyService{repo: repo}
}

// Helper function to extract user ID and roles from context
func getUserAndRolesFromContext(ctx context.Context) (userID int, roles []string, err error) {
	userIDVal := ctx.Value(UserIDKey)
	uid, ok := userIDVal.(int)
	if !ok {
		return 0, nil, errors.New("user ID not found or invalid type in context")
	}

	rolesVal := ctx.Value(UserRolesKey)
	rs, ok := rolesVal.([]string)
	if !ok {
		// If roles are not critical for a specific operation or can be empty,
		// this might return nil for roles instead of an error.
		// For authorization, missing roles might be an issue.
		return uid, nil, errors.New("user roles not found or invalid type in context")
	}
	return uid, rs, nil
}

// Helper function to check if a slice contains a string
func containsString(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// authorizeSurveyAccess checks if the user in context can perform an action on the survey.
// It returns the survey (if fetched and authorized), a boolean indicating if the user is an admin, and an error.
func (s *SurveyService) authorizeSurveyAccess(ctx context.Context, surveyID int) (survey *models.Survey, isUserAdmin bool, err error) {
	userID, roles, err := getUserAndRolesFromContext(ctx)
	if err != nil {
		log.Printf("[SVC_AUTH_ERROR] authorizeSurveyAccess: Error getting user/roles from context for surveyID %d: %v", surveyID, err)
		return nil, false, fmt.Errorf("authorization context error: %w", err)
	}
	log.Printf("[SVC_AUTH_DEBUG] authorizeSurveyAccess: For surveyID %d - Context UserID: %d, Roles: %v", surveyID, userID, roles)

	isUserAdmin = containsString(roles, "admin")
	log.Printf("[SVC_AUTH_DEBUG] authorizeSurveyAccess: For surveyID %d - Is context user admin? %t", surveyID, isUserAdmin)

	survey, err = s.repo.GetSurvey(ctx, surveyID)
	if err != nil {
		log.Printf("[SVC_ERROR] authorizeSurveyAccess: GetSurvey failed for surveyID %d: %v", surveyID, err)
		return nil, isUserAdmin, ErrNotFound
	}
	if survey == nil {
		log.Printf("[SVC_WARN] authorizeSurveyAccess: Survey ID %d not found by repo.GetSurvey, but no error returned from repo.", surveyID)
		return nil, isUserAdmin, ErrNotFound
	}
	log.Printf("[SVC_AUTH_DEBUG] authorizeSurveyAccess: For surveyID %d - Fetched survey.CreatorID: %d", surveyID, survey.CreatorID)

	if isUserAdmin {
		log.Printf("[SVC_AUTH_INFO] authorizeSurveyAccess: Access GRANTED for surveyID %d (User is ADMIN)", surveyID)
		return survey, true, nil // Admin has access
	}

	if survey.CreatorID == userID {
		log.Printf("[SVC_AUTH_INFO] authorizeSurveyAccess: Access GRANTED for surveyID %d (User is OWNER - survey.CreatorID %d == context.userID %d)", surveyID, survey.CreatorID, userID)
		return survey, false, nil // Owner has access
	}

	log.Printf("[SVC_AUTH_INFO] authorizeSurveyAccess: Access DENIED for surveyID %d (User is NOT ADMIN and NOT OWNER - survey.CreatorID %d != context.userID %d)", surveyID, survey.CreatorID, userID)
	return nil, false, ErrForbidden // Neither admin nor owner
}

// CreateSurvey creates a new survey and its questions/options if provided
func (s *SurveyService) CreateSurvey(ctx context.Context, survey *models.Survey, requestedQuestions []models.QuestionUpdateRequest) (int, error) {
	userID, _, err := getUserAndRolesFromContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("CreateSurvey: %w", err) // Error getting user from context
	}
	survey.CreatorID = userID // Set CreatorID from context

	log.Printf("[SVC_DEBUG] CreateSurvey CALLED for Survey Title: %s by UserID: %d", survey.Title, survey.CreatorID)
	log.Printf("[SVC_DEBUG] SurveyData: %+v", survey)
	log.Printf("[SVC_DEBUG] RequestedQuestions for new survey: %+v", requestedQuestions)

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		log.Printf("[SVC_ERROR] CreateSurvey: Failed to begin transaction: %v", err)
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	log.Printf("[SVC_DEBUG] CreateSurvey: Transaction BEGAN")

	var surveyID int // Declare surveyID here to be accessible in defer
	defer func() {
		if p := recover(); p != nil {
			log.Printf("[SVC_PANIC] CreateSurvey: Recovered panic. Rolling back. Panic: %v", p)
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			log.Printf("[SVC_ERROR] CreateSurvey: Error occurred. Rolling back. Error: %v", err)
			rbErr := tx.Rollback(ctx)
			if rbErr != nil {
				log.Printf("[SVC_ERROR] CreateSurvey: Transaction rollback FAILED: %v", rbErr)
			} else {
				log.Printf("[SVC_DEBUG] CreateSurvey: Transaction ROLLED BACK")
			}
		} else {
			log.Printf("[SVC_DEBUG] CreateSurvey: Attempting to COMMIT transaction")
			err = tx.Commit(ctx)
			if err != nil {
				log.Printf("[SVC_ERROR] CreateSurvey: Transaction commit FAILED: %v", err)
			} else {
				log.Printf("[SVC_DEBUG] CreateSurvey: Transaction COMMITTED")
			}
		}
	}()

	log.Printf("[SVC_DEBUG] CreateSurvey: Step 1: Creating survey entry for Title: %s, CreatorID: %d", survey.Title, survey.CreatorID)
	surveyID, err = s.repo.CreateSurveyTx(ctx, tx, survey)
	if err != nil {
		log.Printf("[SVC_ERROR] CreateSurvey: Step 1 FAILED: CreateSurveyTx: %v", err)
		return 0, fmt.Errorf("failed to create survey entry in transaction: %w", err)
	}
	survey.ID = surveyID
	log.Printf("[SVC_DEBUG] CreateSurvey: Step 1 SUCCESS: Survey entry created with ID: %d", surveyID)

	if len(requestedQuestions) > 0 {
		log.Printf("[SVC_DEBUG] CreateSurvey: Step 2: Processing %d requested questions for new Survey ID: %d", len(requestedQuestions), surveyID)
		for i, reqQuestion := range requestedQuestions {
			log.Printf("[SVC_DEBUG] CreateSurvey: Processing requested question #%d: %+v", i+1, reqQuestion)
			questionModel := &models.Question{
				SurveyID: surveyID,
				Text:     reqQuestion.Text,
				Type:     reqQuestion.Type,
				Required: reqQuestion.Required,
				OrderNum: i + 1,
			}

			log.Printf("[SVC_DEBUG] CreateSurvey: Creating NEW question in DB: %+v", questionModel)
			newQuestionID, errCreate := s.repo.CreateQuestionTx(ctx, tx, questionModel)
			if errCreate != nil {
				log.Printf("[SVC_ERROR] CreateSurvey: CreateQuestionTx FAILED for new question: %v", errCreate)
				err = fmt.Errorf("failed to create new question in transaction: %w", errCreate)
				return 0, err
			}
			questionModel.ID = newQuestionID
			log.Printf("[SVC_DEBUG] CreateSurvey: CreateQuestionTx SUCCESS. New Question ID: %d", newQuestionID)

			if len(reqQuestion.Options) > 0 {
				log.Printf("[SVC_DEBUG] CreateSurvey: %d options provided for New Question ID %d. Creating them.", len(reqQuestion.Options), newQuestionID)
				for optIdx, optText := range reqQuestion.Options {
					optionModel := &models.QuestionOption{
						QuestionID: newQuestionID,
						Text:       optText,
						OrderNum:   optIdx + 1,
					}
					log.Printf("[SVC_DEBUG] CreateSurvey: Creating option #%d for New Question ID %d: %+v", optIdx+1, newQuestionID, optionModel)
					_, errCreateOpt := s.repo.CreateQuestionOptionTx(ctx, tx, optionModel)
					if errCreateOpt != nil {
						log.Printf("[SVC_ERROR] CreateSurvey: CreateQuestionOptionTx FAILED for New Question ID %d, Option Text '%s': %v", newQuestionID, optText, errCreateOpt)
						err = fmt.Errorf("failed to create option for new question ID %d: %w", newQuestionID, errCreateOpt)
						return 0, err
					}
					log.Printf("[SVC_DEBUG] CreateSurvey: CreateQuestionOptionTx SUCCESS for option '%s', New Question ID %d", optText, newQuestionID)
				}
			} else {
				log.Printf("[SVC_DEBUG] CreateSurvey: No options provided for New Question ID %d.", newQuestionID)
			}
		}
		log.Printf("[SVC_DEBUG] CreateSurvey: Step 2 SUCCESS: All requested questions processed for new Survey ID: %d", surveyID)
	} else {
		log.Printf("[SVC_DEBUG] CreateSurvey: No questions provided in the request for new Survey ID: %d", surveyID)
	}

	log.Printf("[SVC_DEBUG] CreateSurvey function ENDING for new Survey ID: %d. Final 'err' before defer: %v", surveyID, err)
	return surveyID, err
}

// GetSurvey gets a survey by ID, checking for active status or ownership/admin rights
func (s *SurveyService) GetSurvey(ctx context.Context, id int) (*models.Survey, error) {
	// Call authorizeSurveyAccess. We don't need isUserAdmin directly in this function scope
	// as the authorization decision is handled by the error or returned survey.
	survey, _, err := s.authorizeSurveyAccess(ctx, id)
	if err != nil {
		// If error is Forbidden, but we want to allow public access to active surveys:
		if errors.Is(err, ErrForbidden) {
			// Fetch survey directly to check IsActive (potential re-fetch, consider optimizing)
			publicSurvey, publicErr := s.repo.GetSurvey(ctx, id)
			if publicErr != nil {
				return nil, publicErr // Error fetching for public check
			}
			if publicSurvey == nil {
				return nil, ErrNotFound
			}
			if publicSurvey.IsActive {
				return publicSurvey, nil // Publicly accessible active survey
			}
		}
		return nil, err // Original error (ErrForbidden if not active, ErrNotFound, or other)
	}
	// If authorizeSurveyAccess passed, survey is not nil and user is authorized (owner or admin)
	return survey, nil
}

// GetSurveys gets surveys based on user role (all for admin, own for user)
func (s *SurveyService) GetSurveys(ctx context.Context) ([]*models.Survey, error) {
	userID, roles, err := getUserAndRolesFromContext(ctx)
	if err != nil {
		// Log the error but proceed to fetch all surveys for any authenticated user if context is partially available or for basic view
		// Depending on strictness, could return error here.
		// For the new requirement "all users see all surveys", we might not strictly need userID/roles here if we always fetch all.
		log.Printf("[SERVICE_WARN] GetSurveys: Error getting full user context: %v. Proceeding to fetch all surveys.", err)
		// Fallback to fetching all surveys, frontend will handle edit/delete visibility.
		// If err meant no user context at all, this might be an issue if some auth is still expected.
		// However, API gateway already ensures user is authenticated with jwtAuthMiddleware for GET /surveys
	}

	log.Printf("[SERVICE_INFO] GetSurveys: Called by UserID: %d, Roles: %v", userID, roles) // userID might be 0 if context error occurred

	// New requirement: All authenticated users see all surveys.
	// Admins effectively have this already. For regular users, change from GetSurveysByCreatorID to GetAllSurveys.
	// The authorization for editing/deleting is handled by authorizeSurveyAccess in those respective methods.
	return s.repo.GetAllSurveys(ctx)
}

// UpdateSurvey updates a survey - DEPRECATED in favor of UpdateSurveyWithQuestions?
// If still used, it needs authorization.
func (s *SurveyService) UpdateSurvey(ctx context.Context, survey *models.Survey) error {
	// This method would need to fetch the existing survey to check CreatorID if survey.CreatorID isn't reliable
	// or ensure survey.CreatorID is set correctly by the caller based on existing record.
	// For now, assuming UpdateSurveyWithQuestions is the primary update path.
	_, _, err := s.authorizeSurveyAccess(ctx, survey.ID)
	if err != nil {
		return err
	}
	return s.repo.UpdateSurvey(ctx, survey)
}

// UpdateSurveyWithQuestions updates a survey and its questions/options
func (s *SurveyService) UpdateSurveyWithQuestions(ctx context.Context, surveyToUpdate *models.Survey, requestedQuestions []models.QuestionUpdateRequest) error {
	existingSurvey, isUserAdmin, err := s.authorizeSurveyAccess(ctx, surveyToUpdate.ID)
	if err != nil {
		return err // Handles ErrForbidden, ErrNotFound, or other errors
	}
	// User is authorized (owner or admin)
	surveyToUpdate.CreatorID = existingSurvey.CreatorID // Ensure CreatorID is not changed from original

	log.Printf("[SVC_DEBUG] UpdateSurveyWithQuestions CALLED for Survey ID: %d by UserID: %d (Admin: %t)", surveyToUpdate.ID, existingSurvey.CreatorID, isUserAdmin)
	log.Printf("[SVC_DEBUG] SurveyToUpdate (basic fields): %+v", surveyToUpdate)
	log.Printf("[SVC_DEBUG] RequestedQuestions for update: %+v", requestedQuestions)

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		log.Printf("[SVC_ERROR] Failed to begin transaction for Survey ID %d: %v", surveyToUpdate.ID, err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	log.Printf("[SVC_DEBUG] Transaction BEGAN for Survey ID: %d", surveyToUpdate.ID)

	defer func() {
		if p := recover(); p != nil {
			log.Printf("[SVC_PANIC] Recovered panic during UpdateSurveyWithQuestions for Survey ID %d. Rolling back. Panic: %v", surveyToUpdate.ID, p)
			_ = tx.Rollback(ctx)
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			log.Printf("[SVC_ERROR] Error occurred during UpdateSurveyWithQuestions for Survey ID %d. Rolling back. Error: %v", surveyToUpdate.ID, err)
			rbErr := tx.Rollback(ctx)
			if rbErr != nil {
				log.Printf("[SVC_ERROR] Transaction rollback FAILED for Survey ID %d: %v", surveyToUpdate.ID, rbErr)
			} else {
				log.Printf("[SVC_DEBUG] Transaction ROLLED BACK for Survey ID: %d", surveyToUpdate.ID)
			}
		} else {
			log.Printf("[SVC_DEBUG] Attempting to COMMIT transaction for Survey ID: %d", surveyToUpdate.ID)
			err = tx.Commit(ctx)
			if err != nil {
				log.Printf("[SVC_ERROR] Transaction commit FAILED for Survey ID %d: %v", surveyToUpdate.ID, err)
			} else {
				log.Printf("[SVC_DEBUG] Transaction COMMITTED for Survey ID: %d", surveyToUpdate.ID)
			}
		}
	}()

	// 1. Update the basic survey entry (Title, Description, IsActive, Dates)
	log.Printf("[SVC_DEBUG] UpdateSurveyWithQuestions: Step 1: Updating survey entry for ID: %d", surveyToUpdate.ID)
	err = s.repo.UpdateSurveyTx(ctx, tx, surveyToUpdate) // Assuming UpdateSurveyTx exists
	if err != nil {
		log.Printf("[SVC_ERROR] UpdateSurveyWithQuestions: Step 1 FAILED: UpdateSurveyTx for ID %d: %v", surveyToUpdate.ID, err)
		return fmt.Errorf("failed to update survey entry in transaction: %w", err)
	}
	log.Printf("[SVC_DEBUG] UpdateSurveyWithQuestions: Step 1 SUCCESS: Survey entry updated for ID: %d", surveyToUpdate.ID)

	// 2. Process questions: diff between existing and requested to Add, Update, or Delete questions
	// This logic is complex: fetch existing questions, compare with reqQuestions, then act.
	// For simplicity in this step, the example below assumes full replacement or a more sophisticated repo method.
	// A robust implementation would involve: s.repo.GetQuestionsBySurveyIDTx(...), then diffing.
	// Then s.repo.DeleteQuestionTx, s.repo.UpdateQuestionTx, s.repo.CreateQuestionTx.
	// For now, let's assume a placeholder for this complex diff logic or a simpler approach:
	// Example: Delete all existing questions and recreate from request (simplistic, but shows transaction use)
	log.Printf("[SVC_DEBUG] UpdateSurveyWithQuestions: Step 2a: Deleting existing questions for Survey ID %d before re-adding (simplistic approach).", surveyToUpdate.ID)
	err = s.repo.DeleteQuestionsBySurveyIDTx(ctx, tx, surveyToUpdate.ID)
	if err != nil {
		log.Printf("[SVC_ERROR] UpdateSurveyWithQuestions: Failed to delete existing questions for Survey ID %d: %v", surveyToUpdate.ID, err)
		return fmt.Errorf("failed to delete existing questions: %w", err)
	}
	log.Printf("[SVC_DEBUG] UpdateSurveyWithQuestions: Step 2a SUCCESS: Existing questions deleted for Survey ID %d.", surveyToUpdate.ID)

	if len(requestedQuestions) > 0 {
		log.Printf("[SVC_DEBUG] UpdateSurveyWithQuestions: Step 2b: Processing %d requested questions for Survey ID: %d", len(requestedQuestions), surveyToUpdate.ID)
		for i, reqQuestion := range requestedQuestions {
			questionModel := &models.Question{
				SurveyID: surveyToUpdate.ID,
				Text:     reqQuestion.Text,
				Type:     reqQuestion.Type,
				Required: reqQuestion.Required,
				OrderNum: i + 1, // Or reqQuestion.OrderNum
			}
			// If reqQuestion.ID is present and non-zero, it could imply an update to an existing question if not deleting all.
			// In our simplistic delete-all approach, all are new.
			log.Printf("[SVC_DEBUG] UpdateSurveyWithQuestions: Creating question in DB: %+v", questionModel)
			newQuestionID, errCreate := s.repo.CreateQuestionTx(ctx, tx, questionModel)
			if errCreate != nil {
				log.Printf("[SVC_ERROR] UpdateSurveyWithQuestions: CreateQuestionTx FAILED: %v", errCreate)
				err = fmt.Errorf("failed to create question in transaction: %w", errCreate)
				return err
			}
			log.Printf("[SVC_DEBUG] UpdateSurveyWithQuestions: CreateQuestionTx SUCCESS. Question ID: %d", newQuestionID)

			if len(reqQuestion.Options) > 0 {
				for optIdx, optText := range reqQuestion.Options {
					optionModel := &models.QuestionOption{
						QuestionID: newQuestionID,
						Text:       optText,
						OrderNum:   optIdx + 1, // Or option.OrderNum
					}
					_, errCreateOpt := s.repo.CreateQuestionOptionTx(ctx, tx, optionModel)
					if errCreateOpt != nil {
						log.Printf("[SVC_ERROR] UpdateSurveyWithQuestions: CreateQuestionOptionTx FAILED: %v", errCreateOpt)
						err = fmt.Errorf("failed to create option: %w", errCreateOpt)
						return err
					}
				}
			}
		}
		log.Printf("[SVC_DEBUG] UpdateSurveyWithQuestions: Step 2b SUCCESS: Requested questions processed for Survey ID: %d", surveyToUpdate.ID)
	}

	log.Printf("[SVC_DEBUG] UpdateSurveyWithQuestions function ENDING for Survey ID: %d. Final 'err' before defer: %v", surveyToUpdate.ID, err)
	return err // err will be handled by defer for commit/rollback
}

// DeleteSurvey deletes a survey
func (s *SurveyService) DeleteSurvey(ctx context.Context, id int) error {
	_, _, err := s.authorizeSurveyAccess(ctx, id)
	if err != nil {
		return err
	}
	// User is authorized (owner or admin)
	// Repository needs to handle cascading deletes of questions/options if DB doesn't via FK constraints
	return s.repo.DeleteSurvey(ctx, id)
}

// AddQuestion adds a question to an existing survey
func (s *SurveyService) AddQuestion(ctx context.Context, req *models.CreateQuestionRequest) (int, error) {
	// Authorize access to the survey first
	_, _, err := s.authorizeSurveyAccess(ctx, req.SurveyID)
	if err != nil {
		return 0, fmt.Errorf("AddQuestion: not authorized for survey %d: %w", req.SurveyID, err)
	}
	// User is authorized (owner or admin) to modify this survey

	question := &models.Question{
		SurveyID: req.SurveyID,
		Text:     req.Text,
		Type:     req.Type,
		Required: req.Required,
		OrderNum: req.OrderNum,
	}

	// Transaction for creating question and its options
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	var questionID int
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	questionID, err = s.repo.CreateQuestionTx(ctx, tx, question)
	if err != nil {
		return 0, fmt.Errorf("failed to create question: %w", err)
	}

	for i, optReq := range req.Options {
		option := &models.QuestionOption{
			QuestionID: questionID,
			Text:       optReq.Text,
			OrderNum:   i + 1, // Or optReq.OrderNum
		}
		_, err = s.repo.CreateQuestionOptionTx(ctx, tx, option)
		if err != nil {
			return 0, fmt.Errorf("failed to create question option: %w", err)
		}
	}

	return questionID, err
}

// UpdateSurveyStatus updates a survey's active status
func (s *SurveyService) UpdateSurveyStatus(ctx context.Context, id int, isActive bool) error {
	_, _, err := s.authorizeSurveyAccess(ctx, id)
	if err != nil {
		return err
	}
	// User is authorized (owner or admin)
	return s.repo.UpdateSurveyStatus(ctx, id, isActive) // Repo method needs to exist
}

// UpdateQuestion updates a question and its options
// This also needs authorization at the survey level
func (s *SurveyService) UpdateQuestion(ctx context.Context, question *models.Question, options []*models.QuestionOption) error {
	if question == nil {
		return errors.New("question data cannot be nil")
	}
	_, _, err := s.authorizeSurveyAccess(ctx, question.SurveyID)
	if err != nil {
		return fmt.Errorf("UpdateQuestion: not authorized for survey %d: %w", question.SurveyID, err)
	}
	// User is authorized (owner or admin) to modify this survey's questions

	// Transaction logic would be similar to AddQuestion or UpdateSurveyWithQuestions
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	err = s.repo.UpdateQuestionTx(ctx, tx, question)
	if err != nil {
		return fmt.Errorf("failed to update question: %w", err)
	}

	// Simplistic: delete existing options and recreate
	err = s.repo.DeleteQuestionOptionsTx(ctx, tx, question.ID)
	if err != nil {
		return fmt.Errorf("failed to delete existing options: %w", err)
	}

	for _, opt := range options {
		opt.QuestionID = question.ID
		_, err = s.repo.CreateQuestionOptionTx(ctx, tx, opt)
		if err != nil {
			return fmt.Errorf("failed to create option: %w", err)
		}
	}
	return err
}

// DeleteQuestion deletes a question
// This also needs authorization at the survey level
func (s *SurveyService) DeleteQuestion(ctx context.Context, id int) error {
	// Need to get question first to find its surveyID for authorization
	question, err := s.repo.GetQuestionByID(ctx, id) // Assuming GetQuestionByID exists
	if err != nil {
		return fmt.Errorf("failed to get question %d: %w", id, err)
	}
	if question == nil {
		return ErrNotFound
	}

	_, _, err = s.authorizeSurveyAccess(ctx, question.SurveyID)
	if err != nil {
		return fmt.Errorf("DeleteQuestion: not authorized for survey %d: %w", question.SurveyID, err)
	}
	// User is authorized (owner or admin)
	return s.repo.DeleteQuestion(ctx, id)
}
