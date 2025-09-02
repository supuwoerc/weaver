package response

import (
	"github.com/samber/lo"
	"github.com/supuwoerc/weaver/models"
)

// LoginResponse 登录响应
type LoginResponse struct {
	User         LoginUser `json:"user"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type LoginUser struct {
	ID       uint    `json:"id"`
	Email    string  `json:"email"`
	Nickname *string `json:"nickname"`
}

// RefreshTokenResponse 刷新 token 的响应
type RefreshTokenResponse struct {
	Token string `json:"token"`
}

// ProfileResponse 个人信息响应
type ProfileResponse struct {
	*models.User
	Roles       []*ProfileResponseRole `json:"roles"`
	Departments []*ProfileResponseDept `json:"departments"`
}
type ProfileResponseRole struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
type ProfileResponseDept struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UserListRowResponse struct {
	*models.User
	Roles       []*ProfileResponseRole `json:"roles" gorm:"many2many:user_role;"`
	Departments []*ProfileResponseDept `json:"departments" gorm:"many2many:user_department;"`
}

// ToUserListRowResponse 将 user 转为响应
func ToUserListRowResponse(user *models.User) *UserListRowResponse {
	return &UserListRowResponse{
		User: user,
		Roles: lo.Map(user.Roles, func(item *models.Role, _ int) *ProfileResponseRole {
			return &ProfileResponseRole{
				ID:   item.ID,
				Name: item.Name,
			}
		}),
		Departments: lo.Map(user.Departments, func(item *models.Department, _ int) *ProfileResponseDept {
			return &ProfileResponseDept{
				ID:   item.ID,
				Name: item.Name,
			}
		}),
	}
}
