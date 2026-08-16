package application

type RegisterUserInput struct {
	Email    string
	Password string
}

type RegisterUserOutput struct {
	UserID string
	Email  string
}

type LoginUserInput struct {
	Email    string
	Password string
}

type LoginUserOutput struct {
	AccessToken  string
	RefreshToken string
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type GetMeOutput struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}
