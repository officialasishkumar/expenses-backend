# Expenses Backend

This project is a backend service for managing users and expenses, with support for splitting expenses by equal, exact, or percentage methods. It provides APIs for user and expense management, as well as generating downloadable balance sheets.

video walkthrough: https://youtu.be/KIQ-vyIsIDc

## Prerequisites

Make sure the following are installed on your system:

* **Go** (Version 1.20 or later): [Installation Guide](https://golang.org/doc/install)
* **MongoDB**: [Installation Guide](https://docs.mongodb.com/manual/installation/)
* **jq** (JSON Processor): 
  + **Ubuntu/Debian**: `sudo apt-get install jq`
  + **Fedora**: `sudo dnf install jq`
  + **macOS**: `brew install jq`

Ensure the MongoDB service is running on your machine (default: `mongodb://localhost:27017` ).

## Project Setup

### 1. Clone the Repository

```bash
git clone https://github.com/officialasishkumar/expenses-backend
cd expenses-backend
```

### 2. Install Dependencies

Navigate to the project directory and install Go modules:

```bash
go mod tidy
```

### 3. Start MongoDB (if not already running)

On Linux or macOS:

```bash
sudo systemctl start mongod
```

On Windows (Command Prompt):

```bash
net start MongoDB
```

### 4. Run the Backend Server

Make sure the backend server is running before running the test scripts.

```bash
go run main.go
```

The server should now be running on `http://localhost:8080` .

---

## Running the `test_api.sh` Script

In case you would like to run the curl commands on your own, please have a look at scripts/test_api.sh file.

### 1. Make the Script Executable

If you haven't already, give the script execute permissions:

```bash
chmod +x test_api.sh
```

### 2. Run the Test Script

The script will send multiple API requests to the backend and print the results.

```bash
./test_api.sh
```

### 3. Verify Results

The script will:
* Create users and expenses
* Retrieve individual and all expenses
* Generate and download a balance sheet

If the balance sheet is successfully downloaded, it will be saved as `balance_sheet.csv` .

### 4. API Documentation

## User Endpoints

### **POST /users** – Create a New User

**Request Body:**

```json
{
  "name": "Priya Sharma",
  "email": "priya.sharma@example.com",
  "mobile_number": "9123456789"
}
```

**Response:**

* **201 Created** – Returns user details.  
* **400 Bad Request** – If validation fails.

---

### **GET /users** – Retrieve User Details

**Query Parameters:**  
One of the following must be provided:

* `identifier` (can be **email**, **mobile_number**, or **name**)

**Behavior:**  
* If `email` or `mobile_number` is provided, fetch the unique user.
* If `name` is provided:
  + If the name is **unique**, return the user details.
  + If **multiple users** exist with the same name, return an **error** prompting for email or phone number.

**Response:**

* **200 OK** – Returns user details.  
* **400 Bad Request** – If input is invalid or ambiguous.

---

## Expense Endpoints

### **POST /expenses** – Add a New Expense

**Request Body:**

```json
{
  "description": "Lunch at Cafe",
  "amount": 3000,
  "created_by": "priya.sharma@example.com",
  "split_type": "Equal",
  "participants": ["rajesh.kumar@example.com", "anjali.singh@example.com"]
}
```

**Behavior:**  
* Identify `created_by` and `participants` using **email**, **phone**, or **name**.  
* Validate split details based on the `split_type`.  

**Response:**

* **201 Created** – Returns expense details.  
* **400 Bad Request** – If validation fails.

---

### **GET /expenses/user** – Retrieve All Expenses for a User

**Query Parameter:**

* `identifier` (can be **email**, **phone**, or **name**)

**Response:**

* **200 OK** – Returns a list of expenses.  
* **400 Bad Request** – If the user is ambiguous or not found.

---

### **GET /expenses** – Retrieve Overall Expenses

**Optional Query Parameters:**  
* Pagination parameters like `page` and `limit`.

**Response:**

* **200 OK** – Returns a list of all expenses.

---

## Balance Sheet Endpoint

### **GET /balancesheet/download** – Download Balance Sheet

**Response:**

* **200 OK** – Provides a downloadable **CSV file**.  
* **400 Bad Request** – If generation fails.

---

## Troubleshooting

* **MongoDB Connection Error**: Ensure that MongoDB is running and accessible at `mongodb://localhost:27017`.
* **jq Not Found**: Make sure `jq` is installed. See the [installation guide](https://stedolan.github.io/jq/download/) for more details.
* **Server Not Starting**: Ensure all dependencies are installed and there are no syntax errors in your code.
