package service

import (
	"context"
	"errors"
	"gin-web/models"
	"gin-web/pkg/constant"
	"gin-web/pkg/email"
	"gin-web/pkg/global"
	"gin-web/pkg/jwt"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"gin-web/repository"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"sync"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string, needRoles, needPermissions, needDepts bool) (*models.User, error)
	GetById(ctx context.Context, uid uint, needRoles, needPermissions, needDepts bool) (*models.User, error)
	GetByIds(ctx context.Context, ids []uint, needRoles, needPermissions, needDepts bool) ([]*models.User, error)
	GetList(ctx context.Context, keyword string, limit, offset int) ([]*models.User, int64, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	CacheTokenPair(ctx context.Context, email string, pair *models.TokenPair) error
	GetTokenPairIsExist(ctx context.Context, email string) (bool, error)
	GetTokenPair(ctx context.Context, email string) (*models.TokenPair, error)
}

type UserEmailClient interface {
	SendHTML(to string, subject constant.Subject, templatePath constant.Template, data any) error
}

type UserService struct {
	*BasicService
	*CaptchaService
	userRepository UserRepository
	roleRepository RoleRepository
	emailClient    UserEmailClient
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
			userRepository: repository.NewUserRepository(),
			roleRepository: repository.NewRoleRepository(),
			emailClient:    email.NewEmailClient(),
		}
	})
	return userService
}

func (u *UserService) SignUp(ctx context.Context, id string, code string, user *models.User) error {
	verify := u.CaptchaService.Verify(constant.SignUp, id, code)
	if !verify {
		return response.CaptchaVerifyFail
	}
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	var pwd = string(password)
	user.Password = pwd
	user.Status = constant.Inactive
	emailLock := utils.NewLock(constant.SignUpEmailPrefix, user.Email)
	if err = utils.Lock(ctx, emailLock); err != nil {
		return err
	}
	defer func(lock *utils.RedisLock) {
		if e := utils.Unlock(lock); e != nil {
			global.Logger.Errorf("unlock fail %s", e.Error())
		}
	}(emailLock)
	existUser, err := u.userRepository.GetByEmail(ctx, user.Email, false, false, false)
	if err != nil && !errors.Is(err, response.UserNotExist) {
		return err
	}
	if existUser != nil {
		return response.UserCreateDuplicateEmail
	}
	return u.Transaction(ctx, false, func(ctx context.Context) error {
		if err = u.userRepository.Create(ctx, user); err != nil {
			return err
		}
		// TODO:生成唯一的激活链接 & 重新发送邮件的机制 & 激活账户的机制
		if err = u.emailClient.SendHTML(user.Email, constant.Signup, constant.SignupTemplate, user); err != nil {
			return err
		}
		return nil
	})
}

func (u *UserService) Login(ctx context.Context, email string, password string) (*response.LoginResponse, error) {
	user, err := u.userRepository.GetByEmail(ctx, email, true, false, false)
	switch {
	case err != nil:
		return nil, err
	case user.Status == constant.Inactive:
		return nil, response.UserInactive
	case user.Status == constant.Disabled:
		return nil, response.UserDisabled
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, response.UserLoginFail
	}
	pair, err := u.userRepository.GetTokenPair(ctx, email)
	builder := jwt.NewJwtBuilder()
	if err == nil && pair != nil {
		claims, parseErr := builder.ParseToken(pair.AccessToken)
		if parseErr == nil && claims != nil {
			// 如果缓存的token还未过期,直接返回缓存中的记录
			return &response.LoginResponse{
				User: response.LoginUser{
					ID:       user.ID,
					Email:    user.Email,
					Nickname: user.Nickname,
				},
				Token:        pair.AccessToken,
				RefreshToken: pair.RefreshToken,
			}, nil
		}
	}
	accessToken, refreshToken, err := builder.GenerateAccessAndRefreshToken(&jwt.TokenClaimsBasic{
		ID:       user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
	})
	if err != nil {
		return nil, err
	}
	err = u.userRepository.CacheTokenPair(ctx, user.Email, &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, err
	}
	return &response.LoginResponse{
		User: response.LoginUser{
			ID:       user.ID,
			Email:    user.Email,
			Nickname: user.Nickname,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *UserService) Profile(ctx context.Context, uid uint) (*response.ProfileResponse, error) {
	user, err := u.userRepository.GetById(ctx, uid, true, true, true)
	if err != nil {
		return nil, err
	}
	return &response.ProfileResponse{
		User: user,
		Roles: lo.Map(user.Roles, func(item *models.Role, _ int) *response.ProfileResponseRole {
			return &response.ProfileResponseRole{
				ID:   item.ID,
				Name: item.Name,
			}
		}),
		Departments: lo.Map(user.Departments, func(item *models.Department, _ int) *response.ProfileResponseDept {
			return &response.ProfileResponseDept{
				ID:   item.ID,
				Name: item.Name,
			}
		}),
	}, nil
}

func (u *UserService) GetUserList(ctx context.Context, keyword string, limit, offset int) ([]*response.UserListRowResponse, int64, error) {
	list, total, err := u.userRepository.GetList(ctx, keyword, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return lo.Map(list, func(item *models.User, _ int) *response.UserListRowResponse {
		return response.ToUserListRowResponse(item)
	}), total, nil
}
