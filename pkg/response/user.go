package response

import (
	"github.com/samber/lo"
	"github.com/supuwoerc/weaver/models"
)

// LoginResponse 登录响应
type LoginResponse struct {
	User         LoginUser `json:"user"`          // 用户信息
	Token        string    `json:"token"`         // token
	RefreshToken string    `json:"refresh_token"` // refresh token
}

type LoginUser struct {
	ID       uint    `json:"id"`       // ID
	Email    string  `json:"email"`    // 邮箱
	Nickname *string `json:"nickname"` // 昵称
}

// RefreshTokenResponse 刷新 token 的响应
type RefreshTokenResponse struct {
	Token string `json:"token"` // token
}

// ProfileResponse 个人信息响应
type ProfileResponse struct {
	*models.User
	Roles       []*ProfileResponseRole `json:"roles"`       // 角色
	Departments []*ProfileResponseDept `json:"departments"` // 部门
}
type ProfileResponseRole struct {
	ID   uint   `json:"id"`   //  角色ID
	Name string `json:"name"` // 角色名称
}
type ProfileResponseDept struct {
	ID   uint   `json:"id"`   // 部门ID
	Name string `json:"name"` // 部门名称
}

type UserListRowResponse struct {
	*models.User
	Roles       []*ProfileResponseRole `json:"roles" gorm:"many2many:user_role;"`             // 角色
	Departments []*ProfileResponseDept `json:"departments" gorm:"many2many:user_department;"` // 部门
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
