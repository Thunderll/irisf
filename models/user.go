package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// User 微信小程序用户
type User struct {
	gorm.Model
	Name string `gorm:"comment:名称" json:"name" validate:"required"`

	Username string `gorm:"comment:用户名" json:"username"`
	Password string `gorm:"comment:密码" json:"password"`

	OpenID     string `gorm:"comment:微信用户表示" json:"open_id,omitempty" validate:"required"`
	SessionKey string `gorm:"comment:微信会话密钥" json:"session_key,omitempty"`
}

func (u *User) ParseValidationErrors(errs validator.ValidationErrors) error {
	return nil
}
