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
	*CaptchaService
	repository *repository.UserRepository
}

var userService *UserService

func NewUserService(ctx *gin.Context) *UserService {
	if userService == nil {
		userService = &UserService{
			BasicService:   NewBasicService(ctx),
			CaptchaService: NewCaptchaService(ctx),
			repository:     repository.NewUserRepository(ctx),
		}
	}
	return userService
}

func (u *UserService) SignUp(context context.Context, id string, code string, user models.User) error {
	verify := u.CaptchaService.Verify(id, code)
	if !verify {
		return constant.GetError(u.ctx, response.CAPTCHA_VERIFY_FAIL)
	}
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

func (u *UserService) SetRoles(ctx context.Context, uid uint, role_ids []uint) error {
	// TODO:确认有效的用户和角色id
	return u.repository.AssociateRoles(ctx, uid, role_ids)
}
