# 📝 Blog API - REST API dengan Go, PostgreSQL, dan JWT

REST API sederhana untuk aplikasi **Blog** dengan fitur **authentication**, **posts**, dan **comments**.

---

## 🚀 Fitur

* ✅ User Registration & Login (JWT Authentication)
* ✅ CRUD Posts (Create, Read, Update, Delete)
* ✅ CRUD Comments
* ✅ Validasi Input
* ✅ Error Handling Terpusat
* ✅ Database Transaction
* ✅ Soft Delete
* ✅ Dokumentasi Swagger/OpenAPI
* ✅ Unit Tests
* ✅ Docker & Docker Compose Support

---

## 🔧 Tech Stack

| Komponen     | Teknologi              |
| ------------ | ---------------------- |
| **Bahasa**   | Go 1.21+               |
| **Database** | PostgreSQL 15          |
| **ORM**      | GORM                   |
| **Router**   | Gorilla Mux            |
| **Auth**     | JWT (`golang-jwt/jwt`) |
| **Docs**     | Swagger/OpenAPI        |

---

## 📁 Struktur Project

```
blog-api/
├── cmd/
│   └── server/
│       └── main.go              # Entry point aplikasi
├── internal/
│   ├── config/
│   │   └── config.go            # Konfigurasi aplikasi
│   ├── database/
│   │   ├── database.go          # Koneksi database
│   │   └── migration.go         # Migration & seed
│   ├── models/
│   │   ├── user.go              # User model
│   │   ├── post.go              # Post model
│   │   └── comment.go           # Comment model
│   ├── handlers/
│   │   ├── auth.go              # Auth handlers
│   │   ├── post.go              # Post handlers
│   │   ├── comment.go           # Comment handlers
│   │   ├── health.go            # Health check
│   │   ├── swagger.go           # Swagger handlers
│   │   ├── validator.go         # Input validation
│   │   └── errors.go            # Error handling
│   └── middleware/
│       └── auth.go              # JWT middleware
├── docs/
│   └── docs.go                  # Swagger documentation
├── .env                         # Environment variables
├── .env.example                 # Environment template
├── Dockerfile                   # Docker image definition
├── docker-compose.yml           # Docker Compose config
├── go.mod                       # Go module
└── README.md                    # Dokumentasi
```

---

## 🔧 Setup & Instalasi

### ⚙️ Prasyarat

* Go 1.21 atau lebih tinggi
* PostgreSQL 15 (untuk lokal)
* Docker & Docker Compose (untuk mode container)

---

### 🧩 1. Clone Repository

```bash
git clone <repository-url>
cd blog-api
```

---

### 🧩 2. Setup Environment Variables

```bash
cp .env.example .env
# Edit .env sesuai kebutuhan
```

---

## 🐳 Menjalankan dengan Docker (Recommended)

Cara termudah untuk menjalankan aplikasi:

```bash
# Build dan jalankan
docker-compose up --build

# Jalankan di background
docker-compose up -d

# Lihat logs
docker-compose logs -f app

# Stop aplikasi
docker-compose down

# Reset database (hapus volume)
docker-compose down -v
```

