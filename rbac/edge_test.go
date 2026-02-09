package rbac

import (
	"log"
	"testing"

	"github.com/google/uuid"

	"github.com/aid297/aid/gormPool"
)

func init() {
	dbSetting, err := gormPool.APP.DBSetting.New("./db.yaml")
	if err != nil {
		log.Fatal(err)
	}
	pool := gormPool.MySqlPoolApp.Once(dbSetting)
	db := pool.GetConn()

	tr := APP.Edge.Once(TablePrefix("my_rbac"), DB(db))

	if err = tr.AutoMigrate(); err != nil {
		log.Fatal(err)
	}
}

func Test1(t *testing.T) {
	t.Run("创建 Permission", func(t *testing.T) {
		if err := APP.Permission.New().Store(PermissionName("权限1"), PermissionIdentity("permission1")); err != nil {
			log.Fatal(err)
		}

		if err := APP.Permission.New().Store(PermissionName("权限2"), PermissionIdentity("permission2")); err != nil {
			log.Fatal(err)
		}
	})
}

func Test2(t *testing.T) {
	t.Run("创建 Role", func(t *testing.T) {
		if err := APP.Role.New().Store(RoleName("角色1")); err != nil {
			log.Fatal(err)
		}
	})
}

func Test3(t *testing.T) {
	t.Run("绑定 PermissionGroup", func(t *testing.T) {
		var role *Role
		if err := APP.Role.New().DB().First(&role).Error; err != nil {
			log.Fatal(err)
		}

		var permissions []*Permission
		if err := APP.Permission.New().DB().Find(&permissions).Error; err != nil {
			log.Fatal(err)
		}

		if err := APP.Group.New().BindPermissions(role, permissions); err != nil {
			log.Fatal(err)
		}
	})
}

func Test4(t *testing.T) {
	t.Run("绑定 Role", func(t *testing.T) {
		var role *Role
		if err := APP.Role.New().DB().First(&role).Error; err != nil {
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

		if err := APP.Edge.Bind(role.UUID, intersections); err != nil {
			log.Fatal(err)
		}
	})
}

func Test5(t *testing.T) {
	t.Run("检查用户权限：通过", func(t *testing.T) {
		ok, err := APP.Edge.CheckPermission("01f083b2-8ee4-646a-9282-1a4e34a9320f", "level1", "permission1")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(ok)
	})

	t.Run("检查用户权限：不通过", func(t *testing.T) {
		ok, err := APP.Edge.CheckPermission("01f083b2-8ee4-646a-9282-1a4e34a9320f", "level1", "permission3")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(ok)
	})

	t.Run("检查用户权限：通过（必须符合某一边）", func(t *testing.T) {
		ok, err := APP.Edge.CheckRolePermission("01f083b2-8ee4-646a-9282-1a4e34a9320f", "level1", "permission2", uuid.MustParse("01f083b2-5497-6ff0-8d7b-1a4e34a9320f"))
		if err != nil {
			t.Fatal(err)
		}
		t.Log(ok)
	})
}
