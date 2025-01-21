package service

import (
	"context"
	"gin-web/models"
	"gin-web/pkg/jwt"
	"gin-web/pkg/response"
	"gin-web/repository"
	"golang.org/x/crypto/bcrypt"
	"sync"
)

type UserService struct {
	*BasicService
	*CaptchaService
	userRepository *repository.UserRepository
	roleRepository *repository.RoleRepository
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
		}
	})
	return userService
}

func (u *UserService) SignUp(ctx context.Context, id string, code string, user *models.User) error {
	verify := u.CaptchaService.Verify(SignUp, id, code)
	if !verify {
		return response.CaptchaVerifyFail
	}
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	var pwd = string(password)
	user.Password = pwd
	return u.userRepository.Create(ctx, user)
}

func (u *UserService) Login(ctx context.Context, email string, password string) (*models.User, *models.TokenPair, error) {
	user, err := u.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, nil, response.UserLoginFail
	}
	pair, err := u.userRepository.GetTokenPair(ctx, email)
	builder := jwt.NewJwtBuilder()
	if err == nil && pair != nil {
		claims, parseErr := builder.ParseToken(pair.AccessToken)
		if parseErr == nil && claims != nil {
			// 如果缓存的token还未过期,直接返回缓存中的记录
			return user, pair, nil
		}
	}
	accessToken, refreshToken, err := builder.GenerateAccessAndRefreshToken(&jwt.TokenClaimsBasic{
		UID:      user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
	})
	if err != nil {
		return nil, nil, err
	}
	err = u.userRepository.CacheTokenPair(ctx, user.Email, &models.TokenPair{
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

func (u *UserService) Profile(ctx context.Context, uid uint) (*models.User, error) {
	return u.userRepository.GetById(ctx, uid, true, true)
}
