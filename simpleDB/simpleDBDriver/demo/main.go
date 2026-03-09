package main

import (
	"fmt"

	"github.com/aid297/aid/simpleDB/simpleDBDriver"
)

const demoDatabase = "demo_runtime"

func main() {
	cleanupAll()
	defer cleanupAll()

	run("Create User with UUID v7 Primary Key", demoUserCRUD)
	run("Batch Insert and Query", demoBatchInsertAndQuery)
}

func demoUserCRUD() error {
	table := "users"
	cleanupTable(table)
	defer cleanupTable(table)

	db, err := simpleDBDriver.New.SimpleDB(demoDatabase, table)
	if err != nil {
		return err
	}
	defer db.Close()

	// Set custom config: UUID v7, cascade depth 4
	db.SetConfig(simpleDBDriver.DatabaseConfig{
		DefaultUUIDVersion:     7,
		DefaultCascadeMaxDepth: 4,
	})
	fmt.Printf("Config: %+v\n", db.GetConfig())

	// Configure schema with UUID v7 primary key (using uuid:v7 format)
	err = db.Configure(simpleDBDriver.TableSchema{
		Columns: []simpleDBDriver.Column{
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
	})
	if err != nil {
		return err
	}

	fmt.Println("Schema configured successfully")

	// Insert a user
	newUser, err := db.InsertRow(simpleDBDriver.Row{
		"username": "alice_2026",
		"nickname": "Alice Wang",
		"age":      28,
		"gender":   true,
	})
	if err != nil {
		return err
	}

	fmt.Println("\nInserted user:")
	printJSON(newUser)

	// Read by unique field (username)
	fmt.Println("\nFind by username (unique index):")
	foundUser, exists, err := db.FindByUnique("username", "alice_2026")
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("user not found")
	}
	printJSON(foundUser)

	// Read by indexed field (nickname)
	fmt.Println("\nFind by nickname (indexed):")
	users, err := db.FindByIndex("nickname", "Alice Wang")
	if err != nil {
		return err
	}
	printJSON(users)

	// Read by primary key (UUID)
	userID := newUser["id"]
	fmt.Printf("\nFind by primary key (UUID: %v):\n", userID)
	userByPK, exists, err := db.FindByUnique("id", userID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("user not found by primary key")
	}
	printJSON(userByPK)

	return nil
}

func demoBatchInsertAndQuery() error {
	table := "users_batch"
	cleanupTable(table)
	defer cleanupTable(table)

	db, err := simpleDBDriver.New.SimpleDB(demoDatabase, table)
	if err != nil {
		return err
	}
	defer db.Close()

	// Set config: UUID v7
	db.SetConfig(simpleDBDriver.DatabaseConfig{
		DefaultUUIDVersion: 7,
	})

	// Configure schema
	err = db.Configure(simpleDBDriver.TableSchema{
		Columns: []simpleDBDriver.Column{
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
	})
	if err != nil {
		return err
	}

	fmt.Println("Schema configured successfully")

	// Batch insert two users
	fmt.Println("\nBatch inserting 2 users:")
	rows := []simpleDBDriver.Row{
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

	insertedRows, err := db.InsertRows(rows)
	if err != nil {
		return err
	}

	fmt.Println("Inserted rows:")
	printJSON(insertedRows)

	// Query all rows (by scanning with empty condition)
	fmt.Println("\nQuery all rows (using FindByConditions with empty conditions):")
	allRows, err := db.FindByConditions([]simpleDBDriver.QueryCondition{})
	if err != nil {
		return err
	}
	printJSON(allRows)

	// Query by indexed field (nickname)
	fmt.Println("\nQuery by nickname = 'Bob Smith':")
	bobRows, err := db.FindByIndex("nickname", "Bob Smith")
	if err != nil {
		return err
	}
	printJSON(bobRows)

	// Query by condition (age >= 28)
	fmt.Println("\nQuery by condition (age >= 28):")
	ageFilteredRows, err := db.FindByConditions([]simpleDBDriver.QueryCondition{
		{
			Field:    "age",
			Operator: simpleDBDriver.QueryOpGTE,
			Value:    28,
		},
	})
	if err != nil {
		return err
	}
	printJSON(ageFilteredRows)

	return nil
}
