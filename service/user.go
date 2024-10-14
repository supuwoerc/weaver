package service

import (
	"gin-web/models"
	"gin-web/pkg/constant"
	"gin-web/pkg/jwt"
	"gin-web/pkg/response"
	"gin-web/repository"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	*BasicService
	*CaptchaService
	*RoleService
	repository *repository.UserRepository
}

func NewUserService(ctx *gin.Context) *UserService {
	return &UserService{
		BasicService:   NewBasicService(ctx),
		CaptchaService: NewCaptchaService(ctx),
		RoleService:    NewRoleService(ctx),
		repository:     repository.NewUserRepository(ctx),
	}
}

func (u *UserService) SignUp(id string, code string, user models.User) error {
	verify := u.CaptchaService.Verify(id, code)
	if !verify {
		return constant.GetError(u.ctx, response.CaptchaVerifyFail)
	}
	password, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	var pwd = string(password)
	user.Password = &pwd
	return u.repository.Create(u.ctx.Request.Context(), user)
}

func (u *UserService) Login(email string, password string) (*models.User, *models.TokenPair, error) {
	user, err := u.repository.FindByEmail(u.ctx.Request.Context(), email)
	if err != nil {
		return nil, nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password))
	if err != nil {
		return nil, nil, constant.GetError(u.ctx, response.UserLoginFail)
	}
	pair, err := u.repository.GetTokenPair(u.ctx.Request.Context(), email)
	builder := jwt.NewJwtBuilder(u.ctx)
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
	err = u.repository.CacheTokenPair(u.ctx.Request.Context(), user.Email, &models.TokenPair{
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

func (u *UserService) SetRoles(uid uint, roleIds []uint) error {
	// TODO:配置ADMIN账户，限制ADMIN账户被更改角色
	user, err := u.repository.FindByUid(u.ctx.Request.Context(), uid, false)
	if err != nil {
		return err
	}
	if user.ID == 0 {
		return constant.GetError(u.ctx, response.UserNotExist)
	}
	validIds, err := u.RoleService.FilterValidRoles(roleIds)
	if err != nil {
		return err
	}
	if len(validIds) == 0 {
		return constant.GetError(u.ctx, response.NoValidRoles)
	}
	return u.repository.AssociateRoles(u.ctx.Request.Context(), uid, validIds)
}

func (u *UserService) GetRoles(uid uint) ([]*models.Role, error) {
	return u.repository.FindRolesByUid(u.ctx.Request.Context(), uid)
}

func (u *UserService) Profile(uid uint) (*models.User, error) {
	return u.repository.FindByUid(u.ctx.Request.Context(), uid, true)
}
