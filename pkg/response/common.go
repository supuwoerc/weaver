package response

import (
	"github.com/supuwoerc/weaver/models"
)

type Creator struct {
	*models.User
	Status    any `json:"status,omitempty"`
	Gender    any `json:"gender,omitempty"`
	Birthday  any `json:"birthday,omitempty"`
	Roles     any `json:"roles,omitempty"`
	CreatedAt any `json:"created_at,omitempty"`
	UpdatedAt any `json:"updated_at,omitempty"`
}

type Updater = Creator
type SimpleUser = Creator
