package rbac_test

import (
	"log"
	"testing"

	"github.com/google/uuid"

	"github.com/aid297/aid/v2/gormPool"
	"github.com/aid297/aid/v2/gormPool/mysqlPool"
	"github.com/aid297/aid/v2/rbac"
)

func init() {
	dbSetting, err := gormPool.APP.DBSetting.New("./db.yaml")
	if err != nil {
		log.Fatal(err)
	}
	pool, err := mysqlPool.New(
		dbSetting.MySQL.Database,
		dbSetting.MySQL.Charset,
		dbSetting.MySQL.Rws,
		mysqlPool.Username(dbSetting.MySQL.Main.Username),
		mysqlPool.Password(dbSetting.MySQL.Main.Password),
		mysqlPool.Host(dbSetting.MySQL.Main.Host),
		mysqlPool.Port(dbSetting.MySQL.Main.Port),
		mysqlPool.Sources(dbSetting.MySQL.Sources),
		mysqlPool.Replicas(dbSetting.MySQL.Replicas),
		mysqlPool.MaxIdleTime(dbSetting.Common.MaxIdleTime),
		mysqlPool.MaxLifetime(dbSetting.Common.MaxLifetime),
		mysqlPool.MaxIdleConnections(dbSetting.Common.MaxIdleConnections),
		mysqlPool.MaxOpenConnections(dbSetting.Common.MaxOpenConnections),
	)
	if err != nil {
		log.Fatalf("创建mysql连接池错误：%v", err)
	}
	db := pool.GetConn()

	tr := rbac.APP.Edge.Once(rbac.TablePrefix("my_rbac"), rbac.DB(db))

	if err = tr.AutoMigrate(); err != nil {
		log.Fatal(err)
	}
}

func Test1(t *testing.T) {
	t.Run("创建 Permission", func(t *testing.T) {
		if err := rbac.APP.Permission.New().Store(rbac.PermissionName("权限1"), rbac.PermissionIdentity("permission1")); err != nil {
			log.Fatal(err)
		}

		if err := rbac.APP.Permission.New().Store(rbac.PermissionName("权限2"), rbac.PermissionIdentity("permission2")); err != nil {
			log.Fatal(err)
		}
	})
}

func Test2(t *testing.T) {
	t.Run("创建 Role", func(t *testing.T) {
		if err := rbac.APP.Role.New().Store(rbac.RoleName("角色1")); err != nil {
			log.Fatal(err)
		}
	})
}

func Test3(t *testing.T) {
	t.Run("绑定 PermissionGroup", func(t *testing.T) {
		var role *rbac.Role
		if err := rbac.APP.Role.New().DB().First(&role).Error; err != nil {
			log.Fatal(err)
		}

		var permissions []*rbac.Permission
		if err := rbac.APP.Permission.New().DB().Find(&permissions).Error; err != nil {
			log.Fatal(err)
		}

		if err := rbac.APP.Group.New().BindPermissions(role, permissions); err != nil {
			log.Fatal(err)
		}
	})
}

func Test4(t *testing.T) {
	t.Run("绑定 Role", func(t *testing.T) {
		var role *rbac.Role
		if err := rbac.APP.Role.New().DB().First(&role).Error; err != nil {
			log.Fatal(err)
		}

		userUUIDs := []string{
			uuid.Must(uuid.NewV6()).String(),
			uuid.Must(uuid.NewV6()).String(),
		}
		levels := []string{"level1", "level2"}

		intersections := map[string][]string{}
		for idx := range userUUIDs {
			intersections[userUUIDs[idx]] = append([]string{}, levels...)
		}

		if err := rbac.APP.Edge.Bind(role.UUID, intersections); err != nil {
			log.Fatal(err)
		}
	})
}

func Test5(t *testing.T) {
	t.Run("检查用户权限：通过", func(t *testing.T) {
		ok, err := rbac.APP.Edge.CheckPermission("01f083b2-8ee4-646a-9282-1a4e34a9320f", "level1", "permission1")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(ok)
	})

	t.Run("检查用户权限：不通过", func(t *testing.T) {
		ok, err := rbac.APP.Edge.CheckPermission("01f083b2-8ee4-646a-9282-1a4e34a9320f", "level1", "permission3")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(ok)
	})

	t.Run("检查用户权限：通过（必须符合某一边）", func(t *testing.T) {
		ok, err := rbac.APP.Edge.CheckRolePermission("01f083b2-8ee4-646a-9282-1a4e34a9320f", "level1", "permission2", uuid.MustParse("01f083b2-5497-6ff0-8d7b-1a4e34a9320f"))
		if err != nil {
			t.Fatal(err)
		}
		t.Log(ok)
	})
}
