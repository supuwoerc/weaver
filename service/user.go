package service

import (
	"context"
	"gin-web/models"
	"gin-web/pkg/constant"
	"gin-web/pkg/jwt"
	"gin-web/pkg/response"
	"gin-web/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	*BasicService
	repository *repository.UserRepository
}

var userService *UserService

func NewUserService(ctx *gin.Context) *UserService {
	if userService == nil {
		userService = &UserService{
			BasicService: NewBasicService(ctx),
			repository:   repository.NewUserRepository(ctx),
		}
	}
	return userService
}

func (u *UserService) SignUp(context context.Context, user models.User) error {
	password, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	var pwd = string(password)
	user.Password = &pwd
	return u.repository.Create(context, user)
}

func (u *UserService) Login(ctx context.Context, email string, password string) (models.User, *models.TokenPair, error) {
	user, err := u.repository.FindByEmail(ctx, email)
	if err != nil {
		return models.User{}, nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password))
	if err != nil {
		return models.User{}, nil, constant.GetError(u.ctx, response.USER_LOGIN_FAIL)
	}
	builder := jwt.NewJwtBuilder(u.ctx)
	accessToken, refreshToken, err := builder.GenerateAccessAndRefreshToken(&jwt.TokenClaimsBasic{
		UID:      user.ID,
		Email:    user.Email,
		NickName: user.NickName,
		Gender:   user.Gender,
		About:    user.About,
		Birthday: user.Birthday,
	})
	if err != nil {
		return models.User{}, nil, err
	}
	err = u.repository.CacheTokenPair(ctx, user.Email, &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return models.User{}, nil, err
	}
	return user, &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
