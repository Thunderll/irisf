package models

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"iris_project_foundation/config"
	"log"
	"time"
)

var DB *gorm.DB

func InitDatabase() {
	var (
		err   error
		dsn   string
		db    *gorm.DB
		sqlDB *sql.DB
	)
	dsn = fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.GConfig.Database.User,
		config.GConfig.Database.Password,
		config.GConfig.Database.Host,
		config.GConfig.Database.Name,
	)
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix:   setting.DatabaseSetting.TablePrefix, // 表名前缀，`User`表为`t_users`
			SingularTable: true, // 使用单数表名，启用该选项后，`User` 表将是`user`
			//NameReplacer:  strings.NewReplacer("CID", "Cid"),   // 在转为数据库名称之前，使用NameReplacer更改结构/字段名称。
		},
	})

	if err != nil {
		log.Fatalf("[ERR] 数据库连接失败！ api_error: %v", err)
	} else {
		log.Println("[INFO] 数据库连接成功！")
	}

	if sqlDB, err = db.DB(); err != nil {
		log.Printf("[ERR] 数据库连接失败!")
		log.Fatal(err)
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	//自动迁移数据库
	_ = db.AutoMigrate(&User{})
	//添加&删除数据库
	//db.Migrator().CreateTable(&BasicInfo{})
	//db.Migrator().DropTable(&NewHouseMap{})

	DB = db.Debug()
}
