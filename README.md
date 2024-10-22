# Expenses Backend

This project is a backend service for managing users and expenses, with support for splitting expenses by equal, exact, or percentage methods. It provides APIs for user and expense management, as well as generating downloadable balance sheets.

## Prerequisites

Make sure the following are installed on your system:

- **Go** (Version 1.20 or later): [Installation Guide](https://golang.org/doc/install)
- **MongoDB**: [Installation Guide](https://docs.mongodb.com/manual/installation/)
- **jq** (JSON Processor): 
  - **Ubuntu/Debian**: `sudo apt-get install jq`
  - **Fedora**: `sudo dnf install jq`
  - **macOS**: `brew install jq`

Ensure the MongoDB service is running on your machine (default: `mongodb://localhost:27017`).

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

The server should now be running on `http://localhost:8080`.

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
- Create users and expenses
- Retrieve individual and all expenses
- Generate and download a balance sheet

If the balance sheet is successfully downloaded, it will be saved as `balance_sheet.csv`.

---

## Troubleshooting

- **MongoDB Connection Error**: Ensure that MongoDB is running and accessible at `mongodb://localhost:27017`.
- **jq Not Found**: Make sure `jq` is installed. See the [installation guide](https://stedolan.github.io/jq/download/) for more details.
- **Server Not Starting**: Ensure all dependencies are installed and there are no syntax errors in your code.