API akan berjalan di **[http://localhost:8080](http://localhost:8080)**

---

## 💻 Menjalankan Secara Lokal (Tanpa Docker)

### 1. Install Dependencies

```bash
go mod download
```

### 2. Setup PostgreSQL

Pastikan PostgreSQL berjalan dan buat database:

```sql
CREATE DATABASE blogdb;
```

### 3. Konfigurasi `.env`

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=blogdb
```

### 4. Jalankan Aplikasi

```bash
go run ./cmd/server/main.go
```

> Migration dan seed data akan otomatis berjalan saat startup.

---

## 📝 API Documentation

Akses Swagger UI di:
👉 **[http://localhost:8080/swagger](http://localhost:8080/swagger)**

---

## 📚 Endpoints

| Method | Endpoint                                     | Auth | Deskripsi          |
| ------ | -------------------------------------------- | ---- | ------------------ |
| GET    | `/health`                                    | ❌    | Health check       |
| POST   | `/api/register`                              | ❌    | Register user baru |
| POST   | `/api/login`                                 | ❌    | Login user         |
| GET    | `/api/posts`                                 | ❌    | Get semua posts    |
| GET    | `/api/posts/{id}`                            | ❌    | Get post by ID     |
| POST   | `/api/posts`                                 | ✅    | Create post baru   |
| PUT    | `/api/posts/{id}`                            | ✅    | Update post        |
| DELETE | `/api/posts/{id}`                            | ✅    | Delete post        |
| GET    | `/api/posts/{post_id}/comments`              | ❌    | Get comments       |
| POST   | `/api/posts/{post_id}/comments`              | ✅    | Create comment     |
| DELETE | `/api/posts/{post_id}/comments/{comment_id}` | ✅    | Delete comment     |

---

## 🧪 Testing

### Jalankan Unit Tests

```bash
# Run all tests
go test ./...

# Run dengan verbose
go test -v ./internal/handlers/

# Run dengan coverage
go test -cover ./internal/handlers/

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 📦 Contoh Penggunaan API

### 1. Register User

```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe"
  }'
```

**Response:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2025-01-15T10:30:00Z"
  }
}
```

---

### 2. Login

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Response sama dengan register (token + user data).**

---

### 3. Create Post (Authenticated)

```bash
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "title": "My First Post",
    "content": "This is the content of my first blog post."
  }'
```

**Response:**

```json
{
  "id": 1,
  "title": "My First Post",
  "content": "This is the content of my first blog post.",
  "user_id": 1,
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe"
  },
  "created_at": "2025-01-15T10:35:00Z"
}
```

---

### 4. Get All Posts

```bash
curl http://localhost:8080/api/posts
```

### 5. Get Single Post

```bash
curl http://localhost:8080/api/posts/1
```

### 6. Update Post (Authenticated)

```bash
curl -X PUT http://localhost:8080/api/posts/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "title": "Updated Title",
    "content": "Updated content here."
  }'
```

### 7. Delete Post (Authenticated)

```bash
curl -X DELETE http://localhost:8080/api/posts/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

### 8. Create Comment (Authenticated)

```bash
curl -X POST http://localhost:8080/api/posts/1/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "content": "Great post! Very informative."
  }'
```

### 9. Get Comments for Post

```bash
curl http://localhost:8080/api/posts/1/comments
```

### 10. Delete Comment (Authenticated)

```bash
curl -X DELETE http://localhost:8080/api/posts/1/comments/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## 🔐 Authentication

API menggunakan **JWT (JSON Web Token)** untuk authentication.

* Register atau login untuk mendapatkan token
* Sertakan token di header untuk protected endpoints:

  ```
  Authorization: Bearer YOUR_TOKEN_HERE
  ```
* Token berlaku selama **24 jam**

---

## 🗄️ Database Schema

### Users Table

| Kolom                              | Keterangan  |
| ---------------------------------- | ----------- |
| id                                 | Primary Key |
| email                              | Unique      |
| password                           | Hashed      |
| name                               | Nama User   |
| created_at, updated_at, deleted_at | Timestamp   |

### Posts Table

| Kolom                              | Keterangan          |
| ---------------------------------- | ------------------- |
| id                                 | Primary Key         |
| title                              | Judul               |
| content                            | Isi                 |
| user_id                            | Foreign Key → users |
| created_at, updated_at, deleted_at | Timestamp           |

### Comments Table

| Kolom                              | Keterangan          |
| ---------------------------------- | ------------------- |
| id                                 | Primary Key         |
| content                            | Isi komentar        |
| user_id                            | Foreign Key → users |
| post_id                            | Foreign Key → posts |
| created_at, updated_at, deleted_at | Timestamp           |

---

## 🐛 Troubleshooting

### Database Connection Failed

* Pastikan PostgreSQL **running**
* Cek kredensial di `.env`
* Jika pakai Docker: pastikan service **postgres** berstatus *healthy*

### Port Already in Use

* Ubah `SERVER_PORT` di `.env`
* Atau stop aplikasi lain yang menggunakan port **8080**

### Migration Error

* Drop database dan buat ulang
* Atau jalankan:

  ```bash
  docker-compose down -v
  ```

---

✨ **Selesai!**
API siap digunakan di `http://localhost:8080` 🚀
