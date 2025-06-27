package implementations

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/Mahaveer86619/ImaginAI/src/database"
	types "github.com/Mahaveer86619/ImaginAI/src/types"
)

func GetAllUsers() ([]*types.UserSafeResponse, int, error) {
	conn := db.GetDBConnection()

	query := `SELECT id, name, email, password, gemini_api_key FROM users`
	rows, err := conn.Query(query)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error querying users: %w", err)
	}
	defer rows.Close()

	var users []*types.UserSafeResponse
	for rows.Next() {
		var user types.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.GeminiAPIKey); err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("error scanning row: %w", err)
		}
		users = append(users, user.ToUserSafeResponse())
	}

	return users, http.StatusOK, nil
}

func GetUserByID(userID string) (*types.UserSafeResponse, int, error) {
	conn := db.GetDBConnection()

	query := `SELECT id, name, email, password, gemini_api_key FROM users WHERE id = $1`
	var user types.User

	err := conn.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.GeminiAPIKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, fmt.Errorf("user not found")
		}
		return nil, http.StatusInternalServerError, fmt.Errorf("error querying user: %w", err)
	}

	return user.ToUserSafeResponse(), http.StatusOK, nil
}

func UpdateUser(user *types.UserSafeResponse) (*types.UserSafeResponse, int, error) {
	conn := db.GetDBConnection()

	search_query := `SELECT id FROM users WHERE id = $1`

	update_query := `UPDATE users 
	SET name = $1, email = $2, gemini_api_key = $3
	WHERE id = $4
	RETURNING id, name, email, gemini_api_key`

	// Check if the user exists
	_, err := conn.Exec(search_query, user.ID)

	if err != nil {
		// Some error occurred
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, fmt.Errorf("user not found with id: %s", user.ID)
		}
		return nil, http.StatusInternalServerError, fmt.Errorf("error scanning row: %w", err)
	} else {
		// Update the user
		var userResp types.User

		err = conn.QueryRow(
			update_query,
			user.Name,
			user.Email,
			user.GeminiAPIKey,
			user.ID,
		).Scan(
			&userResp.ID,
			&userResp.Name,
			&userResp.Email,
			&userResp.GeminiAPIKey,
		)

		if err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("error updating user profile: %w", err)
		}

		// Generate user response
		userResponse := userResp.ToUserSafeResponse()
		return userResponse, http.StatusOK, nil
	}
}

func DeleteUser(userId string) (int, error) {
	conn := db.GetDBConnection()

	search_query := `SELECT id FROM users WHERE id = $1`
	del_query := "DELETE FROM users WHERE id = $1"

	// Check if the user exists
	_, err := conn.Exec(search_query, userId)

	if err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("user not found with id: %s", userId)
		}
		return http.StatusInternalServerError, fmt.Errorf("error scanning row: %w", err)
	} else {
		_, err = conn.Exec(del_query, userId)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error deleting row: %w", err)
		}
	}

	return http.StatusOK, nil
}
