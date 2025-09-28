# **File Sharing Web App 🚀**

A full-stack web application for uploading, managing, and sharing files securely. This application allows users to upload files, view their own files, and download them. The backend is written in **Go** (using the Gin framework) with a **MySQL** database.

Phase B of the project integrates user authentication and ownership-based access control, ensuring that each user can manage only their own files.

## **✨ Features**

* **File Operations**: Upload images, PDFs, and text files.  
* **File Management**: List and download your uploaded files.  
* **Secure Authentication**: User registration and login using JSON Web Tokens (JWT).  
* **Ownership Enforcement**: Users can only access, manage, and delete their own files.  
* **Efficient Storage**: File metadata is stored in a MySQL database, while the files themselves are stored on the local disk for quick retrieval.  
* **RESTful API**: A clean and simple API for easy integration with frontend frameworks or other services.

## **🛠️ Tech Stack**

* **Backend**: Go (Golang)  
* **Framework**: Gin Web Framework  
* **Database**: MySQL  
* **Authentication**: JWT

## **⚙️ Installation & Setup**

Follow these steps to get the application up and running on your local machine.

### **1\. Prerequisites**

* [Go](https://go.dev/doc/install) (version 1.18 or higher)  
* \[suspicious link removed\]  
* [Git](https://git-scm.com/downloads/)

### **2\. Clone the Repository**

git clone \[https://github.com/Soumil-2007/file-sharing-webApp.git\](https://github.com/Soumil-2007/file-sharing-webApp.git)  
cd file-sharing-webApp

### **3\. Install Dependencies**

go mod tidy

### **4\. Setup MySQL Database**

First, ensure your MySQL server is running. Then, connect to the MySQL shell and run the following commands:

\-- Connect to MySQL  
mysql \-u root \-p

\-- Create the database  
CREATE DATABASE filesharing;  
USE filesharing;

\-- Run the migration scripts to create tables  
SOURCE db/migrations/001\_create\_files\_table.sql;  
SOURCE db/migrations/002\_create\_users\_table.sql;

### **5\. Create Uploads Directory**

This directory will store all the uploaded files.

mkdir \-p ./uploads

### **6\. Configure Environment Variables**

Create a .env file in the root of the project or export the following variables in your terminal:

DB\_USER=root  
DB\_PASS=yourpassword  
DB\_HOST=127.0.0.1  
DB\_PORT=3306  
DB\_NAME=filesharing  
JWT\_SECRET=your\_jwt\_secret

### **7\. Run the Application**

go run ./cmd

The server will start, typically on http://localhost:8080.

## **🗄️ Database Schema**

The application uses two main tables in the filesharing database.

#### **users Table**

Stores user credentials and information.

CREATE TABLE IF NOT EXISTS users (  
  id INT AUTO\_INCREMENT PRIMARY KEY,  
  username VARCHAR(100) UNIQUE NOT NULL,  
  email VARCHAR(255),  
  password\_hash VARCHAR(255) NOT NULL,  
  created\_at TIMESTAMP DEFAULT CURRENT\_TIMESTAMP  
);

#### **files Table**

Stores metadata about the uploaded files and links them to a user via owner\_id.

CREATE TABLE IF NOT EXISTS files (  
  id INT AUTO\_INCREMENT PRIMARY KEY,  
  owner\_id INT NULL,  
  original\_name VARCHAR(255) NOT NULL,  
  stored\_name VARCHAR(255) NOT NULL,  
  mime\_type VARCHAR(100),  
  size\_bytes BIGINT,  
  path VARCHAR(1024),  
  created\_at TIMESTAMP DEFAULT CURRENT\_TIMESTAMP,  
  INDEX (owner\_id),  
  FOREIGN KEY (owner\_id) REFERENCES users(id) ON DELETE SET NULL  
);

## **🔌 REST API Endpoints**

The API is divided into public and authenticated endpoints.

### **Authenticated APIs (Phase B)**

All routes require an Authorization: Bearer \<JWT\> header.

| Method | Endpoint | Description |
| :---- | :---- | :---- |
| POST | /api/users/register | Register a new user. |
| POST | /api/users/login | Login to get a JWT token. |
| POST | /api/files | Upload a file (owner is auto-assigned). |
| GET | /api/files | List files belonging to the current user. |
| GET | /api/files/:id | Download a file if you are the owner. |
| DELETE | /api/files/:id | Delete a file if you are the owner. |

### **Public APIs (Phase A \- Legacy)**

| Method | Endpoint | Description |
| :---- | :---- | :---- |
| POST | /api/files | Upload a file. |
| GET | /api/files | List all files. |
| GET | /api/files/:id | Download a file by ID. |

## **🔑 Authentication Flow**

1. **Register**: A new user signs up using the /api/users/register endpoint.  
2. **Login**: The user logs in via /api/users/login with their credentials and receives a JWT token in response.  
3. **Authorize**: For all protected endpoints, the user must include the JWT token in the Authorization header as a Bearer token.  
4. **Verify**: A middleware on the server intercepts the request, verifies the JWT, and extracts the user\_id to enforce file ownership rules.

## **📂 Folder Structure**

The project follows a modular structure for better organization and scalability.

file-sharing-webApp/  
├── cmd/                \# main.go, router registration  
├── configs/            \# configuration files  
├── db/                 \# migration SQL files  
├── services/  
│   ├── db/             \# DB initialization  
│   └── files/          \# file handlers (upload, list, download, delete)  
├── uploads/            \# stored uploaded files  
├── go.mod  
└── README.md

## **🧪 Testing APIs with cURL**

Here are some cURL commands to test the authenticated API endpoints. Replace \<JWT\> with the token you receive after logging in.

#### **Upload a file**

curl \-H "Authorization: Bearer \<JWT\>" \\  
     \-F "file=@/path/to/your/file.pdf" \\  
     http://localhost:8080/api/files

#### **List your files**

curl \-H "Authorization: Bearer \<JWT\>" http://localhost:8080/api/files

#### **Download a file (e.g., file with ID 5\)**

curl \-H "Authorization: Bearer \<JWT\>" \\  
     http://localhost:8080/api/files/5 \-o downloaded.pdf

#### **Delete a file (e.g., file with ID 5\)**

curl \-X DELETE \-H "Authorization: Bearer \<JWT\>" http://localhost:8080/api/files/5

## **🚀 Future Enhancements**

* **File Versioning**: Implement a system to track and manage different versions of a file.  
* **Pagination & Search**: Add pagination and search functionality for large lists of files.  
* **Cloud Storage**: Integrate with cloud storage providers like Amazon S3 or Google Cloud Storage.  
* **User Quotas**: Implement rate limiting and file storage quotas per user.  
* **Frontend Integration**: Develop a frontend client using a framework like React or Vue.js.