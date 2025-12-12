package rbac

import (
	"log"
	"testing"

	"github.com/google/uuid"

	"github.com/aid297/aid/gormPool"
)

func init() {
	pool := gormPool.MySqlPoolApp.Once(gormPool.APP.DBSetting.New("./db.yaml"))
	db := pool.GetConn()

	tr := APP.Edge.Once(TablePrefix("my_rbac"), DB(db))

	if err := tr.AutoMigrate(); err != nil {
		log.Fatal(err)
	}
}

func Test1(t *testing.T) {
	t.Run("创建Dot", func(t *testing.T) {
		if err := APP.Permission.New().Store(PermissionName("权限1"), PermissionIdentity("permission1")); err != nil {
			log.Fatal(err)
		}

		if err := APP.Permission.New().Store(PermissionName("权限2"), PermissionIdentity("permission2")); err != nil {
			log.Fatal(err)
		}
	})
}

func Test2(t *testing.T) {
	t.Run("创建Edge", func(t *testing.T) {
		if err := APP.Role.New().Store(RoleName("角色1")); err != nil {
			log.Fatal(err)
		}
	})
}

func Test3(t *testing.T) {
	t.Run("绑定Face", func(t *testing.T) {
		var edge *Role
		if err := APP.Role.New().DB().First(&edge).Error; err != nil {
			log.Fatal(err)
		}

		var dots []*Permission
		if err := APP.Permission.New().DB().Find(&dots).Error; err != nil {
			log.Fatal(err)
		}

		if err := APP.Group.New().BindPermissions(edge, dots); err != nil {
			log.Fatal(err)
		}
	})
}

func Test4(t *testing.T) {
	t.Run("绑定Edge", func(t *testing.T) {
		var edge *Role
		if err := APP.Role.New().DB().First(&edge).Error; err != nil {
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

		if err := APP.Edge.Bind(edge.UUID, intersections); err != nil {
			log.Fatal(err)
		}
	})
}

func Test5(t *testing.T) {
	t.Run("检查用户权限：通过", func(t *testing.T) {
		if ok, err := APP.Edge.CheckPermission("01f083b2-8ee4-646a-9282-1a4e34a9320f", "level1", "permission1"); err != nil {
			t.Fatal(err)
		} else {
			t.Log(ok)
		}
	})

	t.Run("检查用户权限：不通过", func(t *testing.T) {
		if ok, err := APP.Edge.CheckPermission("01f083b2-8ee4-646a-9282-1a4e34a9320f", "level1", "permission3"); err != nil {
			t.Fatal(err)
		} else {
			t.Log(ok)
		}
	})

	t.Run("检查用户权限：通过（必须符合某一边）", func(t *testing.T) {
		if ok, err := APP.Edge.CheckRolePermission("01f083b2-8ee4-646a-9282-1a4e34a9320f", "level1", "permission2", uuid.MustParse("01f083b2-5497-6ff0-8d7b-1a4e34a9320f")); err != nil {
			t.Fatal(err)
		} else {
			t.Log(ok)
		}
	})
}
