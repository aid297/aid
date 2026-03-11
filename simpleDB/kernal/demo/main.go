package main

import (
	"fmt"

	"github.com/aid297/aid/simpleDB/kernal"
)

const demoDatabase = "demo_runtime"

func main() {
	cleanupAll()
	defer cleanupAll()

	run("Create User with UUID v7 Primary Key", demoUserCRUD)
	run("Batch Insert and Query", demoBatchInsertAndQuery)
	run("Delete by Condition", demoDelete)
}

func demoUserCRUD() error {
	var (
		err                error
		table              = "users"
		db                 *kernal.SimpleDB
		newUser, foundUser kernal.Row
		exists             bool
		users              []kernal.Row
		userByPK           kernal.Row
	)

	cleanupTable(table)
	defer cleanupTable(table)

	if db, err = kernal.New.DB(demoDatabase, table); err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	// Set custom config: UUID v7, cascade depth 4
	db.SetAttrs(
		kernal.UUIDVersion(7),       // uuid 版本（默认6）
		kernal.UUIDWithHyphen(true), // 是否带有“-”（默认）
		kernal.UUIDUpper(true),      // 是否使用大写字母（默认）
		kernal.CascadeMaxDepth(4),   // 级联查询最大深度（默认6，过大可能导致性能问题）
		kernal.MaxCPUCores(4),       // 最大CPU核心数（默认：0，表示不限制，如果超过实际CPU核心数则按照最大计算）
		kernal.MaxMemoryGB(16),      // 最大内存使用量（默认：0，表示不限制，如果超出最大内存则按照不限制计算）
	)
	fmt.Printf("Config: %+v\n", db.GetConfig())

	// 配置数据库表结构
	if err = db.Configure(kernal.TableSchema{
		Columns: []kernal.Column{
			{
				Name:          "id",
				Type:          "uuid:v7",
				PrimaryKey:    true,
				AutoIncrement: true,
			},
			{
				Name:        "createdAt",
				Type:        "timestamp",
				Required:    true,
				DefaultExpr: kernal.ColumnExprCurrentTimestamp,
			},
			{
				Name:         "updatedAt",
				Type:         "timestamp",
				Required:     true,
				DefaultExpr:  kernal.ColumnExprCurrentTimestamp,
				OnUpdateExpr: kernal.ColumnExprCurrentTimestamp,
			},
			{
				Name:      "username",
				Type:      "string",
				MaxLength: 64,
				Unique:    true,
				Required:  true,
			},
			{
				Name:      "nickname",
				Type:      "string",
				MaxLength: 64,
				Required:  true,
				Default:   "",
				Indexed:   true,
			},
			{
				Name:     "age",
				Type:     "int",
				Required: true,
				Default:  0,
			},
			{
				Name:     "gender",
				Type:     "bool",
				Required: true,
				Default:  false,
			},
		},
	}); err != nil {
		return err
	}

	fmt.Println("Schema configured successfully")

	// Insert a user
	if newUser, err = db.InsertRow(kernal.Row{
		"username": "alice_2026",
		"nickname": "Alice Wang",
		"age":      28,
		"gender":   true,
	}); err != nil {
		return err
	}

	fmt.Println("\nInserted user:")
	printJSON(newUser)

	// Read by condition (engine auto-plans index usage)
	fmt.Println("\nFindOne by username (auto index planning):")
	if foundUser, exists, err = db.FindOne(kernal.QueryCondition{
		Field:    "username",
		Operator: kernal.QueryOpEQ,
		Value:    "alice_2026",
	}); err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("user not found")
	}
	printJSON(foundUser)

	// Read by condition (indexed field)
	fmt.Println("\nFind by nickname (auto index planning):")
	if users, err = db.Find(kernal.QueryCondition{
		Field:    "nickname",
		Operator: kernal.QueryOpEQ,
		Value:    "Alice Wang",
	}); err != nil {
		return err
	}
	printJSON(users)

	// Read by primary key (UUID)
	fmt.Printf("\nFind by primary key (UUID: %v):\n", newUser["id"])
	if userByPK, exists, err = db.FindOne(kernal.QueryCondition{
		Field:    "id",
		Operator: kernal.QueryOpEQ,
		Value:    newUser["id"],
	}); err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("user not found by primary key")
	}
	printJSON(userByPK)

	return nil
}

