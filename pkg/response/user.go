package response

import "gin-web/models"

type LoginResponse struct {
	User         models.User `json:"user"`
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
}
