package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/VitaliySynytskyi/survey-platform/auth-service/internal/models"
)

func TestGetUserByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)
	ctx := context.Background()

	t.Run("User exists", func(t *testing.T) {
		username := "testuser"
		userID := 1
		email := "test@example.com"
		firstName := "Test"
		lastName := "User"
		isActive := true
		passwordHash := "hashed_password"

		rows := sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "first_name", "last_name", "is_active", "created_at", "updated_at"}).
			AddRow(userID, username, email, passwordHash, firstName, lastName, isActive, time.Now(), time.Now())

		mock.ExpectQuery("SELECT id, username, email, password_hash, first_name, last_name, is_active, created_at, updated_at FROM users WHERE username = \\$1").
			WithArgs(username).
			WillReturnRows(rows)

		// Mock roles query
		roleRows := sqlmock.NewRows([]string{"name"}).
			AddRow("user")
		mock.ExpectQuery("SELECT r.name FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = \\$1").
			WithArgs(userID).
			WillReturnRows(roleRows)

		user, err := repo.GetUserByUsername(ctx, username)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, username, user.Username)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, passwordHash, user.PasswordHash)
		assert.Equal(t, []string{"user"}, user.Roles)
	})

	t.Run("User does not exist", func(t *testing.T) {
		username := "nonexistentuser"

		mock.ExpectQuery("SELECT id, username, email, password_hash, first_name, last_name, is_active, created_at, updated_at FROM users WHERE username = \\$1").
			WithArgs(username).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByUsername(ctx, username)

		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)
	ctx := context.Background()

	t.Run("Successful user creation", func(t *testing.T) {
		user := &models.User{
			Username:     "newuser",
			Email:        "new@example.com",
			PasswordHash: "hashed_password",
			FirstName:    "New",
			LastName:     "User",
			IsActive:     true,
		}

		mock.ExpectQuery("INSERT INTO users").
			WithArgs(user.Username, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.IsActive).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		userID, err := repo.CreateUser(ctx, user)

		assert.NoError(t, err)
		assert.Equal(t, 1, userID)
	})

	t.Run("Failed user creation", func(t *testing.T) {
		user := &models.User{
			Username:     "existinguser",
			Email:        "existing@example.com",
			PasswordHash: "hashed_password",
			IsActive:     true,
		}

		mock.ExpectQuery("INSERT INTO users").
			WithArgs(user.Username, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.IsActive).
			WillReturnError(err)

		userID, err := repo.CreateUser(ctx, user)

		assert.Error(t, err)
		assert.Equal(t, 0, userID)
	})
}

func TestGetRoleByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)
	ctx := context.Background()

	t.Run("Role exists", func(t *testing.T) {
		roleName := "admin"
		roleID := 1

		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(roleID, roleName)

		mock.ExpectQuery("SELECT id, name FROM roles WHERE name = \\$1").
			WithArgs(roleName).
			WillReturnRows(rows)

		role, err := repo.GetRoleByName(ctx, roleName)

		assert.NoError(t, err)
		assert.NotNil(t, role)
		assert.Equal(t, roleID, role.ID)
		assert.Equal(t, roleName, role.Name)
	})

	t.Run("Role does not exist", func(t *testing.T) {
		roleName := "nonexistentrole"

		mock.ExpectQuery("SELECT id, name FROM roles WHERE name = \\$1").
			WithArgs(roleName).
			WillReturnError(sql.ErrNoRows)

		role, err := repo.GetRoleByName(ctx, roleName)

		assert.Error(t, err)
		assert.Nil(t, role)
	})
}
