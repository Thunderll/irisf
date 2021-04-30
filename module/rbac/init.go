package rbac

import (
	"iris_project_foundation/models"
	"log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

var GEnforcer *casbin.Enforcer

func InitRBAC() {
	var (
		err      error
		adapter  *gormadapter.Adapter
		enforcer *casbin.Enforcer
	)

	// 迁移权限和角色表
	_ = models.DB.AutoMigrate(&Permission{}, &Role{})

	if adapter, err = gormadapter.NewAdapterByDB(models.DB); err != nil {
		log.Fatalf("[ERR] casbin初始化错误.\n%v", err)
	}

	if enforcer, err = casbin.NewEnforcer("./module/rbac/rbac_model.conf", adapter, true); err != nil {
		log.Fatalf("[ERR] casbin初始化错误.\n%v", err)
	}

	_ = enforcer.LoadPolicy()

	GEnforcer = enforcer

	//GEnforcer.AddPolicy("manager", "/web-api/v1/web-user", "READ")
	//GEnforcer.AddPolicy("manager", "/web-api/v1/web-user", "POST")
	//GEnforcer.AddRoleForUser("admin", "manager")

}
