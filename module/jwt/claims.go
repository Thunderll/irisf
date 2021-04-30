package jwt

import (
	"iris_project_foundation/models"
	"strconv"
)

type UserClaims struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	user *models.User
}

func (u *UserClaims) GetAuthorization() string {
	return "JWT"
}

func (u *UserClaims) GetID() string {
	return strconv.FormatInt(u.ID, 10)
}

func (u *UserClaims) GetUsername() string {
	var user models.User
	if u.ID != 0 && u.user == nil {
		models.DB.First(&user, strconv.FormatInt(u.ID, 10))
		u.user = &user
	}

	return u.user.Username
}
