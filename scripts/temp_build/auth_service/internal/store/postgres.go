package store

import (
	"context"
	"fmt"
	"time"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/auth_service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresStore implements the UserStore interface with PostgreSQL
type PostgresStore struct {
	db *pgxpool.Pool
}

// NewPostgresStore creates a new PostgreSQL store
func NewPostgresStore(connString string) (*PostgresStore, error) {
	// Create a connection pool
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	// Set connection pool settings
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	// Connect to the database
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Check the connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &PostgresStore{
		db: pool,
	}, nil
}

// Close closes the database connection
func (s *PostgresStore) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

// EnsureSchema ensures that the required database schema exists
func (s *PostgresStore) EnsureSchema(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		hashed_password VARCHAR(255) NOT NULL,
		role VARCHAR(50) NOT NULL DEFAULT 'user',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	`

	_, err := s.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to ensure schema: %w", err)
	}

	return nil
}

// CreateUser creates a new user in the database
func (s *PostgresStore) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
	INSERT INTO users (id, email, hashed_password, role, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := s.db.Exec(
		ctx,
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// FindUserByEmail finds a user by email
func (s *PostgresStore) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
	SELECT id, email, hashed_password, role, created_at, updated_at
	FROM users
	WHERE email = $1
	`

	var user domain.User
	var role string

	err := s.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	user.Role = domain.Role(role)

	return &user, nil
}

// FindUserByID finds a user by ID
func (s *PostgresStore) FindUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
	SELECT id, email, hashed_password, role, created_at, updated_at
	FROM users
	WHERE id = $1
	`

	var user domain.User
	var role string

	err := s.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	user.Role = domain.Role(role)

	return &user, nil
}

// EmailExists checks if an email already exists
func (s *PostgresStore) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `
	SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`

	var exists bool
	err := s.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if email exists: %w", err)
	}

	return exists, nil
}

// Ping checks if the database connection is alive
func (s *PostgresStore) Ping(ctx context.Context) error {
	return s.db.Ping(ctx)
}
