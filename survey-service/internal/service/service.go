package service

import (
	"context"
	"fmt"
	"log"

	"github.com/survey-app/survey-service/internal/models"
	"github.com/survey-app/survey-service/internal/repository"
)

// SurveyServiceInterface defines the interface for survey operations
type SurveyServiceInterface interface {
	// Survey operations
	CreateSurvey(ctx context.Context, survey *models.Survey, questions []models.QuestionUpdateRequest) (int, error)
	GetSurveys(ctx context.Context, creatorID int) ([]*models.Survey, error)
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

// CreateSurvey creates a new survey and its questions/options if provided
func (s *SurveyService) CreateSurvey(ctx context.Context, survey *models.Survey, requestedQuestions []models.QuestionUpdateRequest) (int, error) {
	log.Printf("[SVC_DEBUG] CreateSurvey CALLED for Survey Title: %s", survey.Title)
	log.Printf("[SVC_DEBUG] SurveyData: %+v", survey)
	log.Printf("[SVC_DEBUG] RequestedQuestions for new survey: %+v", requestedQuestions)

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		log.Printf("[SVC_ERROR] CreateSurvey: Failed to begin transaction: %v", err)
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	log.Printf("[SVC_DEBUG] CreateSurvey: Transaction BEGAN")

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

	// 1. Create the basic survey entry
	log.Printf("[SVC_DEBUG] CreateSurvey: Step 1: Creating survey entry for Title: %s", survey.Title)
	surveyID, err := s.repo.CreateSurveyTx(ctx, tx, survey) // Assuming CreateSurveyTx exists or is created
	if err != nil {
		log.Printf("[SVC_ERROR] CreateSurvey: Step 1 FAILED: CreateSurveyTx: %v", err)
		return 0, fmt.Errorf("failed to create survey entry in transaction: %w", err)
	}
	survey.ID = surveyID // Set the ID for the survey object for question linking
	log.Printf("[SVC_DEBUG] CreateSurvey: Step 1 SUCCESS: Survey entry created with ID: %d", surveyID)

	// 2. Process incoming questions if any
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
				return 0, err // Return 0 for surveyID as creation failed
			}
			questionModel.ID = newQuestionID
			log.Printf("[SVC_DEBUG] CreateSurvey: CreateQuestionTx SUCCESS. New Question ID: %d", newQuestionID)

			// Process options for this new question
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
						return 0, err // Return 0 for surveyID
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
	return surveyID, err // err will be handled by defer for commit/rollback
}

// GetSurvey gets a survey by ID
func (s *SurveyService) GetSurvey(ctx context.Context, id int) (*models.Survey, error) {
	return s.repo.GetSurvey(ctx, id)
}

// GetSurveys gets all surveys for a creator
func (s *SurveyService) GetSurveys(ctx context.Context, creatorID int) ([]*models.Survey, error) {
	return s.repo.GetSurveys(ctx, creatorID)
}

// UpdateSurvey updates a survey
func (s *SurveyService) UpdateSurvey(ctx context.Context, survey *models.Survey) error {
	return s.repo.UpdateSurvey(ctx, survey)
}

