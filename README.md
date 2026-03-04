Live
https://golang-server-crud-auth-learning.onrender.com/

AWS
https://haiyk2bz5hpupwy6hs24e52s3q0oxzxz.lambda-url.ap-south-1.on.aws/
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

## ☁️ Deployment: AWS Lambda & CI/CD

This project is configured to run on both local servers (like Railway/Render) and **AWS Lambda**.

### 1. How it works (Lambda Adapter)
In `main.go`, we use a **Lambda Adapter** that transforms AWS API Gateway events into standard HTTP requests that the **Gin** framework can understand.
*   **Local**: Starts a standard HTTP server on port 8080.
*   **Lambda**: Starts the `lambda.Start()` handler when it detects the `LAMBDA_TASK_ROOT` environment variable.

### 2. CI/CD Pipeline (GitHub Actions)
The file `.github/workflows/deploy.yml` automatically deploys your code to AWS Lambda whenever you push to the `main` branch.

#### 🛠️ Steps to enable:
1.  **On GitHub**: Go to your repository **Settings** > **Secrets and variables** > **Actions**.
2.  Add the following **Repository Secrets**:
    *   `AWS_ACCESS_KEY_ID`: Your AWS Access Key ID.
    *   `AWS_SECRET_ACCESS_KEY`: Your AWS Secret Access Key.
3.  Ensure your **AWS Lambda function** has the following environment variables set in the AWS Console:
    *   `DB_URL`: Your Neon DB connection string.
    *   `JWT_SECRET`: Your secret key for JWT.
    *   `REFRESH_SECRET`: Your secret key for Refresh tokens.

### 3. Deploying to AWS Lambda (Manual check)
The CI/CD pipeline performs these commands:
```bash
# Build for Linux (Lambda Runtime)
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap main.go

# Zip the binary
zip deployment.zip bootstrap

# Update Lambda function code
aws lambda update-function-code --function-name GoLang_Server_CRUD_AUTH --zip-file fileb://deployment.zip
```

---
**Learning Note**: 
- **Serverless**: Lambda functions only run when a request comes in, making them extremely cost-effective.
- **`bootstrap`**: In the new `provided.al2023` Lambda runtime, the binary must be named `bootstrap`.
