### GORMPool 使用说明

1. 初始化
   ```go
   package main
   
   import `github.com/aid297/aid/gormPool`
   
   func main() {
   	dbSetting, err := gormPool.APP.DBSetting.New("db.yaml") // 读取配置文件
   	if err != nil {
   		panic(err) // 读取配置文件失败，程序无法继续运行
   	}
   	mysqlPool := gormPool.APP.MySQLPool.Once(dbSetting) // 创建 orm 连接池
   	db := mysqlPool.GetConn()
   
   	var ret any
   	if err = db.Find(&ret).Error; err != nil { // 进行数据库操作
   		panic(err) // 异常处理
   	}
   }
   ```

2. 配置文件
   ```yaml
   common:
     driver: "mysql"
     maxOpenConns: 100
     maxIdleConns: 20
     maxLifetime: 100
     maxIdleTime: 10
   ar-sql:
     database: "cbit_db"
     rws: false
     main:
       username: "yjz"
       password: "123123"
       host: 127.0.0.1
       port: 12344
      sources:
      replicas:
   mysql:
     database: "tbl_test"
     charset: "utf8mb4"
     collation: "utf8mb4_general_ci"
     rws: true
     main:
       username: "root"
       password: "root"
       host: 127.0.0.1
       port: 3308
     sources:
       conn1:
         username: "root"
         password: "root"
         host: 127.0.0.1
         port: 3308
       conn2:
         username: "root"
         password: "root"
         host: 127.0.0.1
         port: 3308
     replicas:
       conn3:
         username: "root"
         password: "root"
         host: 127.0.0.1
         port: 3308
       conn4:
         username: "root"
         password: "root"
         host: 127.0.0.1
         port: 3308
       conn5:
         username: "root"
         password: "root"
         host: 127.0.0.1
         port: 3308
   postgres:
     main:
       username: "postgres"
       password: "postgres"
       host: 127.0.0.1
       port: 5432
       database: "tbl_test"
       sslmode: "disable"
       timezone: "Asia/Shanghai"
   sql-server:
     main:
       username: "admin"
       password: "Admin@1234"
       host: 127.0.0.1
       port: 9930
       database: "tbl_test"
   ```

3. 高级用法：`Finder`
   ```golang
   package main
   
   import (
   	`github.com/aid297/aid/gormPool`
   	`gorm.io/gorm`
   )
   
   type (
   	UserModel struct {
   		ID       uint            `gorm:"type:bigint unsigned;primaryKey;autoIncrement"`
   		Username string          `gorm:"type:varchar(64);not null;unique;default:''"`
   		Password string          `gorm:"type:varchar(255);not null;default:''"`
   		Age      uint            `gorm:"type:bigint unsigned;not null;default:0"`
   		Articles []*ArticleModel `gorm:"foreignKey:AuthorID;references:ID"`
   	}
   
   	ArticleModel struct {
   		Title    string     `gorm:"type:varchar(255);not null;unique;default:''"`
   		AuthorID int        `gorm:"type:bigint unsigned;not null;default:0"`
   		Author   *UserModel `gorm:"foreignKey:AuthorID;references:ID"`
   	}
   
   	UserCondition struct {
   		Page         int      `json:"page"`
   		PageSize     int      `json:"page_ize"`
   		Username     string   `json:"username"`
   		Age          uint     `json:"age"`
   		OrderBy      []string `json:"order_by"`
   		ArticleTitle string   `json:"article_title"`
   	}
   )
   
   func main() {
   	dbSetting, err := gormPool.APP.DBSetting.New("db.yaml") // 读取配置文件
   	if err != nil {
   		panic(err) // 读取配置文件失败，程序无法继续运行
   	}
   	mysqlPool := gormPool.APP.MySQLPool.Once(dbSetting) // 创建 orm 连接池
   	db := mysqlPool.GetConn()
   
   	var (
   		users     []UserModel
   		condition = UserCondition{Username: "admin", Age: 18, Page: 1, PageSize: 100}
   	)
   
   	gormPool.APP.Finder.New(db.Model(UserModel{})).
   		WhenFunc(condition.ArticleTitle != "", func(db *gorm.DB) {
   			db.Preload("Articles").
   				Joins("articles", "articles.author_id = users.id").
   				Where("articles.title = ?", condition.ArticleTitle)
   		}). // 如果带有级联属性则使用 WhenFunc 来构建复杂的 SQL 条件
   		When(condition.Age > 20, "age > ?", condition.Age). // 通过条件判断是否要增加对应 SQL 的条件
   		WhenIn(condition.Username != "", "username", condition.Username). // 通过条件判断是否要增加对应 SQL 的条件 → IN 查询（更多查询： WhenNotIn、WhenBetween、WhenNotBetween、WhenLike、WhenLikeRight、WhenFunc）
   		TryPagination(condition.Page, condition.PageSize). // 尝试使用分页（当 page 和 pageSize 都大于 0 时，会尝试进行分页。结果也由简单的列表改为分页格式列表）
   		TryOrder(condition.OrderBy...). // 尝试进行排序
   		Find(&users) // 构建查询条件并执行查询
   }
   ```