// UpdateSurveyWithQuestions updates a survey and its questions/options
func (s *SurveyService) UpdateSurveyWithQuestions(ctx context.Context, surveyToUpdate *models.Survey, requestedQuestions []models.QuestionUpdateRequest) error {
	log.Printf("[SVC_DEBUG] UpdateSurveyWithQuestions CALLED for Survey ID: %d", surveyToUpdate.ID)
	log.Printf("[SVC_DEBUG] SurveyToUpdate: %+v", surveyToUpdate)
	log.Printf("[SVC_DEBUG] RequestedQuestions: %+v", requestedQuestions)

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

	// 1. Update basic survey details
	log.Printf("[SVC_DEBUG] Step 1: Updating survey details for Survey ID: %d", surveyToUpdate.ID)
	if err = s.repo.UpdateSurveyTx(ctx, tx, surveyToUpdate); err != nil {
		log.Printf("[SVC_ERROR] Step 1 FAILED: UpdateSurveyTx for Survey ID %d: %v", surveyToUpdate.ID, err)
		return fmt.Errorf("failed to update survey details in transaction: %w", err)
	}
	log.Printf("[SVC_DEBUG] Step 1 SUCCESS: Survey details updated for Survey ID: %d", surveyToUpdate.ID)

	// 2. Get existing questions for this survey to compare
	log.Printf("[SVC_DEBUG] Step 2: Getting existing questions for Survey ID: %d", surveyToUpdate.ID)
	existingQuestions, err := s.repo.GetQuestionsBySurveyIDTx(ctx, tx, surveyToUpdate.ID)
	if err != nil {
		log.Printf("[SVC_ERROR] Step 2 FAILED: GetQuestionsBySurveyIDTx for Survey ID %d: %v", surveyToUpdate.ID, err)
		return fmt.Errorf("failed to get existing questions in transaction: %w", err)
	}
	log.Printf("[SVC_DEBUG] Step 2 SUCCESS: Found %d existing questions for Survey ID %d: %+v", len(existingQuestions), surveyToUpdate.ID, existingQuestions)

	existingQuestionMap := make(map[int]*models.Question)
	for _, q := range existingQuestions {
		existingQuestionMap[q.ID] = q
	}

	requestedQuestionIDs := make(map[int]bool)

	// 3. Process incoming questions: update existing or create new
	log.Printf("[SVC_DEBUG] Step 3: Processing %d requested questions for Survey ID: %d", len(requestedQuestions), surveyToUpdate.ID)
	for i, reqQuestion := range requestedQuestions {
		log.Printf("[SVC_DEBUG] Processing requested question #%d (Original ID: %v): %+v", i+1, reqQuestion.ID, reqQuestion)
		questionModel := &models.Question{
			SurveyID: surveyToUpdate.ID,
			Text:     reqQuestion.Text,
			Type:     reqQuestion.Type,
			Required: reqQuestion.Required,
			OrderNum: i + 1,
		}

		if reqQuestion.ID != nil && *reqQuestion.ID != 0 { // Existing question
			questionModel.ID = *reqQuestion.ID
			requestedQuestionIDs[questionModel.ID] = true
			log.Printf("[SVC_DEBUG] Updating EXISTING question (ID: %d): %+v", questionModel.ID, questionModel)

			if _, ok := existingQuestionMap[questionModel.ID]; !ok {
				log.Printf("[SVC_ERROR] Attempted to update non-existent question ID %d for Survey %d", questionModel.ID, surveyToUpdate.ID)
				err = fmt.Errorf("attempted to update non-existent question ID %d", questionModel.ID)
				return err
			}

			if err = s.repo.UpdateQuestionTx(ctx, tx, questionModel); err != nil {
				log.Printf("[SVC_ERROR] UpdateQuestionTx FAILED for Question ID %d (Survey %d): %v", questionModel.ID, surveyToUpdate.ID, err)
				return fmt.Errorf("failed to update question ID %d in transaction: %w", questionModel.ID, err)
			}
			log.Printf("[SVC_DEBUG] UpdateQuestionTx SUCCESS for Question ID %d", questionModel.ID)
		} else { // New question
			log.Printf("[SVC_DEBUG] Creating NEW question: %+v", questionModel)
			newQuestionID, errCreate := s.repo.CreateQuestionTx(ctx, tx, questionModel)
			if errCreate != nil {
				log.Printf("[SVC_ERROR] CreateQuestionTx FAILED for new question (Survey %d): %v", surveyToUpdate.ID, errCreate)
				err = fmt.Errorf("failed to create new question in transaction: %w", errCreate) // assign to outer err
				return err
			}
			questionModel.ID = newQuestionID
			log.Printf("[SVC_DEBUG] CreateQuestionTx SUCCESS. New Question ID: %d", newQuestionID)
		}

		log.Printf("[SVC_DEBUG] Processing options for Question ID %d (Text: '%s')", questionModel.ID, questionModel.Text)
		if err = s.repo.DeleteQuestionOptionsTx(ctx, tx, questionModel.ID); err != nil {
			log.Printf("[SVC_ERROR] DeleteQuestionOptionsTx FAILED for Question ID %d (Survey %d): %v", questionModel.ID, surveyToUpdate.ID, err)
			return fmt.Errorf("failed to delete old options for question ID %d: %w", questionModel.ID, err)
		}
		log.Printf("[SVC_DEBUG] DeleteQuestionOptionsTx SUCCESS for Question ID %d", questionModel.ID)

		if len(reqQuestion.Options) > 0 {
			log.Printf("[SVC_DEBUG] %d options provided for Question ID %d. Creating them.", len(reqQuestion.Options), questionModel.ID)
			for optIdx, optText := range reqQuestion.Options {
				optionModel := &models.QuestionOption{
					QuestionID: questionModel.ID,
					Text:       optText,
					OrderNum:   optIdx + 1,
				}
				log.Printf("[SVC_DEBUG] Creating option #%d for Question ID %d: %+v", optIdx+1, questionModel.ID, optionModel)
				_, errCreateOpt := s.repo.CreateQuestionOptionTx(ctx, tx, optionModel)
				if errCreateOpt != nil {
					log.Printf("[SVC_ERROR] CreateQuestionOptionTx FAILED for Question ID %d (Survey %d), Option Text '%s': %v", questionModel.ID, surveyToUpdate.ID, optText, errCreateOpt)
					err = fmt.Errorf("failed to create option for question ID %d: %w", questionModel.ID, errCreateOpt) // assign to outer err
					return err
				}
				log.Printf("[SVC_DEBUG] CreateQuestionOptionTx SUCCESS for option '%s', Question ID %d", optText, questionModel.ID)
			}
		} else {
			log.Printf("[SVC_DEBUG] No options provided for Question ID %d.", questionModel.ID)
		}
	}
	log.Printf("[SVC_DEBUG] Step 3 SUCCESS: All requested questions processed for Survey ID: %d", surveyToUpdate.ID)

	// 4. Delete questions that were in DB but not in the request
	log.Printf("[SVC_DEBUG] Step 4: Deleting questions not present in the request for Survey ID: %d", surveyToUpdate.ID)
	deletedCount := 0
	for _, existingQ := range existingQuestions {
		if _, foundInRequest := requestedQuestionIDs[existingQ.ID]; !foundInRequest {
			log.Printf("[SVC_DEBUG] Deleting Question ID %d (Text: '%s') as it's no longer in request for Survey %d.", existingQ.ID, existingQ.Text, surveyToUpdate.ID)
			if err = s.repo.DeleteQuestionTx(ctx, tx, existingQ.ID); err != nil {
				log.Printf("[SVC_ERROR] DeleteQuestionTx FAILED for Question ID %d (Survey %d): %v", existingQ.ID, surveyToUpdate.ID, err)
				return fmt.Errorf("failed to delete question ID %d: %w", existingQ.ID, err)
			}
			log.Printf("[SVC_DEBUG] DeleteQuestionTx SUCCESS for Question ID %d", existingQ.ID)
			deletedCount++
		}
	}
	log.Printf("[SVC_DEBUG] Step 4 SUCCESS: Deleted %d questions for Survey ID: %d", deletedCount, surveyToUpdate.ID)

	log.Printf("[SVC_DEBUG] UpdateSurveyWithQuestions function ENDING for Survey ID: %d. Final 'err' before defer: %v", surveyToUpdate.ID, err)
	return err
}

