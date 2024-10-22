#!/bin/bash

# test_api.sh
# Ensure the script exits if any command fails
set -e

# Function to print status messages
print_status() {
  echo -e "\n=== $1 ==="
}

# 1. Create Users
print_status "Creating Users"

curl -s -X POST http://localhost:8080/users \
-H "Content-Type: application/json" \
-d '{
  "name": "Alice Smith",
  "email": "alice@example.com",
  "mobile_number": "+1234567890"
}' | jq .

curl -s -X POST http://localhost:8080/users \
-H "Content-Type: application/json" \
-d '{
  "name": "Bob Johnson",
  "email": "bob@example.com",
  "mobile_number": "+19876543210"
}' | jq .

curl -s -X POST http://localhost:8080/users \
-H "Content-Type: application/json" \
-d '{
  "name": "Carol Williams",
  "email": "carol@example.com",
  "mobile_number": "+11234567890"
}' | jq .

curl -s -X POST http://localhost:8080/users \
-H "Content-Type: application/json" \
-d '{
  "name": "Dave Brown",
  "email": "dave@example.com",
  "mobile_number": "+10987654321"
}' | jq .

print_status "Users Created Successfully"

# 2. Retrieve Users
print_status "Retrieving Users"

echo "Retrieve Alice by Email:"
curl -s -X GET "http://localhost:8080/users?identifier=alice@example.com" | jq .

echo "Retrieve Bob by Mobile Number:"
curl -s -X GET "http://localhost:8080/users?identifier=+19876543210" | jq .

echo "Retrieve Carol by Name:"
curl -s -X GET "http://localhost:8080/users?identifier=Carol%20Williams" | jq .

print_status "Users Retrieved Successfully"

# 3. Add Expenses
print_status "Adding Expenses"

echo "Adding Expense with Equal Split:"
curl -s -X POST http://localhost:8080/expenses \
-H "Content-Type: application/json" \
-d '{
  "description": "Lunch at Cafe",
  "amount": 3000,
  "created_by": "alice@example.com",
  "split_type": "Equal",
  "participants": ["bob@example.com", "carol@example.com"]
}' | jq .

echo "Adding Expense with Exact Split:"
curl -s -X POST http://localhost:8080/expenses \
-H "Content-Type: application/json" \
-d '{
  "description": "Shopping",
  "amount": 4299,
  "created_by": "alice@example.com",
  "split_type": "Exact",
  "participants": ["bob@example.com", "carol@example.com", "alice@example.com"],
  "split_details": {
    "bob@example.com": 799,
    "carol@example.com": 2000,
    "alice@example.com": 1500
  }
}' | jq .

echo "Adding Expense with Percentage Split:"
curl -s -X POST http://localhost:8080/expenses \
-H "Content-Type: application/json" \
-d '{
  "description": "Party",
  "amount": 4000,
  "created_by": "alice@example.com",
  "split_type": "Percentage",
  "participants": ["bob@example.com", "carol@example.com", "dave@example.com"],
  "split_details": {
    "alice@example.com": 50,
    "bob@example.com": 25,
    "carol@example.com": 25
  }
}' | jq .

print_status "Expenses Added Successfully"

# 4. Retrieve Expenses for a User
print_status "Retrieving Expenses for Alice"

curl -s -X GET "http://localhost:8080/expenses/user?identifier=alice@example.com" | jq .

print_status "Expenses Retrieved Successfully"

# 5. Retrieve All Expenses with Pagination
print_status "Retrieving All Expenses (Page 1, Limit 10)"

curl -s -X GET "http://localhost:8080/expenses?page=1&limit=10" | jq .

print_status "All Expenses Retrieved Successfully"

# 6. Download Balance Sheet
print_status "Downloading Balance Sheet"

curl -s -X GET http://localhost:8080/balancesheet/download -o balance_sheet.csv
echo "Balance Sheet downloaded as 'balance_sheet.csv'"

print_status "All Tests Completed Successfully"
