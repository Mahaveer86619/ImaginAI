package types

type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	GeminiAPIKey string `json:"gemini_api_key"`
}

type UserResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	GeminiAPIKey string `json:"gemini_api_key"`
}

type UserSafeResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	GeminiAPIKey string `json:"gemini_api_key"`
}

func (u *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		GeminiAPIKey: u.GeminiAPIKey,
	}
}

func (u *User) ToUserSafeResponse() *UserSafeResponse {
	return &UserSafeResponse{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
	}
}

func (u *User) ToUserResponseWithTokens(token string, refreshToken string) *UserResponse {
	return &UserResponse{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		Token:        token,
		RefreshToken: refreshToken,
		GeminiAPIKey: u.GeminiAPIKey,
	}
}
