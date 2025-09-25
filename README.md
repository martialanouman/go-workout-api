# FEM Project - Workout Tracking API

A RESTful API service for tracking workouts and managing users, built with Go and PostgreSQL. This project provides endpoints for user registration, authentication via JWT tokens, and comprehensive workout management.

## Features

- 🔐 **User Authentication**: JWT-based authentication system with token management
- 👤 **User Management**: User registration with secure password hashing
- 🏋️ **Workout Tracking**: Full CRUD operations for workout management
- 🗄️ **Database Migrations**: Automated database schema management with Goose
- 🐳 **Docker Support**: Easy deployment with Docker Compose
- 🔒 **Security**: Input validation, middleware protection, and secure password handling

## API Endpoints

### Authentication

- `POST /api/tokens` - Create authentication token (login)
- `POST /api/tokens/revoke-all` - Revoke all tokens for authenticated user

### Users

- `POST /api/users` - Register new user

### Workouts

- `GET /api/workouts` - Get all workouts for authenticated user
- `GET /api/workouts/{id}` - Get specific workout by ID
- `POST /api/workouts` - Create new workout
- `PUT /api/workouts/{id}` - Update existing workout
- `DELETE /api/workouts/{id}` - Delete workout

## Tech Stack

- **Language**: Go 1.24.6
- **Database**: PostgreSQL 17
- **Router**: Chi v5
- **Authentication**: JWT (golang-jwt/jwt/v4)
- **Password Hashing**: bcrypt
- **Migrations**: Goose
- **Database Driver**: pq (PostgreSQL driver)

## Project Structure

```
fem_project/
├── main.go                    # Application entry point
├── compose.yml               # Docker Compose configuration
├── go.mod                    # Go module definition
├── internal/
│   ├── api/                  # HTTP handlers
│   │   ├── token_handler.go  # Authentication endpoints
│   │   ├── user_handler.go   # User registration
│   │   └── workout_handler.go # Workout CRUD operations
│   ├── app/
│   │   └── app.go           # Application initialization
│   ├── middleware/
│   │   └── middleware.go    # Authentication middleware
│   ├── routes/
│   │   └── routes.go        # Route configuration
│   ├── store/               # Data access layer
│   │   ├── database.go      # Database connection
│   │   ├── tokens.go        # Token operations
│   │   ├── user_store.go    # User operations
│   │   └── workout_store.go # Workout operations
│   ├── tokens/
│   │   └── tokens.go        # JWT utilities
│   └── utils/
│       └── utils.go         # Common utilities
└── migrations/              # Database migrations
    ├── fs.go               # Embedded migrations
    └── *.sql               # Migration files
```

## Getting Started

### Prerequisites

- Go 1.24.6 or later
- Docker and Docker Compose
- PostgreSQL 17 (if running without Docker)

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/martialanouman/femProject.git
   cd fem_project
   ```

2. **Start the database with Docker Compose**

   ```bash
   docker compose up -d db
   ```

3. **Install dependencies**

   ```bash
   go mod download
   ```

4. **Set up environment variables** (optional)

   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=workout_user
   export DB_PASSWORD=workout_password
   export DB_NAME=workout_db
   ```

5. **Run database migrations**

   ```bash
   go run main.go migrate
   ```

6. **Start the application**
   ```bash
   go run main.go
   # or with custom port
   go run main.go -port=3000
   ```

The API will be available at `http://localhost:8080` (or your specified port).

### Running with Docker Compose

To run the entire application stack:

```bash
docker compose up -d
```

This will start both the PostgreSQL database and the application.

## Configuration

The application uses the following default configuration:

- **Port**: 8080 (configurable via `-port` flag)
- **Database**: PostgreSQL connection via environment variables
- **Timeouts**:
  - Idle: 1 minute
  - Read: 10 seconds
  - Write: 30 seconds

## Database Schema

### Users Table

- `id` - Unique identifier
- `username` - Unique username (50 chars max)
- `email` - Unique email address
- `password_hash` - Hashed password
- `bio` - User biography
- `created_at`, `updated_at` - Timestamps

### Workouts Table

- `id` - Unique identifier
- `user_id` - Reference to users table
- `title` - Workout title (100 chars max)
- `description` - Workout description
- `duration_minutes` - Workout duration
- `calories_burned` - Calories burned during workout
- `created_at`, `updated_at` - Timestamps

### Tokens Table

- Authentication tokens with expiration and user association

## API Usage Examples

### Register a User

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepassword",
    "bio": "Fitness enthusiast"
  }'
```

### Login (Get Token)

```bash
curl -X POST http://localhost:8080/api/tokens \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword"
  }'
```

### Create a Workout

```bash
curl -X POST http://localhost:8080/api/workouts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Morning Run",
    "description": "5km run in the park",
    "duration_minutes": 30,
    "calories_burned": 300
  }'
```

## Development

### Running Tests

```bash
go test ./...
```

### Database Migrations

To create a new migration:

```bash
goose -dir migrations create migration_name sql
```

To apply migrations:

```bash
goose -dir migrations postgres "user=workout_user password=workout_password dbname=workout_db sslmode=disable" up
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Martial Anouman** - [GitHub](https://github.com/martialanouman)

---

_Built with ❤️ and Go_
