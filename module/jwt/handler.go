package jwt

import (
	"encoding/json"
	"errors"
	"iris_project_foundation/common"
	"iris_project_foundation/common/api_error"
	"iris_project_foundation/config"
	"iris_project_foundation/models"
	"iris_project_foundation/module/libs"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	jwt_ "github.com/kataras/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type TokenResult struct {
	UserID       int64           `json:"user_id"`
	AccessToken  json.RawMessage `json:"access_token"`
	RefreshToken json.RawMessage `json:"refresh_token,omitempty"`
}

// generateToken 根据用户信息生成token
func generateToken(user *models.User) (result *TokenResult, err error) {
	var (
		claims    *UserClaims
		tokenPair jwt.TokenPair
		token     []byte
	)

	claims = &UserClaims{
		ID:   int64(user.ID),
		Name: user.Name,
	}

	// 是否启用Refresh-Token
	if config.GConfig.App.TokenPair {
		if tokenPair, err = signer.NewTokenPair(
			claims, jwt.Claims{Subject: "WechatToken"},
			time.Duration(config.GConfig.App.RefreshTokenExpiration)*time.Minute,
			jwt.Claims{ID: uuid.NewString()},
		); err != nil {
			return nil, api_error.UnauthorizedError
		}
		return &TokenResult{
			UserID:       int64(user.ID),
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
		}, nil
	} else {
		if token, err = signer.Sign(claims, jwt.Claims{ID: uuid.NewString()}); err != nil {
			return nil, api_error.UnauthorizedError
		}
		return &TokenResult{int64(user.ID), jwt_.BytesQuote(token), nil}, nil
	}
}

// WechatObtainToken 微信小程序获取token
// 通过微信临时凭证code到微信服务器换取openid和sessionkey, 根据openid到数据库获取用户信息
func WechatObtainToken(ctx iris.Context) {
	var (
		err error

		code   string
		user   *models.User
		result *TokenResult
	)

	code = ctx.PostValueDefault("code", "")

	if user = wechatAuth(code); user == nil {
		ctx.JSON(common.Failed(api_error.UnauthorizedError))
		return
	}

	if result, err = generateToken(user); err != nil {
		ctx.JSON(common.Failed(err))
	}

	ctx.JSON(common.Success(result, 0, 0))
}

func wechatAuth(code string) *models.User {
	var (
		err    error
		openID string
		user   models.User
	)

	if openID, _, err = libs.WechatAuthorize(code); err != nil {
		return nil
	}

	if err = models.DB.Where("open_id = ?", openID).First(&user).Error; err != nil {
		return nil
	}
	return &user
}

// WebObtainToken 后端登录
// 传统的用户名/密码登录方式
func WebObtainToken(ctx iris.Context) {
	var (
		err     error
		webUser *models.User
		result  *TokenResult
	)

	if webUser, err = webAuth(ctx); err != nil {
		ctx.JSON(common.Failed(err))
	}

	if result, err = generateToken(webUser); err != nil {
		ctx.JSON(common.Failed(err))
	}

	ctx.JSON(common.Success(result, 0, 0))
}

type LoginForm struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func webAuth(ctx iris.Context) (*models.User, error) {
	var (
		err       error
		loginFrom LoginForm
		errs      validator.ValidationErrors
		ok        bool
		webUser   models.User
	)

	if err = ctx.ReadJSON(&loginFrom); err != nil {
		if errs, ok = err.(validator.ValidationErrors); ok {
			for _, err := range errs {
				switch err.StructField() {
				case "Username":
					return nil, api_error.WebAuthFailedError
				case "Password":
					return nil, api_error.WebAuthFailedError
				}
			}
		}
	}

	if err = models.DB.Where("username = ?", loginFrom.Username).First(&webUser).Error; err != nil {
		if ok := errors.Is(err, gorm.ErrRecordNotFound); ok {
			return nil, api_error.WebAuthFailedError
		} else {
			return nil, api_error.SqlQueryError
		}
	}

	if err = bcrypt.CompareHashAndPassword([]byte(webUser.Password), []byte(loginFrom.Password)); err != nil {
		return nil, api_error.WebAuthFailedError
	}

	return &webUser, nil
}

// RefreshToken 刷新token
func RefreshToken(ctx iris.Context) {
	var (
		err           error
		tokenPair     jwt.TokenPair
		refreshToken  json.RawMessage
		verifiedToken *jwt.VerifiedToken
		userClaims    UserClaims
		user          models.User
		result        *TokenResult
	)

	if err = ctx.ReadJSON(&tokenPair); err != nil {
		ctx.JSON(common.Failed(api_error.RefreshTokenError))
		return
	}
	refreshToken = tokenPair.RefreshToken

	// 验证RefreshToken
	if verifiedToken, err = verifier.VerifyToken(refreshToken); err != nil {
		ctx.JSON(common.Failed(api_error.RefreshTokenInvalidError))
		return
	}
	// 取出RefreshToken中的claims
	if err = verifiedToken.Claims(&userClaims); err != nil {
		ctx.JSON(common.Failed(api_error.RefreshTokenError))
		return
	}
	// 根据claims获取用户对象
	if err = models.DB.First(&user, "id = ?", userClaims.ID).Error; err != nil {
		if ok := errors.Is(err, gorm.ErrRecordNotFound); ok {
			ctx.JSON(common.Failed(&api_error.BaseAPIError{ErrorCode: 10016, Message: "刷新token失败,用户不能存在"}))
		} else {
			ctx.JSON(common.Failed(api_error.SqlQueryError))
		}
		return
	}
	// 生成新的token对
	if result, err = generateToken(&user); err != nil {
		ctx.JSON(common.Failed(err))
		return
	}

	ctx.JSON(common.Success(result, 0, 0))
}

// Logout 注销登录
// 在配置了Blocklist时使用, 未配置情况下注销只需要前端清除token
func Logout(ctx iris.Context) {
	var err error
	if err = ctx.Logout(); err != nil {
		ctx.JSON(common.Failed(api_error.UserLogoutFailedError))
	}

	ctx.JSON(common.Success(nil, 0, 0))
}
