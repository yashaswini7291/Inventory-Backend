#  Inventory Management System - Backend API

This is a backend service to manage product inventory for a small business. It includes secure user login using JWT tokens and endpoints to add, update, and retrieve product data. The project uses Go (Golang), the Gin web framework, MongoDB for storage, and Swagger for API documentation.

---

##  Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/yashaswini7291/Inventory-Backend.git
cd Inventory-Backend

```
### 2. Install Go Modules

```bash
go mod tidy
```

### 3.  Set Up MongoDB

Ensure MongoDB is running locally, or update the MongoDB connection URI in database/db.go.

```bash
mongodb://localhost:27017
```

### 4.  Generate Swagger Docs 
Swagger docs are already generated in the docs/ folder. If you update the annotations, regenerate with:

```bash
swag init
```
### 5.  Run the Server

```bash
go run main.go
```

By default, the server runs at: http://localhost:8080

