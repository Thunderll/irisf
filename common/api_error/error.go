package api_error

type BaseAPIError struct {
	ErrorCode int64
	Message   string
}

func (e *BaseAPIError) Error() string {
	return e.Message
}

var (
	UnknownError          = &BaseAPIError{10001, "未知错误"}
	ResourceNotFoundError = &BaseAPIError{10002, "资源未找到"}

	DataCreateFailedError  = &BaseAPIError{10003, "新增数据失败"}
	DataDeleteFailedError  = &BaseAPIError{10004, "删除数据失败"}
	DataUpdateFailedError  = &BaseAPIError{10005, "修改数据失败"}
	ImageUrlInvalidError   = &BaseAPIError{10006, "无效的图片路径"}
	PrimaryKeyInvalidError = &BaseAPIError{10007, "主键无效"}
	SqlQueryError          = &BaseAPIError{10008, "数据库查询错误"}
	QueryParamError        = &BaseAPIError{10009, "查询参数错误"}

	WechatAuthorizeError     = &BaseAPIError{10010, "微信认证请求错误"}
	UnauthorizedError        = &BaseAPIError{10011, "认证失败"}
	ForbiddenError           = &BaseAPIError{10012, "权限不足"}
	WebAuthFailedError       = &BaseAPIError{10013, "用户名或密码错误"}
	WebUserCreateError       = &BaseAPIError{10014, "用户注册失败"}
	UserLogoutFailedError    = &BaseAPIError{10015, "用户注销失败"}
	RefreshTokenError        = &BaseAPIError{10016, "刷新token失败"}
	RefreshTokenInvalidError = &BaseAPIError{10017, "refresh token无效"}

	PermissionQueryError       = &BaseAPIError{10020, "权限数据查询失败"}
	RoleLabelMissError         = &BaseAPIError{10021, "角色名称是必填项"}
	RoleNotationMissError      = &BaseAPIError{10022, "角色标识是必填项"}
	RoleDeleteError            = &BaseAPIError{10023, "角色删除失败"}
	RoleIDMissError            = &BaseAPIError{10024, "缺少要编辑的角色"}
	RolePermsMissError         = &BaseAPIError{10025, "缺少角色权限"}
	RoleAllocatePermError      = &BaseAPIError{10026, "为角色分配权限失败"}
	RoleDuplicateNotationError = &BaseAPIError{10027, "角色标识重复"}
	RoleUserMissError          = &BaseAPIError{10028, "用户是必选项"}
	RoleRolesMissError         = &BaseAPIError{10029, "角色是必选项"}
	RoleSetUserError           = &BaseAPIError{10030, "为用户分配角色失败"}
	RoleForUserQueryError      = &BaseAPIError{10031, "用户角色查询失败"}
)
