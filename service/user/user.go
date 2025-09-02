package user

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/supuwoerc/weaver/models"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/jwt"
	"github.com/supuwoerc/weaver/pkg/response"
	"github.com/supuwoerc/weaver/pkg/utils"
	"github.com/supuwoerc/weaver/service"
	"github.com/supuwoerc/weaver/service/captcha"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
)

type DAO interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string, preload ...string) (*models.User, error)
	GetById(ctx context.Context, uid uint, preload ...string) (*models.User, error)
	GetByIds(ctx context.Context, ids []uint, preload ...string) ([]*models.User, error)
	GetList(ctx context.Context, keyword string, limit, offset int) ([]*models.User, int64, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	UpdateAccountStatus(ctx context.Context, id uint, status constant.UserStatus) error
}
type Cache interface {
	CacheRefreshToken(ctx context.Context, email, refreshToken string, expiration time.Duration) error
	DeleteRefreshToken(ctx context.Context, email string) error
	GetRefreshToken(ctx context.Context, email string) (string, error)
	CacheActiveAccountCode(ctx context.Context, id uint, code string, duration time.Duration) error
	GetActiveAccountCode(ctx context.Context, id uint) (string, error)
	RemoveActiveAccountCode(ctx context.Context, id uint) error
}

type EmailClient interface {
	SendHTML(ctx context.Context, to string, subject constant.Subject, templatePath constant.Template, data any) error
}

type Service struct {
	*service.BasicService
	*captcha.Service
	userDAO      DAO
	userCache    Cache
	tokenBuilder *jwt.TokenBuilder
}

func NewUserService(
	basic *service.BasicService,
	captchaService *captcha.Service,
	userDAO DAO,
	userCache Cache,
	tb *jwt.TokenBuilder,
) *Service {
	return &Service{
		BasicService: basic,
		Service:      captchaService,
		userDAO:      userDAO,
		userCache:    userCache,
		tokenBuilder: tb,
	}
}

func (u *Service) SignUp(ctx context.Context, id string, code string, user *models.User) error {
	verify := u.Service.Verify(constant.SignUp, id, code)
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
	emailLock := u.Locksmith.NewLock(constant.SignUpEmailPrefix, user.Email)
	if err = emailLock.Lock(ctx, true); err != nil {
		return err
	}
	defer func(lock *utils.RedisLock) {
		if e := lock.Unlock(); e != nil {
			u.Logger.WithContext(ctx).Errorf("unlock fail %s", e.Error())
		}
	}(emailLock)
	existUser, err := u.userDAO.GetByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, response.UserNotExist) {
		return err
	}
	if existUser != nil {
		return response.UserCreateDuplicateEmail
	}
	return u.Transaction(ctx, false, func(ctx context.Context) error {
		if err = u.userDAO.Create(ctx, user); err != nil {
			return err
		}
		return u.sendActiveEmail(ctx, user.ID, user.Email)
	})
}

func (u *Service) sendActiveEmail(ctx context.Context, uid uint, email string) error {
	activeURL, makeErr := u.generateActiveURL(ctx, uid)
	if makeErr != nil {
		return makeErr
	}
	variable := models.SignUpVariable{ActiveURL: activeURL}
	err := u.EmailClient.SendHTML(ctx, email, constant.Signup, constant.SignupTemplate, variable)
	if err != nil {
		return err
	}
	return nil
}

func (u *Service) generateActiveURL(ctx context.Context, uid uint) (string, error) {
	activeCode := lo.RandomString(constant.UserActiveCodeLength, lo.LettersCharset)
	baseURL := u.Conf.System.BaseURL
	expiration := u.Conf.Account.Expiration
	if err := u.userCache.CacheActiveAccountCode(ctx, uid, activeCode, expiration*time.Second); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/view/v1/public/user/active?active_code=%s&id=%d", baseURL, activeCode, uid), nil
}

func (u *Service) Login(ctx context.Context, email string, password string) (*response.LoginResponse, error) {
	user, err := u.userDAO.GetByEmail(ctx, email, "Roles")
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
	accessToken, refreshToken, err := u.tokenBuilder.GenerateAccessAndRefreshToken(&jwt.TokenClaimsBasic{
		ID:       user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
	})
	if err != nil {
		return nil, err
	}
	err = u.userCache.CacheRefreshToken(ctx, user.Email, refreshToken, u.tokenBuilder.GetRefreshTokenExpiration())
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

func (u *Service) Logout(ctx context.Context, email string) error {
	return u.userCache.DeleteRefreshToken(ctx, email)
}

func (u *Service) Profile(ctx context.Context, uid uint) (*response.ProfileResponse, error) {
	user, err := u.userDAO.GetById(ctx, uid, "Roles", "Departments")
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

func (u *Service) GetUserList(
	ctx context.Context, keyword string, limit, offset int,
) ([]*response.UserListRowResponse, int64, error) {
	list, total, err := u.userDAO.GetList(ctx, keyword, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return lo.Map(list, func(item *models.User, _ int) *response.UserListRowResponse {
		return response.ToUserListRowResponse(item)
	}), total, nil
}

func (u *Service) ActiveAccount(ctx context.Context, uid uint, activeCode string) error {
	userLock := u.Locksmith.NewLock(constant.UserIdPrefix, strconv.Itoa(int(uid)))
	if lockErr := userLock.Lock(ctx, true); lockErr != nil {
		return lockErr
	}
	defer func(lock *utils.RedisLock) {
		if e := lock.Unlock(); e != nil {
			u.Logger.WithContext(ctx).Errorf("unlock fail %s", e.Error())
		}
	}(userLock)
	code, err := u.userCache.GetActiveAccountCode(ctx, uid)
	if err != nil {
		// key 过期的情况需要重新发送邮件
		if errors.Is(err, redis.Nil) {
			var user *models.User
			user, err = u.userDAO.GetById(ctx, uid)
			if err == nil && user.Status == constant.Inactive {
				go func() {
					if temp := u.sendActiveEmail(context.Background(), user.ID, user.Email); temp != nil {
						u.Logger.WithContext(ctx).Errorf("send active email fail %s", err.Error())
					}
				}()
			}
		}
		return err
	}
	if activeCode != code {
		return response.InvalidActiveCode
	} else {
		return u.Transaction(ctx, false, func(ctx context.Context) error {
			var user *models.User
			user, err = u.userDAO.GetById(ctx, uid)
			if err != nil {
				return err
			}
			if user.Status == constant.Normal {
				return response.ReActiveErr
			}
			if user.Status == constant.Disabled {
				return response.UserDisabled
			}
			if err = u.userDAO.UpdateAccountStatus(ctx, uid, constant.Normal); err != nil {
				go func() {
					if temp := u.sendActiveEmail(context.Background(), user.ID, user.Email); temp != nil {
						u.Logger.WithContext(ctx).Errorf("send active email fail %s", err.Error())
					}
				}()
				return err
			}
			if err = u.userCache.RemoveActiveAccountCode(ctx, uid); err != nil {
				go func() {
					if temp := u.sendActiveEmail(context.Background(), user.ID, user.Email); temp != nil {
						u.Logger.WithContext(ctx).Errorf("send active email fail %s", err.Error())
					}
				}()
				return err
			}
			return nil
		})
	}
}
