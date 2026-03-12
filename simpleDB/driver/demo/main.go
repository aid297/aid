package main

import (
	"errors"

	"github.com/aid297/aid/debugLogger"
	"github.com/aid297/aid/simpleDB/driver"
)

func main() {
	db, err := driver.New.DB("demo", "driver_demo")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	baseSchema := driver.TableSchema{Columns: []driver.Column{
		{Name: "id", Type: "uuid:v7", PrimaryKey: true, AutoIncrement: true},
		{Name: "username", Type: "string", Unique: true, Required: true},
		{Name: "age", Type: "int", Required: true, Default: int64(0)},
	}}

	if err = db.AutoMigrate(baseSchema); err != nil {
		panic(err)
	}

	// 模拟 schema 演进：新增 email 列（普通索引）
	evolvedSchema := driver.TableSchema{Columns: []driver.Column{
		{Name: "id", Type: "uuid:v7", PrimaryKey: true, AutoIncrement: true},
		{Name: "username", Type: "string", Unique: true, Required: true},
		{Name: "age", Type: "int", Required: true, Default: int64(0)},
		{Name: "email", Type: "string", Indexed: true, Default: ""},
	}}

	plan, exists, err := db.SchemaDiff(evolvedSchema)
	if err != nil {
		panic(err)
	}
	debugLogger.Print("schema exists=%v, diff=%+v\n", exists, plan)

	if err = db.SyncSchema(evolvedSchema); err != nil {
		panic(err)
	}

	// 使用 AlterTable 做显式 DDL：再新增一个 nickname 列
	if err = db.AlterTable(driver.AlterTablePlan{
		AddColumns: []driver.Column{{Name: "nickname", Type: "string", Default: ""}},
		AddIndexes: []string{"nickname"},
	}); err != nil {
		panic(err)
	}

	if err = db.WithTx(func(tx *driver.Tx) error {
		_, err := tx.InsertRow(driver.Row{"username": "demo_user", "age": 24, "email": "demo_user@aid.dev", "nickname": "demo"})
		return err
	}); err != nil {
		panic(err)
	}

	_, err = db.InsertRow(driver.Row{"username": "demo_user", "age": 25, "email": "duplicate@aid.dev", "nickname": "dup"})
	if err != nil {
		var dErr *driver.DriverError
		if errors.As(err, &dErr) {
			debugLogger.Print("driver error code=%s err=%v\n", dErr.Code, dErr.Err)
		}
	}

	rows, err := db.Find(driver.QueryCondition{Field: "age", Operator: driver.QueryOpGTE, Value: 20})
	if err != nil {
		panic(err)
	}
	debugLogger.Print("rows=%v\n", rows)

	if err = demoLinkLayer(); err != nil {
		panic(err)
	}
}

func demoLinkLayer() error {
	users, err := driver.New.DB("demo", "driver_users")
	if err != nil {
		return err
	}
	defer users.Close()

	orders, err := driver.New.DB("demo", "driver_orders")
	if err != nil {
		return err
	}
	defer orders.Close()

	if err = users.SyncSchema(driver.TableSchema{Columns: []driver.Column{
		{Name: "id", Type: "uuid:v7", PrimaryKey: true, AutoIncrement: true},
		{Name: "username", Type: "string", Unique: true, Required: true},
	}}); err != nil {
		return err
	}

	if err = orders.SyncSchema(driver.TableSchema{
		Columns: []driver.Column{
			{Name: "id", Type: "uuid:v7", PrimaryKey: true, AutoIncrement: true},
			{Name: "user_id", Type: "uuid:v7", Required: true},
			{Name: "amount", Type: "float", Required: true, Default: float64(0)},
		},
		ForeignKeys: []driver.ForeignKey{
			{Name: "fk_orders_user", Field: "user_id", RefTable: "driver_users", RefField: "id", Alias: "user"},
		},
	}); err != nil {
		return err
	}

	user, err := users.InsertRow(driver.Row{"username": "link_user"})
	if err != nil {
		return err
	}

	if _, err = orders.InsertRows([]driver.Row{
		{"user_id": user["id"], "amount": float64(19.9)},
		{"user_id": user["id"], "amount": float64(88.8)},
	}); err != nil {
		return err
	}

	cascadeRows, err := orders.QueryCascade(driver.CascadeQuery{
		Conditions: []driver.QueryCondition{{Field: "user_id", Operator: driver.QueryOpEQ, Value: user["id"]}},
		Includes:   []driver.CascadeInclude{{Table: "driver_users", Alias: "user"}},
	})
	if err != nil {
		return err
	}

	debugLogger.Print("link-layer cascade rows=%v\n", cascadeRows)
	return nil
}
