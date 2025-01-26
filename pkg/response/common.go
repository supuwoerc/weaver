package response

import (
	"gin-web/models"
)

type Creator struct {
	*models.User
	Gender    any `json:"gender,omitempty"`
	Birthday  any `json:"birthday,omitempty"`
	Roles     any `json:"roles,omitempty"`
	CreatedAt any `json:"created_at,omitempty"`
	UpdatedAt any `json:"updated_at,omitempty"`
}

type Updater = Creator
