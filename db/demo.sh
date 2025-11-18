#!/bin/bash

echo "Creating users table..."
echo "CREATE TABLE users (id INT, name VARCHAR, age INT)" | go run .

echo ""
echo "Inserting data..."
echo "INSERT INTO users VALUES (1, 'Alice', 25)" | go run .
echo "INSERT INTO users VALUES (2, 'Bob', 30)" | go run .
echo "INSERT INTO users VALUES (3, 'Charlie', 35)" | go run .

echo ""
echo "Listing all data..."
echo "SELECT * FROM users" | go run .

echo ""
echo "Showing table schema..."
echo "DESC users" | go run .