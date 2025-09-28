# File Sharing Web App

A secure file-sharing web application built in Go, allowing users to upload, manage, and share files with JWT-based authentication and ownership enforcement. Files are stored locally on disk, with metadata in MySQL. Includes background cleanup for expired files and in-memory caching for performance.

## Features

- **User Authentication**: Register/login with email and password; JWT tokens for secure access.
- **File Operations**: Upload (images, PDFs, text), list, and download files; users can only access their own files.
- **Concurrency**: Goroutines for file processing and background tasks.
- **Caching**: In-memory caching for file metadata to reduce DB load.
- **Background Jobs**: Periodic cleanup of expired files using goroutines.
- **RESTful APIs**: Clean endpoints for integration.

## Tech Stack

- **Backend**: Go (Golang)
- **Framework**: Gorilla Mux for routing
- **Database**: MySQL
- **Authentication**: JWT (JSON Web Tokens)
- **Caching**: In-memory (map-based)
- **Other**: bcrypt for password hashing, UUID for file naming

## Installation & Setup

### Prerequisites
- Go 1.25.1 or higher
- MySQL server
- Git

### Steps
1. **Clone the Repository**:
   ```
   git clone https://github.com/Soumil-2007/file-sharing-webApp.git
   cd file-sharing-webApp
   ```

2. **Install Dependencies**:
   ```
   go mod tidy
   ```

3. **Setup MySQL Database**:
   - Start MySQL and create the database:
     ```
     mysql -u root -p
     CREATE DATABASE filesharing;
     CREATE USER 'web-user'@'localhost' IDENTIFIED BY 'password';
     GRANT ALL PRIVILEGES ON filesharing.* TO 'web-user'@'localhost';
     FLUSH PRIVILEGES;
     EXIT;
     ```
   - Run migrations:
     ```
     go run cmd/migrate/migrations/main.go up
     ```

4. **Create Uploads Directory**:
   ```
   mkdir -p ./uploads
   ```

5. **Configure Environment Variables**:
   Create a `.env` file in the root:
   ```
   DB_USER=web-user
   DB_PASSWORD=password
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_NAME=filesharing
   JWT_SECRET=your_strong_secret_here
   ```

6. **Run the Application**:
   ```
   go run cmd/main.go
   ```
   Server starts on `http://localhost:8080`.

## Database Schema

### Users Table
```sql
CREATE TABLE users (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  firstName VARCHAR(255) NOT NULL,
  lastName VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);
```

### Files Table
```sql
CREATE TABLE files (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  owner_id INT,
  original_name VARCHAR(255) NOT NULL,
  stored_name VARCHAR(255) NOT NULL,
  mime_type VARCHAR(100),
  size_bytes BIGINT,
  path VARCHAR(1024),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE SET NULL,
  INDEX (owner_id)
);
```

## API Endpoints

Base URL: `http://localhost:8080/api/v1`

### Public Endpoints
- **POST /register**: Register a user.
  - Body: `{"firstName": "John", "lastName": "Doe", "email": "john@example.com", "password": "password123"}`
  - Response: 201 Created
- **POST /login**: Login and get JWT.
  - Body: `{"email": "john@example.com", "password": "password123"}`
  - Response: `{"token": "<jwt_token>"}`

### Authenticated Endpoints (Require `Authorization: Bearer <token>`)
- **POST /files**: Upload a file.
  - Body: Form-data with `file` (e.g., image.pdf)
  - Response: File metadata JSON
- **GET /files**: List user's files (cached).
  - Response: Array of file objects
- **GET /files/{id}**: Download a file.
  - Response: File binary

## Usage Examples

### With cURL
- **Register**:
  ```
  curl -X POST http://localhost:8080/api/v1/register -H "Content-Type: application/json" -d '{"firstName":"John","lastName":"Doe","email":"john@example.com","password":"password123"}'
  ```
- **Login**:
  ```
  curl -X POST http://localhost:8080/api/v1/login -H "Content-Type: application/json" -d '{"email":"john@example.com","password":"password123"}'
  ```
- **Upload File** (replace `<token>`):
  ```
  curl -X POST http://localhost:8080/api/v1/files -H "Authorization: Bearer <token>" -F "file=@/path/to/file.pdf"
  ```
- **List Files**:
  ```
  curl -H "Authorization: Bearer <token>" http://localhost:8080/api/v1/files
  ```

### With Postman
1. Set method to POST, URL to `http://localhost:8080/api/v1/register`, body as raw JSON.
2. For login, get token from response.
3. For protected routes, add header `Authorization: Bearer <token>`.

## Folder Structure

```
file-sharing-webApp/
├── cmd/
│   ├── main.go              # Entry point with background cleanup
│   ├── api/api.go           # Router setup
│   └── migrate/migrations/  # DB migrations
├── configs/envs.go          # Environment config
├── db/db.go                 # DB connection
├── services/
│   ├── auth/                # JWT and password handling
│   ├── files/files.go       # File upload/list/download with caching
│   ├── middleware/auth.go   # JWT middleware
│   └── user/                # User routes and store
├── types/types.go           # Structs and interfaces
├── utils/utils.go           # Helpers
├── uploads/                 # Stored files
├── .env                     # Environment vars
├── go.mod                   # Dependencies
└── README.md
```

## Background Jobs & Caching

- **Cleanup**: Goroutine runs hourly to delete expired files from disk and DB.
- **Caching**: In-memory map caches file lists for 5 minutes to avoid repeated DB queries.

## Future Enhancements

- Add search by name/type/date.
- Implement file sharing with expiration.
- Integrate Redis for distributed caching.
- Add file versioning and quotas.
- Build a frontend (React/Vue).

For issues, check server logs or DB connections.

---

This README is accurate to the current project state. If you need modifications, let me know!