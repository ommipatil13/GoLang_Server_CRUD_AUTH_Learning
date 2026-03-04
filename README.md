Live
https://golang-server-crud-auth-learning.onrender.com/
--------------------------------------------------------------------------
# 🚀 Golang Auth API (Gin + GORM + Postgres)

An industry-standard authentication API built with Go, Gin, and GORM. 

## ✨ Features
- **User CRUD**: Create, Read, Update, and Delete user profiles.
- **Authentication**: JWT-based authentication with Access and Refresh tokens.
- **Security**: 
  - Password hashing via `bcrypt`.
  - **Refresh Token Rotation**: New refresh tokens are issued on every refresh.
  - **HttpOnly Cookies**: Refresh tokens are stored securely in cookies.
- **Live Reloading**: Hot reloading during development with `CompileDaemon`.
- **Database**: PostgreSQL (Neon Tech) with GORM Auto-migration.

## 🛠️ Tech Stack
- **Framework**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/)
- **Auth**: JWT (JSON Web Token)
- **Hashing**: Bcrypt

## 📦 Deep Dive: Packages & Commands

Since you are learning, here is a detailed breakdown of every tool we used and why.

### 1. The Packages (Libraries)
We used `go get` to install these. You can see them listed in your `go.mod` file.

| Package | Purpose | Why it's "Industry Standard" |
| :--- | :--- | :--- |
| **`gin-gonic/gin`** | Web Framework | The most popular and fastest web framework for Go. It handles routing and JSON perfectly. |
| **`gorm.io/gorm`** | ORM (Object Relational Mapper) | Lets you interact with your database using Go structs instead of writing raw SQL strings. |
| **`postgres` driver** | Database Driver | Specialized library that allows GORM to talk specifically to PostgreSQL (Neon). |
| **`golang-jwt/jwt`** | Security Tokens | The standard library for creating and verifying JSON Web Tokens. |
| **`crypto/bcrypt`** | Password Hashing | A secure, slow hashing algorithm that makes it nearly impossible for hackers to "guess" passwords. |
| **`joho/godotenv`** | Env Management | Loads variables from your `.env` file into your computer's memory so your code can use them. |

### 2. The Commands (CLI)
Here are the exact commands I ran to set this up for you:

*   **`go mod init golang-auth-api`**: 
    - Initializes your project. It creates the `go.mod` file (which is like `package.json` in Node.js).
*   **`go get -u <package-name>`**:
    - Downloads the library and adds it to your project. The `-u` flag ensures it gets the latest version.
*   **`go install github.com/githubnemo/CompileDaemon@latest`**:
    - This installs a tool on your **system** (not just the project). Think of this like `npm install -g nodemon`.
*   **`go mod tidy`**:
    - Cleans up your `go.mod` file, removing unused packages and adding missing ones. It's good practice to run this often.

### 3. Running with Hot Reload (Nodemon style)
To run your project so it restarts automatically on save:
```bash
CompileDaemon --build="go build -o server.exe main.go" --command="./server.exe"
```
- **Build Step**: `go build` compiles your code into a "Binary" (a fast machine code file).
- **Command Step**: This tells the daemon to execute that binary file once the build is successful.

---

## 📡 API Endpoints

### Public Routes
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/auth/register` | Create a new account |
| `POST` | `/auth/login` | Login and get Access Token (+ Refresh Cookie) |
| `POST` | `/auth/refresh` | Get new tokens using Refresh Cookie |

### Protected Routes (Authorization: Bearer <token>)
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/profile` | Get current user's profile |
| `PUT` | `/api/user/update` | Update user details |
| `DELETE`| `/api/user/delete` | Delete account |
| `POST` | `/api/logout` | Invalidate session |
| `GET` | `/api/users` | List all registered users |

---
**Learning Note**: 
- **JWT (Access Token)** is short-lived.
- **Refresh Token** is stored in an **HttpOnly Cookie** so that JavaScript cannot access it (protects against XSS).
- **Refresh Token Rotation** makes sure even if a refresh token is stolen, it becomes invalid as soon as it's used or replaced.
