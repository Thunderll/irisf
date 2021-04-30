package rbac

import (
	"iris_project_foundation/common/api_error"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Permission struct {
	Label       string `gorm:"not null;comment:权限名称" json:"label"`
	Description string `gorm:"default:'';comment:权限描述" json:"description"`
	Object      string `gorm:"not null;uniqueIndex:idx_perm;comment:权限实体" json:"obj"`
	Action      string `gorm:"not null;uniqueIndex:idx_perm;comment:权限行为" json:"act"`
	gorm.Model
}

type Role struct {
	gorm.Model
	Label    string `gorm:"not null;comment:角色名称" json:"label" validate:"required"`
	Notation string `gorm:"not null;unique;comment:角色标识" json:"notation" validate:"required"`
}

func (m *Role) ParseValidationErrors(errs validator.ValidationErrors) error {
	for _, err := range errs {
		switch err.StructField() {
		case "Label":
			return api_error.RoleLabelMissError
		case "Notation":
			return api_error.RoleNotationMissError
		}
	}
	return nil
}