// DeleteSurvey deletes a survey
func (s *SurveyService) DeleteSurvey(ctx context.Context, id int) error {
	return s.repo.DeleteSurvey(ctx, id)
}

// AddQuestion adds a question to a survey
func (s *SurveyService) AddQuestion(ctx context.Context, req *models.CreateQuestionRequest) (int, error) {
	// First, get the current max order number for questions in this survey
	questions, err := s.repo.GetQuestionsBySurveyID(ctx, req.SurveyID)
	if err != nil {
		return 0, err
	}

	orderNum := 1
	if len(questions) > 0 {
		// Find the maximum order number and increment
		maxOrderNum := 0
		for _, q := range questions {
			if q.OrderNum > maxOrderNum {
				maxOrderNum = q.OrderNum
			}
		}
		orderNum = maxOrderNum + 1
	}

	// Create the question
	question := &models.Question{
		SurveyID: req.SurveyID,
		Text:     req.Text,
		Type:     req.Type,
		Required: req.Required,
		OrderNum: orderNum,
	}

	questionID, err := s.repo.CreateQuestion(ctx, question)
	if err != nil {
		return 0, err
	}

	// If it's a single_choice question, add options
	if question.Type == "single_choice" && len(req.Options) > 0 {
		for i, optReq := range req.Options {
			option := &models.QuestionOption{
				QuestionID: questionID,
				Text:       optReq.Text,
				OrderNum:   i + 1,
			}

			_, err := s.repo.CreateQuestionOption(ctx, option)
			if err != nil {
				return 0, err
			}
		}
	}

	return questionID, nil
}

// UpdateQuestion updates a question and its options
func (s *SurveyService) UpdateQuestion(ctx context.Context, question *models.Question, options []*models.QuestionOption) error {
	// Update question
	err := s.repo.UpdateQuestion(ctx, question)
	if err != nil {
		return err
	}

	// If it's a single_choice question, update options
	if question.Type == "single_choice" {
		// Delete existing options
		err = s.repo.DeleteQuestionOptions(ctx, question.ID)
		if err != nil {
			return err
		}

		// Add new options
		for _, option := range options {
			_, err := s.repo.CreateQuestionOption(ctx, option)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DeleteQuestion deletes a question
func (s *SurveyService) DeleteQuestion(ctx context.Context, id int) error {
	return s.repo.DeleteQuestion(ctx, id)
}

// UpdateSurveyStatus updates the status of a survey
func (s *SurveyService) UpdateSurveyStatus(ctx context.Context, id int, isActive bool) error {
	return s.repo.UpdateSurveyStatus(ctx, id, isActive)
}
