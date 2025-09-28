# **File Sharing Web App**

A full-stack web application for uploading, managing, and sharing files.
The application allows users to securely upload files, view their own files, and download them.
The backend is written in Go (Gin framework) with a MySQL database, while files are stored locally on disk.
Phase B integrates user authentication and ownership-based access control, ensuring that each user can manage only their own files.

## **Features**

* Upload images, PDFs, and text files.
* List and download uploaded files.
* User authentication using JWT.
* File ownership enforcement (users can only access their own files).
* Stores file metadata in MySQL (filesharing database).
* Files stored on local disk (./uploads) for efficient retrieval.
* RESTful APIs for integration with frontend or other services.

## **Database Schema**

**Database:** filesharing

### **Files Table**

```sql
CREATE TABLE IF NOT EXISTS files (
  id INT AUTO_INCREMENT PRIMARY KEY,
  owner_id INT NULL,
  original_name VARCHAR(255) NOT NULL,
  stored_name VARCHAR(255) NOT NULL,
  mime_type VARCHAR(100),
  size_bytes BIGINT,
  path VARCHAR(1024),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX (owner_id)
);
```

### **Users Table (Phase B)**

```sql
CREATE TABLE IF NOT EXISTS users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(100) UNIQUE NOT NULL,
  email VARCHAR(255),
  password_hash VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE files
ADD CONSTRAINT fk_files_owner
FOREIGN KEY (owner_id)
REFERENCES users(id) ON DELETE SET NULL;
```

## **Environment Variables**

Create a .env file in the root directory or export the following environment variables:

```
DB_USER=root
DB_PASS=yourpassword
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=filesharing
JWT_SECRET=your_jwt_secret
```

## **Installation & Setup**

1. **Clone the repository:**
   ```
   git clone [https://github.com/Soumil-2007/file-sharing-webApp.git](https://github.com/Soumil-2007/file-sharing-webApp.git)
   cd file-sharing-webApp
   ```

2. **Install Go dependencies:**
   ```
   go mod tidy
   ```

3. Setup MySQL Database:
   Ensure your MySQL server is running.
   ```
   # Connect to MySQL
   mysql -u root -p

   # Inside MySQL shell
   CREATE DATABASE filesharing;
   USE filesharing;

   # Run migration files
   source db/migrations/001_create_files_table.sql;
   source db/migrations/002_create_users_table.sql;
   ```

4. **Create uploads directory:**
   ```
   mkdir -p ./uploads
   ```

5. Set environment variables:
   Create a .env file as shown above, or run this command:
   ```
   export DB_USER=root
   export DB_PASS=yourpassword
   export DB_HOST=127.0.0.1
   export DB_PORT=3306
   export DB_NAME=filesharing
   export JWT_SECRET=your_jwt_secret
   ```

6. **Run the application:**
   ```
   go run ./cmd
   ```

## **REST API Endpoints**

### **Public APIs (Phase A)**

| Method | Endpoint | Description |
| :---- | :---- | :---- |
| POST | /api/files | Upload a file |
| GET | /api/files | List all files |
| GET | /api/files/:id | Download a file by ID |

### **Authenticated APIs (Phase B)**

All routes require an Authorization: Bearer <JWT> header.

| Method | Endpoint | Description |
| :---- | :---- | :---- |
| POST | /api/users/register | Register a new user |
| POST | /api/users/login | Login and get JWT token |
| POST | /api/files | Upload a file (owner_id set automatically) |
| GET | /api/files | List current user's files only |
| GET | /api/files/:id | Download only if owner |
| DELETE | /api/files/:id | Delete only if owner |

## **File Handling Notes**

* Uploaded files are stored in ./uploads.
* Each file is saved with a UUID prefix to prevent collisions.
* File metadata (original name, stored name, MIME type, size, path) is stored in MySQL.
* Access is restricted based on owner_id.

## **Authentication Flow**

1. User registers using /api/users/register.
2. User logs in via /api/users/login and receives a JWT token.
3. Include the token in the Authorization header for all protected endpoints.
4. Middleware verifies the JWT and extracts the user_id for ownership enforcement.

## **Folder Structure**

```
file-sharing-webApp/
├── cmd/         # main.go, router registration
├── configs/     # configuration files
├── db/          # migration SQL files
├── services/
│   ├── db/      # DB initialization
│   └── files/   # file handlers (upload, list, download, delete)
├── uploads/     # stored uploaded files
├── go.mod
└── README.md
```

## **Testing APIs**

### **Upload a file (authenticated)**

```
curl -H "Authorization: Bearer <JWT>" -F "file=@/path/to/file.pdf" http://localhost:8080/api/files
```

### **List files**

```
curl -H "Authorization: Bearer <JWT>" http://localhost:8080/api/files
```

### **Download a file**

```
curl -H "Authorization: Bearer <JWT>" http://localhost:8080/api/files/5 -o downloaded.pdf
```

### **Delete a file**

```
curl -X DELETE -H "Authorization: Bearer <JWT>" http://localhost:8080/api/files/5
```

## **Future Enhancements**

* Implement file versioning.
* Add pagination and search for large file lists.
* Store files on cloud storage (S3) instead of local disk.
* Rate limiting and file quotas per user.
* Frontend integration with React or Vue.js.