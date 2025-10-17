# Event Booking REST API

## Project Description
Implementing a simple REST API in Go.

📚 *This project was created for educational purposes.*

A REST API for an event booking system that allows users to create and register for events.

## Functional Requirements

### User Management
| Endpoint | Description |
|---|---|
| `POST /signup` | Register a new user |
| `POST /login`  | Log in a user and receive a JWT token |
- Email and password validation
- Password hashing

### Event Management
| Endpoint | Description |
|---|---|
| `GET /events` | Get a list of all events (accessible to everyone) |
| `GET /events/:id` | Get a specific event by its ID |
| `POST /events` | Create a new event (authorized users only) |
| `PUT /events/:id` | Update an event (event creator only) |
| `DELETE /events/:id` | Delete an event (event creator only) |

### Booking System
| Endpoint | Description |
|---|---|
| `POST /events/:id/register` | Register for an event (authorized users only) |
| `DELETE /events/:id/register` | Cancel a registration for an event |
- Uniqueness check for registration (a user cannot register for the same event twice)

## Technical Requirements

### Technology Stack
- **Programming Language:** Go 1.24
- **HTTP Framework:** Gin
- **Database:** SQLite (for development simplicity)
- **ORM/Query Builder:** Raw SQL queries
- **Authentication:** JWT tokens

### Database
- **users** - users table
- **events** - events table
- **registrations** - registrations table (M2M relationship between users and events)

### Middleware
- **Authentication middleware** - verifies JWT tokens for protected routes
- **CORS middleware** - to support cross-origin requests
- **Logging middleware** - logs requests

### Project Structure
```
event-booking-api/
├── main.go
├── config/
│   └── config.go
├── models/
│   ├── user.go
│   └── event.go
├── repositories/
│   ├── users.go
│   └── events.go
├── services/
│   ├── users.go
│   └── events.go
├── routes/
│   ├── routes.go
│   ├── users.go
│   ├── events.go
│   └── register.go
├── middlewares/
│   └── auth.go
├── utils/
│   ├── jwt.go
│   └── hash.go
├── db/
│   └── db.go
├── api_test.go
├── api.db (SQLite)
├── go.mod
├── go.sum
└── README.md
```

### Running the Project
```bash
go run main.go
```

The API will be available at: `http://localhost:8080`
