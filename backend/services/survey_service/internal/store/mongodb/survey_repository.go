package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/VitaliySynytskyi/microservices-survey-app/backend/services/survey_service/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// CollectionName ім'я колекції в MongoDB
	CollectionName = "surveys"
)

// Repository інтерфейс репозиторію опитувань
type Repository interface {
	Create(ctx context.Context, survey *models.Survey) error
	GetByID(ctx context.Context, id string) (*models.Survey, error)
	Update(ctx context.Context, survey *models.Survey) error
	Delete(ctx context.Context, id string) error
	GetByOwnerID(ctx context.Context, ownerID string, page, perPage int64) ([]models.Survey, int64, error)
}

// SurveyRepository реалізація репозиторію опитувань
type SurveyRepository struct {
	collection *mongo.Collection
}

// NewSurveyRepository створює новий екземпляр репозиторію опитувань
func NewSurveyRepository(db *mongo.Database) Repository {
	collection := db.Collection(CollectionName)

	// Створення індексів
	_, _ = collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "owner_id", Value: 1}},
			Options: options.Index().SetBackground(true),
		},
	)

	return &SurveyRepository{
		collection: collection,
	}
}

// Create створює нове опитування
func (r *SurveyRepository) Create(ctx context.Context, survey *models.Survey) error {
	// Встановлення часу створення та оновлення
	now := time.Now()
	survey.CreatedAt = now
	survey.UpdatedAt = now

	// Генерація ID для кожного питання, якщо вони відсутні
	for i := range survey.Questions {
		if survey.Questions[i].ID == "" {
			survey.Questions[i].ID = primitive.NewObjectID().Hex()
		}
	}

	// Вставка документа
	result, err := r.collection.InsertOne(ctx, survey)
	if err != nil {
		return err
	}

	// Присвоєння ID створеному опитуванню
	survey.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID отримує опитування за ID
func (r *SurveyRepository) GetByID(ctx context.Context, id string) (*models.Survey, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var survey models.Survey
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&survey)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("survey not found")
		}
		return nil, err
	}

	return &survey, nil
}

// Update оновлює опитування
func (r *SurveyRepository) Update(ctx context.Context, survey *models.Survey) error {
	// Встановлення часу оновлення
	survey.UpdatedAt = time.Now()

	// Генерація ID для нових питань
	for i := range survey.Questions {
		if survey.Questions[i].ID == "" {
			survey.Questions[i].ID = primitive.NewObjectID().Hex()
		}
	}

	// Оновлення документа
	_, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": survey.ID},
		survey,
	)
	return err
}

// Delete видаляє опитування за ID
func (r *SurveyRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// GetByOwnerID отримує список опитувань за ID власника з пагінацією
func (r *SurveyRepository) GetByOwnerID(ctx context.Context, ownerID string, page, perPage int64) ([]models.Survey, int64, error) {
	// Параметри пагінації
	skip := (page - 1) * perPage
	limit := perPage

	// Опції пошуку
	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	// Запит
	cursor, err := r.collection.Find(ctx, bson.M{"owner_id": ownerID}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Декодування результатів
	var surveys []models.Survey
	if err = cursor.All(ctx, &surveys); err != nil {
		return nil, 0, err
	}

	// Отримання загальної кількості
	total, err := r.collection.CountDocuments(ctx, bson.M{"owner_id": ownerID})
	if err != nil {
		return nil, 0, err
	}

	return surveys, total, nil
}
