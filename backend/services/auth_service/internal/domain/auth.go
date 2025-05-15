package domain

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TokenResponse represents the response after successful login or token refresh
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // Expiration time in seconds
}

// RefreshRequest represents the request body for token refresh
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Validate validates the register request
func (r *RegisterRequest) Validate() map[string]string {
	errors := make(map[string]string)

	if r.Email == "" {
		errors["email"] = "Email is required"
	}

	if r.Password == "" {
		errors["password"] = "Password is required"
	} else if len(r.Password) < 6 {
		errors["password"] = "Password must be at least 6 characters long"
	}

	return errors
}

// Validate validates the login request
func (r *LoginRequest) Validate() map[string]string {
	errors := make(map[string]string)

	if r.Email == "" {
		errors["email"] = "Email is required"
	}

	if r.Password == "" {
		errors["password"] = "Password is required"
	}

	return errors
}

// Validate validates the refresh token request
func (r *RefreshRequest) Validate() map[string]string {
	errors := make(map[string]string)

	if r.RefreshToken == "" {
		errors["refresh_token"] = "Refresh token is required"
	}

	return errors
}
