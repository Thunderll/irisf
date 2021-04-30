package services

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"iris_project_foundation/common"
	"iris_project_foundation/common/api_error"
	"iris_project_foundation/models"
	"iris_project_foundation/module/libs"
)

type UserQuery struct {
	Page     int64 `url:"page"`
	PageSize int64 `url:"page_size"`
}

type UserForm struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

func (u *UserForm) ParseValidationErrors(errs validator.ValidationErrors) error {
	for _, err := range errs {
		switch err.StructField() {
		case "Name":
			return &api_error.BaseAPIError{ErrorCode: 10013, Message: "请输入用户名"}
		}
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////

// IUserService UserService接口
type IUserService interface {
	GetUserList() (*common.Response, error)
	GetUser(int64) (*common.Response, error)
	CreateUser() error
}

// UserService _
type UserService struct {
	Ctx iris.Context
}

func NewUserService(ctx iris.Context) *UserService {
	return &UserService{ctx}
}

func (u *UserService) GetUserList() (*common.Response, error) {
	var (
		err    error
		query  UserQuery
		offset int
		users  []*models.User
	)

	offset = int(query.PageSize * query.Page)
	if err = models.DB.Limit(int(query.Page)).Offset(offset).Find(&users).Error; err != nil {
		return nil, api_error.SqlQueryError
	}
	return common.Success(users, 0, 0), nil
}

func (u *UserService) GetUser(uid int64) (*common.Response, error) {
	var (
		err  error
		user models.User
	)

	if err = models.DB.First(&user, uid).Error; err != nil {
		if ok := errors.Is(err, gorm.ErrRecordNotFound); ok {
			return nil, api_error.ResourceNotFoundError
		} else {
			return nil, api_error.SqlQueryError
		}
	}

	return common.Success(user, 0, 0), nil
}

func (u *UserService) CreateUser() error {
	var (
		err            error
		errs           validator.ValidationErrors
		ok             bool
		userForm       UserForm
		user           *models.User
		openID         string
		sessionKey     string
		hashedPassword []byte
	)

	err = u.Ctx.ReadJSON(&userForm)
	if errs, ok = err.(validator.ValidationErrors); ok {
		return userForm.ParseValidationErrors(errs)
	}

	if userForm.Code != "" {
		if openID, sessionKey, err = libs.WechatAuthorize(userForm.Code); err != nil {
			return err
		}

		user = &models.User{
			Name:       userForm.Name,
			OpenID:     openID,
			SessionKey: sessionKey,
		}
	} else {
		if hashedPassword, err = bcrypt.GenerateFromPassword([]byte(userForm.Password), bcrypt.DefaultCost); err != nil {
			return api_error.WebUserCreateError
		}
		user = &models.User{
			Name:     userForm.Name,
			Username: userForm.Username,
			Password: string(hashedPassword),
		}
	}

	if err = models.DB.Create(&user).Error; err != nil {
		return api_error.WebUserCreateError
	}
	return nil
}
