package rbac

import (
	"errors"
	"iris_project_foundation/common/api_error"
	"iris_project_foundation/models"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

type RoleService struct {
	Ctx iris.Context
}

func NewRoleService(ctx iris.Context) *RoleService {
	return &RoleService{Ctx: ctx}
}

// GetPermissions 查询所有可用权限
func (s *RoleService) GetPermissions() ([]*Permission, error) {
	var (
		err   error
		perms []*Permission
	)
	if err = models.DB.Find(&perms).Error; err != nil {
		return nil, api_error.PermissionQueryError
	}
	return perms, nil
}

// CreateRole 新增角色
func (s *RoleService) CreateRole() error {
	var (
		err  error
		errs validator.ValidationErrors
		ok   bool
		role Role
	)

	err = s.Ctx.ReadJSON(&role)
	if errs, ok = err.(validator.ValidationErrors); ok {
		return role.ParseValidationErrors(errs)
	}

	if err = models.DB.Create(&role).Error; err != nil {
		s.Ctx.Application().Logger().Error(err)
		if strings.HasPrefix(err.Error(), "Error 1062: Duplicate entry") {
			return api_error.RoleDuplicateNotationError
		}
		return api_error.DataCreateFailedError
	}
	return nil
}

// DeleteRole 删除角色, 同时删除casbin中对应的policy
func (s *RoleService) DeleteRole(id int64) error {
	var (
		err    error
		ok     bool
		result *gorm.DB
		role   Role
	)

	if result = models.DB.Delete(&role, id); result.Error != nil {
		s.Ctx.Application().Logger().Error(result.Error)
		return api_error.DataDeleteFailedError
	}

	if result.RowsAffected == 0 {
		return api_error.ResourceNotFoundError
	}

	if ok, err = GEnforcer.DeleteRole(role.Notation); !ok {
		s.Ctx.Application().Logger().Error(err)
		return api_error.RoleDeleteError
	}

	return nil
}

type EditRoleForm struct {
	ID          int64   `json:"role" validate:"required"`
	Permissions []int64 `json:"permissions" validate:"required"`
}

func (s *RoleService) SetPermissionsForRole() error {
	var (
		err   error
		form  EditRoleForm
		errs  validator.ValidationErrors
		ok    bool
		role  Role
		perms []*Permission
		perm  *Permission
	)
	err = s.Ctx.ReadJSON(&form)
	if errs, ok = err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			switch e.StructField() {
			case "ID":
				return api_error.RoleIDMissError
			case "Permissions":
				return api_error.RolePermsMissError
			}
		}
	}

	// 获取角色对象
	if err = models.DB.First(&role, strconv.FormatInt(form.ID, 10)).Error; err != nil {
		return api_error.ResourceNotFoundError
	}
	// 清空角色拥有的权限
	if _, err = GEnforcer.DeletePermissionsForUser(role.Notation); err != nil {
		s.Ctx.Application().Logger().Error(err)
		return api_error.RoleAllocatePermError
	}

	if err = models.DB.Find(&perms, "id IN ?", form.Permissions).Error; err != nil {
		s.Ctx.Application().Logger().Error(err)
		return api_error.RoleAllocatePermError
	}
	for _, perm = range perms {
		_, _ = GEnforcer.AddPermissionForUser(role.Notation, perm.Object, perm.Action)
	}

	return nil
}

type UserRoleForm struct {
	User  int64   `json:"user" validate:"required"`
	Roles []int64 `json:"roles" validate:"required"`
}

func (s *RoleService) SetRoleForUser() error {
	var (
		err   error
		errs  validator.ValidationErrors
		ok    bool
		form  UserRoleForm
		roles []*Role
		user  models.User
	)

	err = s.Ctx.ReadJSON(&form)
	if errs, ok = err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			switch e.StructField() {
			case "User":
				return api_error.RoleUserMissError
			case "Roles":
				return api_error.RoleRolesMissError
			}
		}
	}

	if err = models.DB.First(&user, form.User).Error; err != nil {
		if ok := errors.Is(err, gorm.ErrRecordNotFound); ok {
			return api_error.ResourceNotFoundError
		} else {
			return api_error.SqlQueryError
		}
	}

	if err = models.DB.Find(&roles, "id IN ?", form.Roles).Error; err != nil {
		return api_error.SqlQueryError
	}

	if _, err = GEnforcer.DeleteUser(user.Username); err != nil {
		return api_error.RoleSetUserError
	}
	for _, r := range roles {
		_, _ = GEnforcer.AddRoleForUser(user.Username, r.Notation)
	}
	return nil
}

func (s *RoleService) GetRoleList() ([]*Role, error) {
	var (
		err   error
		roles []*Role
	)

	if err = models.DB.Find(&roles).Error; err != nil {
		return nil, api_error.SqlQueryError
	}

	return roles, nil
}

type RolePerms struct {
	Label       string        `json:"label"`
	Notation    string        `json:"notation"`
	Permissions []*Permission `json:"permissions"`
}

func (s *RoleService) GetRoleDetail(id int64) (*RolePerms, error) {
	var (
		err      error
		role     Role
		rowPerms [][]string
		perm     *Permission
		perms    []*Permission
		result   *RolePerms
	)

	if err = models.DB.First(&role, strconv.FormatInt(id, 10)).Error; err != nil {
		if ok := errors.Is(err, gorm.ErrRecordNotFound); ok {
			return nil, api_error.ResourceNotFoundError
		} else {
			return nil, api_error.SqlQueryError
		}
	}

	rowPerms = GEnforcer.GetPermissionsForUser(role.Notation)
	for _, r := range rowPerms {
		perm = new(Permission)
		if err = models.DB.First(perm, "object = ? AND action = ?", r[1], r[2]).Error; err != nil {
			continue
		}
		perms = append(perms, perm)
	}

	result = &RolePerms{
		Label:       role.Label,
		Notation:    role.Notation,
		Permissions: perms,
	}
	return result, nil
}

type UserRoles struct {
	Name     string  `json:"name"`
	UserName string  `json:"username"`
	Roles    []*Role `json:"roles"`
}

func (s *RoleService) GetUserRole(id int64) (*UserRoles, error) {
	var (
		err      error
		user     models.User
		rawRoles []string
		roles    []*Role
		result   *UserRoles
	)

	if err = models.DB.First(&user, id).Error; err != nil {
		if ok := errors.Is(err, gorm.ErrRecordNotFound); ok {
			return nil, api_error.ResourceNotFoundError
		} else {
			return nil, api_error.SqlQueryError
		}
	}

	if rawRoles, err = GEnforcer.GetRolesForUser(user.Username); err != nil {
		return nil, api_error.RoleForUserQueryError
	}

	if err = models.DB.Find(&roles, "notation IN ?", rawRoles).Error; err != nil {
		return nil, api_error.RoleForUserQueryError
	}

	result = &UserRoles{Name: user.Name, UserName: user.Name, Roles: roles}
	return result, nil
}
