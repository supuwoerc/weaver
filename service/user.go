package service

import (
	"context"
	"gin-web/models"
	"gin-web/pkg/jwt"
	"gin-web/pkg/response"
	"gin-web/repository"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"sync"
)

type UserService struct {
	*BasicService
	*CaptchaService
	*RoleService
	repository *repository.UserRepository
}

var (
	userOnce    sync.Once
	userService *UserService
)

func NewUserService() *UserService {
	userOnce.Do(func() {
		userService = &UserService{
			BasicService:   NewBasicService(),
			CaptchaService: NewCaptchaService(),
			RoleService:    NewRoleService(),
			repository:     repository.NewUserRepository(),
		}
	})
	return userService
}

func (u *UserService) SignUp(ctx context.Context, id string, code string, user models.User) error {
	verify := u.CaptchaService.Verify(id, code)
	if !verify {
		return response.CaptchaVerifyFail
	}
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	var pwd = string(password)
	user.Password = pwd
	return u.repository.Create(ctx, user)
}

func (u *UserService) Login(ctx context.Context, email string, password string) (*models.User, *models.TokenPair, error) {
	user, err := u.repository.FindByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, nil, response.UserLoginFail
	}
	pair, err := u.repository.GetTokenPair(ctx, email)
	builder := jwt.NewJwtBuilder()
	if err == nil && pair != nil {
		claims, parseErr := builder.ParseToken(pair.AccessToken)
		if parseErr == nil && claims != nil {
			// 如果缓存的token还未过期,直接返回缓存中的记录
			return user, pair, nil
		}
	}
	roleIds := lo.Map[*models.Role, uint](user.Roles, func(item *models.Role, _ int) uint {
		return item.ID
	})
	accessToken, refreshToken, err := builder.GenerateAccessAndRefreshToken(&jwt.TokenClaimsBasic{
		UID:      user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
		Roles:    roleIds,
	})
	if err != nil {
		return nil, nil, err
	}
	err = u.repository.CacheTokenPair(ctx, user.Email, &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, nil, err
	}
	return user, &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *UserService) SetRoles(ctx context.Context, uid uint, roleIds []uint) error {
	// TODO:配置ADMIN账户，限制ADMIN账户被更改角色
	user, err := u.repository.FindByUid(ctx, uid, false)
	if err != nil {
		return err
	}
	if user.ID == 0 {
		return response.UserNotExist
	}
	validIds, err := u.RoleService.FilterValidRoles(ctx, roleIds)
	if err != nil {
		return err
	}
	if len(validIds) == 0 {
		return response.NoValidRoles
	}
	return u.repository.AssociateRoles(ctx, uid, validIds)
}

func (u *UserService) GetRoles(ctx context.Context, uid uint) ([]*models.Role, error) {
	return u.repository.FindRolesByUid(ctx, uid)
}

func (u *UserService) Profile(ctx context.Context, uid uint) (*models.User, error) {
	return u.repository.FindByUid(ctx, uid, true)
}
