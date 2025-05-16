package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_service/internal/api/handlers/middleware"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_service/internal/domain/models"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/survey_service/internal/store/mongodb"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SurveyHandler обробник опитувань
type SurveyHandler struct {
	repository mongodb.Repository
}

// NewSurveyHandler створює новий обробник опитувань
func NewSurveyHandler(repo mongodb.Repository) *SurveyHandler {
	return &SurveyHandler{
		repository: repo,
	}
}

// Create обробляє запит на створення опитування
func (h *SurveyHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Отримуємо дані користувача з контексту
	userCtx, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Парсинг JSON
	var createReq models.CreateSurveyRequest
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Валідація вхідних даних
	if err := validateSurveyRequest(&createReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Створення об'єкта опитування
	survey := models.Survey{
		ID:          primitive.NewObjectID(),
		Title:       createReq.Title,
		Description: createReq.Description,
		OwnerID:     userCtx.UserID,
		Questions:   createReq.Questions,
	}

	// Збереження в базу даних
	if err := h.repository.Create(r.Context(), &survey); err != nil {
		http.Error(w, "Failed to create survey", http.StatusInternalServerError)
		return
	}

	// Відправлення відповіді
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(survey.ToResponse())
}

// GetByID обробляє запит на отримання опитування за ID
func (h *SurveyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Отримуємо дані користувача з контексту
	userCtx, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Отримання ID з URL
	vars := mux.Vars(r)
	surveyID := vars["surveyId"]

	// Отримання опитування з репозиторію
	survey, err := h.repository.GetByID(r.Context(), surveyID)
	if err != nil {
		http.Error(w, "Survey not found", http.StatusNotFound)
		return
	}

	// Перевірка прав доступу (власник або адміністратор)
	if survey.OwnerID != userCtx.UserID && userCtx.Role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Відправлення відповіді
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(survey.ToResponse())
}

// Update обробляє запит на оновлення опитування
func (h *SurveyHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Отримуємо дані користувача з контексту
	userCtx, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Отримання ID з URL
	vars := mux.Vars(r)
	surveyID := vars["surveyId"]
	objectID, err := primitive.ObjectIDFromHex(surveyID)
	if err != nil {
		http.Error(w, "Invalid survey ID", http.StatusBadRequest)
		return
	}

	// Отримання опитування з репозиторію
	existingSurvey, err := h.repository.GetByID(r.Context(), surveyID)
	if err != nil {
		http.Error(w, "Survey not found", http.StatusNotFound)
		return
	}

	// Перевірка прав доступу (тільки власник)
	if existingSurvey.OwnerID != userCtx.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Парсинг JSON
	var updateReq models.UpdateSurveyRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Оновлення полів
	if updateReq.Title != "" {
		existingSurvey.Title = updateReq.Title
	}

	if updateReq.Description != "" {
		existingSurvey.Description = updateReq.Description
	}

	if updateReq.Questions != nil {
		// Валідація питань
		if err := validateQuestions(updateReq.Questions); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		existingSurvey.Questions = updateReq.Questions
	}

	// Оновлення ID
	existingSurvey.ID = objectID

	// Збереження оновленого опитування
	if err := h.repository.Update(r.Context(), existingSurvey); err != nil {
		http.Error(w, "Failed to update survey", http.StatusInternalServerError)
		return
	}

	// Відправлення відповіді
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingSurvey.ToResponse())
}

// Delete обробляє запит на видалення опитування
func (h *SurveyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Отримуємо дані користувача з контексту
	userCtx, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Отримання ID з URL
	vars := mux.Vars(r)
	surveyID := vars["surveyId"]

	// Отримання опитування з репозиторію для перевірки прав
	existingSurvey, err := h.repository.GetByID(r.Context(), surveyID)
	if err != nil {
		http.Error(w, "Survey not found", http.StatusNotFound)
		return
	}

	// Перевірка прав доступу (тільки власник)
	if existingSurvey.OwnerID != userCtx.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Видалення опитування
	if err := h.repository.Delete(r.Context(), surveyID); err != nil {
		http.Error(w, "Failed to delete survey", http.StatusInternalServerError)
		return
	}

	// Відправлення успішної відповіді
	w.WriteHeader(http.StatusNoContent)
}

// GetUserSurveys обробляє запит на отримання списку опитувань користувача
func (h *SurveyHandler) GetUserSurveys(w http.ResponseWriter, r *http.Request) {
	// Отримуємо дані користувача з контексту
	userCtx, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Отримання ID користувача з URL
	vars := mux.Vars(r)
	userID := vars["userId"]

	// Перевірка прав доступу (власний список або адміністратор)
	if userID != userCtx.UserID && userCtx.Role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Параметри пагінації
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.ParseInt(r.URL.Query().Get("per_page"), 10, 64)
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	// Отримання списку опитувань
	surveys, total, err := h.repository.GetByOwnerID(r.Context(), userID, page, perPage)
	if err != nil {
		http.Error(w, "Failed to get surveys", http.StatusInternalServerError)
		return
	}

	// Конвертація в відповідь API
	surveysResponse := models.SurveyListResponse{
		Surveys:    make([]models.SurveyResponse, 0, len(surveys)),
		TotalCount: total,
		Page:       page,
		PerPage:    perPage,
	}

	for _, survey := range surveys {
		surveysResponse.Surveys = append(surveysResponse.Surveys, survey.ToResponse())
	}

	// Відправлення відповіді
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(surveysResponse)
}

// validateSurveyRequest перевіряє дані запиту на створення опитування
func validateSurveyRequest(req *models.CreateSurveyRequest) error {
	if req.Title == "" {
		return errors.New("title is required")
	}

	if len(req.Questions) == 0 {
		return errors.New("at least one question is required")
	}

	return validateQuestions(req.Questions)
}

// validateQuestions перевіряє валідність питань
func validateQuestions(questions []models.Question) error {
	for i, q := range questions {
		if q.Text == "" {
			return errors.New("question text is required")
		}

		switch q.Type {
		case models.SingleChoice, models.MultipleChoice:
			if len(q.Options) < 2 {
				return errors.New("choice questions must have at least 2 options")
			}
		case models.Scale:
			if q.ScaleSettings == nil {
				return errors.New("scale questions must have scale settings")
			}
			if q.ScaleSettings.Min >= q.ScaleSettings.Max {
				return errors.New("scale min must be less than max")
			}
		case models.MatrixSingle, models.MatrixMultiple:
			if len(q.MatrixRows) == 0 {
				return errors.New("matrix questions must have at least one row")
			}
			if len(q.MatrixColumns) == 0 {
				return errors.New("matrix questions must have at least one column")
			}
		case models.OpenText:
			// Немає особливих вимог для текстових питань
		default:
			return errors.New("invalid question type at index " + strconv.Itoa(i))
		}

		// Валідація логіки відображення, якщо вона вказана
		if q.DisplayLogic != nil && q.DisplayLogic.DependsOnQuestionID != "" {
			found := false
			// Перевіряємо, чи існує питання, від якого залежить це питання
			for _, prevQ := range questions {
				if prevQ.ID == q.DisplayLogic.DependsOnQuestionID {
					found = true
					break
				}
			}
			if !found {
				return errors.New("display logic references non-existent question")
			}
		}
	}

	return nil
}