func demoBatchInsertAndQuery() error {
	var (
		err                                             error
		table                                           = "users_batch"
		db                                              *kernal.SimpleDB
		insertedRows, allRows, bobRows, ageFilteredRows []kernal.Row
	)

	cleanupTable(table)
	defer cleanupTable(table)

	if db, err = kernal.New.DB(demoDatabase, table); err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	// Set config: UUID v7
	db.SetAttrs(kernal.UUIDVersion(7))

	// Configure schema
	if err = db.Configure(kernal.TableSchema{
		Columns: []kernal.Column{
			{
				Name:          "id",
				Type:          "uuid:v7",
				PrimaryKey:    true,
				AutoIncrement: true,
			},
			{
				Name:      "username",
				Type:      "string",
				MaxLength: 64,
				Unique:    true,
				Required:  true,
			},
			{
				Name:      "nickname",
				Type:      "string",
				MaxLength: 64,
				Required:  true,
				Default:   "",
				Indexed:   true,
			},
			{
				Name:     "age",
				Type:     "int",
				Required: true,
				Default:  0,
			},
			{
				Name:     "gender",
				Type:     "bool",
				Required: true,
				Default:  false,
			},
		},
	}); err != nil {
		return err
	}

	fmt.Println("Schema configured successfully")

	// Batch insert two users
	fmt.Println("\nBatch inserting 2 users:")
	rows := []kernal.Row{
		{
			"username": "bob_2026",
			"nickname": "Bob Smith",
			"age":      32,
			"gender":   true,
		},
		{
			"username": "carol_2026",
			"nickname": "Carol Davis",
			"age":      26,
			"gender":   false,
		},
	}

	if insertedRows, err = db.InsertRows(rows); err != nil {
		return err
	}

	fmt.Println("Inserted rows:")
	printJSON(insertedRows)

	// Query all rows (empty conditions)
	fmt.Println("\nQuery all rows (using Find with empty conditions):")
	if allRows, err = db.Find(); err != nil {
		return err
	}
	printJSON(allRows)

	// Query by indexed field (nickname)
	fmt.Println("\nQuery by nickname = 'Bob Smith':")
	if bobRows, err = db.Find(kernal.QueryCondition{
		Field:    "nickname",
		Operator: kernal.QueryOpEQ,
		Value:    "Bob Smith",
	}); err != nil {
		return err
	}
	printJSON(bobRows)

	// Query by condition (age >= 28)
	fmt.Println("\nQuery by condition (age >= 28):")
	if ageFilteredRows, err = db.Find(kernal.QueryCondition{
		Field:    "age",
		Operator: kernal.QueryOpGTE,
		Value:    28,
	}); err != nil {
		return err
	}
	printJSON(ageFilteredRows)

	return nil
}

func demoDelete() error {
	var (
		err                           error
		table                         = "users_delete"
		db                            *kernal.SimpleDB
		rows, insertedRows, allBefore []kernal.Row
		allAfter, finalUsers          []kernal.Row
	)

	cleanupTable(table)
	defer cleanupTable(table)

	if db, err = kernal.New.DB(demoDatabase, table); err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	// Configure schema
	if err = db.Configure(kernal.TableSchema{
		Columns: []kernal.Column{
			{
				Name:          "id",
				Type:          "uuid:v7",
				PrimaryKey:    true,
				AutoIncrement: true,
			},
			{
				Name:      "username",
				Type:      "string",
				MaxLength: 64,
				Unique:    true,
				Required:  true,
			},
			{
				Name:     "age",
				Type:     "int",
				Required: true,
				Default:  0,
			},
		},
	}); err != nil {
		return err
	}

	// Insert test data
	fmt.Println("Inserting 3 test users:")
	rows = []kernal.Row{
		{"username": "user_a", "age": 20},
		{"username": "user_b", "age": 25},
		{"username": "user_c", "age": 30},
	}
	if insertedRows, err = db.InsertRows(rows); err != nil {
		return err
	}
	printJSON(insertedRows)

	// Get a user ID for deletion
	fmt.Println("\nUsers before deletion:")
	if allBefore, err = db.Find(); err != nil {
		return err
	}
	printJSON(allBefore)

	// Delete by ID (using the first user's ID)
	deleteID := insertedRows[0]["id"]
	fmt.Printf("\nDeleting user with ID: %v\n", deleteID)
	deletedCount, err := db.RemoveByCondition(kernal.QueryCondition{
		Field:    "id",
		Operator: kernal.QueryOpEQ,
		Value:    deleteID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Deleted %d row(s)\n", deletedCount)

	// Query remaining users
	fmt.Println("\nUsers after deletion:")
	if allAfter, err = db.Find(); err != nil {
		return err
	}
	printJSON(allAfter)

	// Delete by condition (age >= 25)
	fmt.Println("\nDeleting users with age >= 25:")
	deletedCount, err = db.RemoveByCondition(kernal.QueryCondition{
		Field:    "age",
		Operator: kernal.QueryOpGTE,
		Value:    25,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Deleted %d row(s)\n", deletedCount)

	// Query final state
	fmt.Println("\nFinal users:")
	if finalUsers, err = db.Find(); err != nil {
		return err
	}
	if len(finalUsers) == 0 {
		fmt.Println("(empty)")
	} else {
		printJSON(finalUsers)
	}

	return nil
}
