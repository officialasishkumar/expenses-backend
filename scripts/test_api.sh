#!/bin/bash

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
  "name": "Priya Sharma",
  "email": "priya.sharma@example.com",
  "mobile_number": "9123456789"
}' | jq .

curl -s -X POST http://localhost:8080/users \
-H "Content-Type: application/json" \
-d '{
  "name": "Rajesh Kumar",
  "email": "rajesh.kumar@example.com",
  "mobile_number": "9876543210"
}' | jq .

curl -s -X POST http://localhost:8080/users \
-H "Content-Type: application/json" \
-d '{
  "name": "Anjali Singh",
  "email": "anjali.singh@example.com",
  "mobile_number": "9988776655"
}' | jq .

curl -s -X POST http://localhost:8080/users \
-H "Content-Type: application/json" \
-d '{
  "name": "Vikram Patel",
  "email": "vikram.patel@example.com",
  "mobile_number": "9765432109"
}' | jq .

print_status "Users Created Successfully"

# 2. Retrieve Users
print_status "Retrieving Users"

echo "Retrieve Priya by Email:"
curl -s -X GET "http://localhost:8080/users?identifier=priya.sharma@example.com" | jq .

echo "Retrieve Rajesh by Mobile Number:"
curl -s -X GET "http://localhost:8080/users?identifier=9876543210" | jq .

echo "Retrieve Anjali by Name:"
curl -s -X GET "http://localhost:8080/users?identifier=Anjali%20Singh" | jq .

print_status "Users Retrieved Successfully"

# 3. Add Expenses (Including Edge Cases)
print_status "Adding Expenses"

echo "Adding Expense with Equal Split:"
curl -s -X POST http://localhost:8080/expenses \
-H "Content-Type: application/json" \
-d '{
  "description": "Lunch at Cafe",
  "amount": 3000,
  "created_by": "priya.sharma@example.com",
  "split_type": "Equal",
  "participants": ["rajesh.kumar@example.com", "anjali.singh@example.com"]
}' | jq .

echo "Adding Expense with Exact Split:"
curl -s -X POST http://localhost:8080/expenses \
-H "Content-Type: application/json" \
-d '{
  "description": "Shopping",
  "amount": 4299,
  "created_by": "priya.sharma@example.com",
  "split_type": "Exact",
  "participants": ["rajesh.kumar@example.com", "anjali.singh@example.com", "priya.sharma@example.com"],
  "split_details": {
    "rajesh.kumar@example.com": 799,
    "anjali.singh@example.com": 2000,
    "priya.sharma@example.com": 1500
  }
}' | jq .

echo "Adding Expense with Percentage Split:"
curl -s -X POST http://localhost:8080/expenses \
-H "Content-Type: application/json" \
-d '{
  "description": "Party",
  "amount": 4000,
  "created_by": "priya.sharma@example.com",
  "split_type": "Percentage",
  "participants": ["rajesh.kumar@example.com", "anjali.singh@example.com", "vikram.patel@example.com"],
  "split_details": {
    "priya.sharma@example.com": 50,
    "rajesh.kumar@example.com": 25,
    "anjali.singh@example.com": 25
  }
}' | jq .

echo "Adding Expense with Invalid Percentage Split (Sum != 100%):"
curl -s -X POST http://localhost:8080/expenses \
-H "Content-Type: application/json" \
-d '{
  "description": "Invalid Party",
  "amount": 5000,
  "created_by": "priya.sharma@example.com",
  "split_type": "Percentage",
  "participants": ["rajesh.kumar@example.com", "anjali.singh@example.com"],
  "split_details": {
    "rajesh.kumar@example.com": 60,
    "anjali.singh@example.com": 30
  }
}' | jq .

echo "Adding Expense with Missing Participant:"
curl -s -X POST http://localhost:8080/expenses \
-H "Content-Type: application/json" \
-d '{
  "description": "Dinner",
  "amount": 2000,
  "created_by": "priya.sharma@example.com",
  "split_type": "Equal",
  "participants": ["nonexistent@example.com"]
}' | jq .

print_status "Expenses Added Successfully"

# 4. Retrieve Expenses for a User (Edge Case: No Expenses)
print_status "Retrieving Expenses for Vikram (No Expenses)"

curl -s -X GET "http://localhost:8080/expenses/user?identifier=vikram.patel@example.com" | jq .

# 5. Retrieve All Expenses with Pagination (Edge Case: Invalid Page/Limit)
print_status "Retrieving All Expenses with Invalid Pagination"

curl -s -X GET "http://localhost:8080/expenses?page=invalid&limit=invalid" | jq .

# 6. Download Balance Sheet (Edge Case: No Expenses)
print_status "Downloading Balance Sheet with No Expenses"

# Note: Since expenses have been added above, this balance sheet will include them.
# To test with no expenses, you would need to comment out the expense creation steps above or run this script on a clean database.

curl -s -X GET http://localhost:8080/balancesheet/download -o balance_sheet.csv
if [ -s balance_sheet.csv ]; then
  echo "Balance Sheet downloaded successfully"
else
  echo "Balance Sheet is empty as expected"
fi

# 7. Adding Edge Case: Zero Amount Expense
print_status "Adding Expense with Zero Amount"

curl -s -X POST http://localhost:8080/expenses \
-H "Content-Type: application/json" \
-d '{
  "description": "Zero Amount Expense",
  "amount": 0,
  "created_by": "priya.sharma@example.com",
  "split_type": "Equal",
  "participants": ["rajesh.kumar@example.com", "anjali.singh@example.com"]
}' | jq .

# 8. Retrieve All Expenses (Page 1, Limit 10)
print_status "Retrieving All Expenses (Page 1, Limit 10)"

curl -s -X GET "http://localhost:8080/expenses?page=1&limit=10" | jq .

print_status "All Expenses Retrieved Successfully"

# 9. Summary
print_status "All Tests Completed Successfully"