4. 高级用法：`FinderCondition`
   ```go
   package main
   
   import (
   	`github.com/aid297/aid/gormPool`
   	`gorm.io/gorm`
   )
   
   type (
   	UserModel struct {
   		ID       uint            `gorm:"type:bigint unsigned;primaryKey;autoIncrement"`
   		Username string          `gorm:"type:varchar(64);not null;unique;default:''"`
   		Password string          `gorm:"type:varchar(255);not null;default:''"`
   		Age      uint            `gorm:"type:bigint unsigned;not null;default:0"`
   		Articles []*ArticleModel `gorm:"foreignKey:AuthorID;references:ID"`
   	}
   
   	ArticleModel struct {
   		Title    string     `gorm:"type:varchar(255);not null;unique;default:''"`
   		AuthorID int        `gorm:"type:bigint unsigned;not null;default:0"`
   		Author   *UserModel `gorm:"foreignKey:AuthorID;references:ID"`
   	}
   
   	UserRequest struct {
   		finderCondition *gormPool.FinderCondition
   	}
   )
   
   func main() {
   	var (
   		err       error
   		dbSetting *gormPool.DBSetting
   		mysqlPool *gormPool.MySQLPool
   		db        *gorm.DB
   		users     []UserModel
   		request   = &UserRequest{finderCondition: &gormPool.FinderCondition{}}
   	)
   	if dbSetting, err = gormPool.APP.DBSetting.New("db.yaml"); err != nil { // 读取配置文件
   		panic(err) // 读取配置文件失败，程序无法继续运行
   	}
   	mysqlPool = gormPool.APP.MySQLPool.Once(dbSetting) // 创建 orm 连接池
   	db = mysqlPool.GetConn()
   
   	if err = gormPool.APP.Finder.
   		New(db.Model(UserModel{})).
   		FindOnlyCondition(request.finderCondition, &users).
   		GetDB().
   		Error; err != nil {
   		panic(err) // 异常处理
   	}
   
   	// 参数格式：
   	// {
   	//    "page": 1,
   	//    "pageSize": 10,
   	//    "condition": {
   	//        "queries": [
   	//            {
   	//                "option": "and",
   	//                "conditions": [
   	//                    {
   	//                        "key": "env_kind",
   	//                        "operator": "=",
   	//                        "values": [
   	//                            "coderepo"
   	//                        ]
   	//                    },
   	//                    {
   	//                        "key": "project_uuid",
   	//                        "operator": "=",
   	//                        "values": [
   	//                            "1f0afe0f-b4fa-62c6-bebd-14361e47a020"
   	//                        ]
   	//                    },
   	//                    {
   	//                        "key": "operator_uuid",
   	//                        "operator": "!=",
   	//                        "values": [
   	//                            ""
   	//                        ]
   	//                    }
   	//                ]
   	//            }
   	//        ],
   	//        "orders": [
   	//            "id desc"
   	//        ]
   	//    }
   	// }
   }
   ```

   