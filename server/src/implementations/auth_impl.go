package implementations

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/Mahaveer86619/ImaginAI/src/database"
	helpers "github.com/Mahaveer86619/ImaginAI/src/helpers"
	middleware "github.com/Mahaveer86619/ImaginAI/src/middleware"
	services "github.com/Mahaveer86619/ImaginAI/src/services"
	types "github.com/Mahaveer86619/ImaginAI/src/types"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthenticateUser(credentials *types.AuthenticatingCredentials) (*types.UserResponse, int, error) {
	conn := db.GetDBConnection()

	query := `SELECT id, name, email, password, gemini_api_key FROM users WHERE email = $1`
	var user types.User

	err := conn.QueryRow(query, credentials.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.GeminiAPIKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, fmt.Errorf("user not found")
		}
		return nil, http.StatusInternalServerError, fmt.Errorf("error querying user: %w", err)
	}

	if credentials.Password != user.Password {
		return nil, http.StatusUnauthorized, fmt.Errorf("invalid password")
	}

	token, err := middleware.GenerateToken(user.Email)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error generating token: %w", err)
	}

	refreshToken, err := middleware.GenerateRefreshToken(user.Email)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error generating refresh token: %w", err)
	}

	return user.ToUserResponseWithTokens(token, refreshToken), http.StatusOK, nil
}

func RegisterUser(credentials *types.RegisteringCredentials) (*types.UserResponse, int, error) {
	conn := db.GetDBConnection()

	checkQuery := `SELECT id FROM users WHERE email = $1`
	var existingID string
	err := conn.QueryRow(checkQuery, credentials.Email).Scan(&existingID)
	if err == nil {
		return nil, http.StatusConflict, fmt.Errorf("email already registered")
	} else if err != sql.ErrNoRows {
		return nil, http.StatusInternalServerError, fmt.Errorf("error checking user: %w", err)
	}

	insertQuery := `
		INSERT INTO users (id, name, email, password, gemini_api_key, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING name, email, password, gemini_api_key
	`
	var user types.User
	user.ID = uuid.New().String()
	err = conn.QueryRow(insertQuery, user.ID, credentials.Name, credentials.Email, credentials.Password, credentials.GeminiAPIKey).
		Scan(&user.Name, &user.Email, &user.Password, &user.GeminiAPIKey)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error creating user: %w", err)
	}

	token, err := middleware.GenerateToken(user.Email)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error generating token: %w", err)
	}

	refreshToken, err := middleware.GenerateRefreshToken(user.Email)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error generating refresh token: %w", err)
	}

	// Send email to welcome user
	err = services.SendBasicHTMLEmail(
		[]string{user.Email},
		"Welcome to ImaginAI!",
		services.GenerateWelcomeHTML(user.Email),
	)

	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error sending email: %w", err)
	}

	return user.ToUserResponseWithTokens(token, refreshToken), http.StatusCreated, nil
}

func SendPassResetCode(email string) (int, error) {
	conn := db.GetDBConnection()

	insert_query := `
		INSERT INTO forgot_password (id, email, code)
		VALUES ($1, $2, $3) RETURNING *
	`
	select_user_query := `
	  SELECT *
	  FROM users
	  WHERE email = $1
	`

	var authUser types.User

	// Search for user in database
	err := conn.QueryRow(
		select_user_query,
		email,
	).Scan(
		&authUser.ID,
		&authUser.Name,
		&authUser.Email,
		&authUser.Password,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("provided email is not registered: %s", email)
		}
		return http.StatusInternalServerError, fmt.Errorf("error scanning row: %w", err)
	}

	var forgotPassword types.ForgotPassword
	forgotPassword.ID = uuid.New().String()
	forgotPassword.Email = email
	forgotPassword.Code = helpers.Gen6DigitCode()

	err = conn.QueryRow(
		insert_query,
		forgotPassword.ID,
		forgotPassword.Email,
		forgotPassword.Code,
	).Scan(
		&forgotPassword.ID,
		&forgotPassword.Email,
		&forgotPassword.Code,
	)

	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error generating forgot password code: %w", err)
	}

	// Send email with forgot password code
	err = services.SendBasicHTMLEmail(
		[]string{email},
		"Reset your password",
		services.GeneratePasswordResetHTML(forgotPassword.Code, email),
	)

	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error sending email: %w", err)
	}

	return http.StatusOK, nil
}

func CheckResetPassCode(code string, email string) (string, int, error) {
	conn := db.GetDBConnection()

	select_query := `
	  SELECT *
	  FROM forgot_password
	  WHERE email = $1
	`
	del_query := "DELETE FROM forgot_password WHERE id = $1"

	var forgotPassword types.ForgotPassword

	// Search for forgot password in database
	if err := conn.QueryRow(
		select_query,
		email,
	).Scan(
		&forgotPassword.ID,
		&forgotPassword.Email,
		&forgotPassword.Code,
	); err != nil {
		if err == sql.ErrNoRows {
			return "", http.StatusNotFound, fmt.Errorf("forgot password row not found with email: %s", email)
		}
		return "", http.StatusInternalServerError, fmt.Errorf("error scanning row: %w", err)
	}

	// Delete forgot password row if code is correct
	if forgotPassword.Code != code {
		return "", http.StatusBadRequest, fmt.Errorf("invalid code")
	}

	_, err := conn.Exec(del_query, forgotPassword.ID)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("error deleting row: %w", err)
	}

	return forgotPassword.Code, http.StatusOK, nil
}

func RefreshToken(refreshingToken *types.RefreshTokenBody) (*types.RefreshTokenResp, int, error) {
	claims := &middleware.Claims{}
	token, err := jwt.ParseWithClaims(refreshingToken.RefreshTokenKey, claims, func(token *jwt.Token) (interface{}, error) {
		return middleware.JwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, http.StatusUnauthorized, fmt.Errorf("invalid refresh token")
	}

	newToken, err := middleware.GenerateToken(claims.Email)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error generating new token: %w", err)
	}

	newRefreshToken, err := middleware.GenerateRefreshToken(claims.Email)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error generating new token: %w", err)
	}

	return &types.RefreshTokenResp{
		TokenKey:        newToken,
		RefreshTokenKey: newRefreshToken,
	}, http.StatusOK, nil
}
