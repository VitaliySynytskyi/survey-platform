package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/user_service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
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

// GetUserByID finds a user by ID
func (s *PostgresStore) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
	SELECT id, email, hashed_password, role, first_name, last_name, created_at, updated_at
	FROM users
	WHERE id = $1
	`

	var user domain.User
	var role string
	var firstName, lastName sql.NullString

	err := s.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&role,
		&firstName,
		&lastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	user.Role = domain.Role(role)

	if firstName.Valid {
		user.FirstName = firstName.String
	}

	if lastName.Valid {
		user.LastName = lastName.String
	}

	return &user, nil
}

// GetUserByEmail finds a user by email
func (s *PostgresStore) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
	SELECT id, email, hashed_password, role, first_name, last_name, created_at, updated_at
	FROM users
	WHERE email = $1
	`

	var user domain.User
	var role string
	var firstName, lastName sql.NullString

	err := s.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&role,
		&firstName,
		&lastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	user.Role = domain.Role(role)

	if firstName.Valid {
		user.FirstName = firstName.String
	}

	if lastName.Valid {
		user.LastName = lastName.String
	}

	return &user, nil
}

// UpdateUser updates a user in the database
func (s *PostgresStore) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `
	UPDATE users
	SET email = $1, hashed_password = $2, role = $3, first_name = $4, last_name = $5, updated_at = $6
	WHERE id = $7
	`

	_, err := s.db.Exec(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
		user.Role,
		sql.NullString{String: user.FirstName, Valid: user.FirstName != ""},
		sql.NullString{String: user.LastName, Valid: user.LastName != ""},
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// ListUsers returns a list of all users with pagination
func (s *PostgresStore) ListUsers(ctx context.Context, offset, limit int) ([]*domain.User, error) {
	query := `
	SELECT id, email, hashed_password, role, first_name, last_name, created_at, updated_at
	FROM users
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User

	for rows.Next() {
		var user domain.User
		var role string
		var firstName, lastName sql.NullString

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&role,
			&firstName,
			&lastName,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}

		user.Role = domain.Role(role)

		if firstName.Valid {
			user.FirstName = firstName.String
		}

		if lastName.Valid {
			user.LastName = lastName.String
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}

// CountUsers returns the total number of users
func (s *PostgresStore) CountUsers(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int
	err := s.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// EmailExists checks if an email already exists, excluding the given user ID
func (s *PostgresStore) EmailExists(ctx context.Context, email string, excludeID uuid.UUID) (bool, error) {
	query := `
	SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND id != $2)
	`

	var exists bool
	err := s.db.QueryRow(ctx, query, email, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if email exists: %w", err)
	}

	return exists, nil
}
