### RBAC 套件

1. 创建`Permission`
   ```go
   package main
   
   import (
   	`log`
   	`github.com/aid297/aid/gormPool`
   	`github.com/aid297/aid/rbac`
   )
   
   func init() {
   	dbSetting, err := gormPool.APP.DBSetting.New("./db.yaml")
   	if err != nil {
   		log.Fatal(err)
   	}
   	pool := gormPool.MySqlPoolApp.Once(dbSetting)
   	db := pool.GetConn()
   
   	tr := rbac.APP.Edge.Once(rbac.TablePrefix("my_rbac"), rbac.DB(db))
   
   	if err = tr.AutoMigrate(); err != nil {
   		log.Fatal(err)
   	}
   }
   
   func main() {
   	if err := rbac.APP.Permission.New().Store(rbac.PermissionName("权限1"), rbac.PermissionIdentity("permission1")); err != nil {
   		log.Fatal(err)
   	}
   
   	if err := rbac.APP.Permission.New().Store(rbac.PermissionName("权限2"), rbac.PermissionIdentity("permission2")); err != nil {
   		log.Fatal(err)
   	}
   }
   ```

2. 创建`Role`
   ```go
   package main
   
   import (
   	`log`
   	`github.com/aid297/aid/gormPool`
   	`github.com/aid297/aid/rbac`
   )
   
   func init() {
   	dbSetting, err := gormPool.APP.DBSetting.New("./db.yaml")
   	if err != nil {
   		log.Fatal(err)
   	}
   	pool := gormPool.MySqlPoolApp.Once(dbSetting)
   	db := pool.GetConn()
   
   	tr := rbac.APP.Edge.Once(rbac.TablePrefix("my_rbac"), rbac.DB(db))
   
   	if err = tr.AutoMigrate(); err != nil {
   		log.Fatal(err)
   	}
   }
   
   func main() {
   	if err := rbac.APP.Role.New().Store(rbac.RoleName("角色1")); err != nil {
   		log.Fatal(err)
   	}
   }
   ```

3. 绑定`Permission Group`
   ```go
   package main
   
   import (
   	`log`
   	`github.com/aid297/aid/gormPool`
   	`github.com/aid297/aid/rbac`
   )
   
   func init() {
   	dbSetting, err := gormPool.APP.DBSetting.New("./db.yaml")
   	if err != nil {
   		log.Fatal(err)
   	}
   	pool := gormPool.MySqlPoolApp.Once(dbSetting)
   	db := pool.GetConn()
   
   	tr := rbac.APP.Edge.Once(rbac.TablePrefix("my_rbac"), rbac.DB(db))
   
   	if err = tr.AutoMigrate(); err != nil {
   		log.Fatal(err)
   	}
   }
   
   func main() {
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
   }
   ```

4. 绑定`Role`
   ```go
   package main
   
   import (
   	`log`
   	`github.com/aid297/aid/gormPool`
   	`github.com/aid297/aid/rbac`
   	`github.com/gofrs/uuid/v5`
   )
   
   func init() {
   	dbSetting, err := gormPool.APP.DBSetting.New("./db.yaml")
   	if err != nil {
   		log.Fatal(err)
   	}
   	pool := gormPool.MySqlPoolApp.Once(dbSetting)
   	db := pool.GetConn()
   
   	tr := rbac.APP.Edge.Once(rbac.TablePrefix("my_rbac"), rbac.DB(db))
   
   	if err = tr.AutoMigrate(); err != nil {
   		log.Fatal(err)
   	}
   }
   
   func main() {
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
   }
   ```

5. 检查是否具备权限
   ```go
   package main
   
   import (
   	`log`
   
   	`github.com/aid297/aid/gormPool`
   	`github.com/aid297/aid/rbac`
   	`github.com/google/uuid`
   )
   
   func init() {
   	dbSetting, err := gormPool.APP.DBSetting.New("./db.yaml")
   	if err != nil {
   		log.Fatal(err)
   	}
   	pool := gormPool.MySqlPoolApp.Once(dbSetting)
   	db := pool.GetConn()
   
   	tr := rbac.APP.Edge.Once(rbac.TablePrefix("my_rbac"), rbac.DB(db))
   
   	if err = tr.AutoMigrate(); err != nil {
   		log.Fatal(err)
   	}
   }
   
   func main() {
   	// 检查用户权限：通过
   	func() {
   		ok, err := rbac.APP.Edge.CheckPermission("01f083b2-8ee4-646a-9282-1a4e34a9320f", "level1", "permission1")
   		if err != nil {
   			log.Fatal(err)
   		}
   		log.Println(ok)
   	}()
   
   	// 检查用户权限：不通过
   	func() {
   		ok, err := rbac.APP.Edge.CheckPermission("01f083b2-8ee4-646a-9282-1a4e34a9320f", "level1", "permission3")
   		if err != nil {
   			log.Fatal(err)
   		}
   		log.Println(ok)
   	}()
   
   	// 检查用户权限：通过（必须符合某一边）
   	func() {
   		ok, err := rbac.APP.Edge.CheckRolePermission("01f083b2-8ee4-646a-9282-1a4e34a9320f", "level1", "permission2", uuid.MustParse("01f083b2-5497-6ff0-8d7b-1a4e34a9320f"))
   		if err != nil {
   			log.Fatal(err)
   		}
   		log.Println(ok)
   	}()
   }
   ```

   