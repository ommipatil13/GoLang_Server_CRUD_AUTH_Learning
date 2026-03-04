# 🔐 Industry Standard Auth Explanation

Hi! Since you're new to Go, here's a quick guide on why we used this specific architecture.

## 1. Why Access & Refresh Tokens?
- **Access Token (Short-lived)**: This is like a visitor's pass. It lasts for only 15 minutes. It's stored in memory or local storage. If someone steals it, they only have 15 minutes of access.
- **Refresh Token (Long-lived)**: This is like your primary key to the house. It lasts for 7 days. It's stored in a ****HttpOnly Cookie**.
- **HttpOnly Cookie**: This is a special cookie that **cannot be read by JavaScript**. This prevents **XSS (Cross-Site Scripting)** attacks where a hacker's script might try to steal your token.

## 2. Refresh Token Rotation (Security)
In this project, every time you use a refresh token to get a new access token, we also:
1. Revoke the old refresh token.
2. Issue a **brand new** refresh token.
3. Update the database with this new token.
**Why?** If a hacker steals your refresh token and uses it, you will notice because YOUR next login attempt would fail (as the token was already used). This is the gold standard for security.

## 3. Password Hashing (Bcrypt)
We **never** store passwords in plain text.
- **Bcrypt** adds a "Salt" (random data) and hashes it thousands of times.
- Even if the database is leaked, hackers cannot easily "reverse" the hash to find your real password.

## 4. GORM & Auto-Migration
GORM is an Object-Relational Mapper.
- Instead of writing `CREATE TABLE users (...)`, we define a Go struct `User`.
- `db.AutoMigrate(&User{})` automatically creates the table for us. This keeps the database in sync with our code.

## 5. Gin Framework
Gin is used for routing. It's extremely fast and provides:
- **Middleware**: To check if a user is logged in before allowing access to a route.
- **JSON Binding**: To easily convert JSON input into Go objects.

Happy learning Go! 🚀
