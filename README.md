# Asynchronous File-Based CRUD API

[![Go](https://img.shields.io/badge/Go-1.24.2-blue)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-ready-green)](https://www.docker.com/)

A **Go-based asynchronous CRUD API** for file-based data management.  
The API allows concurrent requests to **Create, Read, Update, and Delete records** in a JSON file while maintaining data consistency.

---

## 🚀 Quick Start

### 1️⃣ Build the Docker image

```bash
docker build -t zabbix-technical-task .
```
This will build docker image from golang:1.24 and install all dependencies

---

### 2️⃣ Run the container
```bash
docker run -p 8080:8080 --rm file-crud-api
```
The server will listen on port 8080.

----

### 3️⃣ Test it in your browser
Open:
```bash
http://localhost:8080
```
Create a record
```bash
POST /records
Content-Type: application/json
Body: {
  "id": 123,
  "name": "Alice",
  "likes": ["apples", "bananas"]
  ... //other fields
}
```
Read a record (Replace :id with the numeric record ID)
```bash
GET /records/:id
```
Update a record (:id in path must match id in JSON body)
```bash
PUT /records/:id
Content-Type: application/json
Body: {
  "id": 123, 
  ... //updated fields
}
```
Delete a record (Deletes the record with the given ID)
```bash
DELETE /records/:id
```
---
### ⚙️Optional: Configure max unbacked records
```bash
const maxUnbackedRecords = 49
```
If 0, then each record goes straight to disk and nothing will be lost in case of a crash.
And if 49, then you can lose a maximum of 50 in case of a sudden crash.
### 🏗️ Project Structure
```
├── cmd/server         # HTTP server entry
├── internal/handler   # requests handler
├── internal/router    # requests multiplexer
├── pkg/cache/         # Cache implementation
├── pkg/storage/       # File storage
├── pkg/userrecord/    # Records implementation
├── go.mod
└── README.md
```