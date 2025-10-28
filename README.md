# Task Management Backend

A scalable and cleanly architected RESTful API for managing team-based projects and tasks — featuring real-time collaboration, role-based access control, and WebSocket notifications.

[![Go Version](https://img.shields.io/badge/Go-1.23%2B-00ADD8?logo=go&logoColor=white)](https://golang.org/)
[![Fiber](https://img.shields.io/badge/Framework-Fiber%20v2-2C8EBB?logo=fiber&logoColor=white)](https://gofiber.io/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)


## Key Features

- **Authentication & Authorization** – Register, login, and secure routes using JWT.  
- **Project & Board Management** – Create, manage, and order boards within projects.  
- **Task Management** – Add tasks with assignees, due dates, color labels, and file attachments.  
- **Team Collaboration** – Invite members to projects and manage their roles.  
- **Nested Comments** – Support for threaded comments up to 2 levels (similar to TikTok/Instagram).  
- **Real-time Notifications** – WebSocket-based instant updates (new comments, invites, etc).  
- **Activity Log** – Automatic project activity tracking.  
- **File Uploads** – Integrated with Cloudinary for task attachments.  

> **Note:** The invitation system is currently in development mode — tokens are returned in API responses for testing purposes.  
> You can enable real email integration in a production environment.

## Tech Stack

- **Language:** Go (Golang) 1.23+  
- **Framework:** Fiber v2  
- **ORM:** GORM  
- **Database:** PostgreSQL  
- **File Storage:** Cloudinary  
- **Authentication:** JWT (JSON Web Tokens)  
- **Real-time:** WebSocket (`github.com/gofiber/websocket/v2`)  
- **Architecture:** Clean Architecture (Repository–Service–Handler)  
- **Error Handling:** Centralized custom sentinel errors with DTOs  

## Prerequisites

Before running the application, make sure you have the following installed:

- [Go](https://golang.org/dl/) v1.23 or later  
- [PostgreSQL](https://www.postgresql.org/download/) v12+  
- [Git](https://git-scm.com/downloads)

### Environment Configuration

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Fill in your .env file: 
    ```bash
    # Database
    DB_HOST=localhost
    DB_USER=postgres
    DB_PASSWORD=your_password
    DB_NAME=task_management
    DB_PORT=5432

    # Cloudinary
    CLOUDINARY_CLOUD_NAME=your_cloud_name
    CLOUDINARY_API_KEY=your_api_key
    CLOUDINARY_API_SECRET=your_api_secret

    # JWT
    JWT_SECRET=your_strong_jwt_secret_here
    ```

## Project Structure

```bash
task-management-backend/
├── cmd/
│   └── main.go                 # Application entry point
├── config/                     # Database & Cloudinary configuration
├── internal/
│   ├── dto/                    # Data Transfer Objects
│   ├── errors/                 # Centralized sentinel errors
│   ├── handlers/               # HTTP request handlers
│   ├── middlewares/            # Authentication & logging middleware
│   ├── models/                 # GORM models
│   ├── repository/             # Database access layer
│   ├── routes/                 # HTTP route definitions
│   ├── services/               # Business logic layer
│   ├── utils/                  # Utility/helper functions
│   └── websocket/              # WebSocket hub & handlers
├── .env.example                # Example environment variables
├── go.mod
├── go.sum
└── README.md
```

## Usage Examples

1. **Run the Application**
    ```bash
    # Install dependencies
    go mod tidy

    # Start the server
    go run cmd/main.go
    ```

The server will run at: http://localhost:8080

2. **Example API Flow (via Postman)**

    1. Register a User

*Request*
```json
POST v1/api/auth/register
Content-Type: application/json
{
    "name": "John Doe",
    "email": "johndoe@example.com",
    "password": "securepassword123"
}
```

*Response*
```json
{
    "message": "user registered successfully"
}
```

    2. Login & Get JWT Token

*Request*
```json
POST v1/api/auth/login
Content-Type: application/json
{
    "email": "johndoe@example.com",
    "password": "securepassword123"
}
```

*Response*
```json
{
    "message": "login successfully",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9…"
}
```

    3. Create a Project

*Request*
```json
POST v1/api/projects
Authorization: Bearer <your-jwt-token>
Content-Type: application/json
{
    "name": "Project Alpha",
    "description": "Task management project"
}
```

*Response*
```json
{
    "success": true,
    "message": "Project created successfully",
    "data": {
        "id": "a4f3e4dc-62ie-478c-9a82-1f8e1e54awsz",
        "name": "Project Alpha",
        "description": "Task management project",
        "owner_id": "c8a1bda1-7d4s-48a5-8810-18e2516d7b65",
        "created_at": "2025-10-28T11:15:09.9293128+07:00",
        "updated_at": "2025-10-28T11:15:09.9293128+07:00"
    }
}
```

    4. Create a Task in a Board

*Request*
```json
POST v1/api/boards/<board_id>/tasks
Authorization: Bearer <your-jwt-token>
Content-Type: application/json
{
    "title": "Implement WebSocket",
    "description": "integrate websocket to project",
    "priority": "high",
    "due_date": "2025-12-31",
    "assignee_id": "c8a1bda1-7d4s-48a5-8810-18e2516d7b65",
    "labels": [
        {
            "name": "backend",
            "color": "blue"
        }
    ]
}
```

*Response*
```json
{
    "success": true,
    "message": "Task created successfully",
    "data": {
        "id": "37d1d00a-399b-4aa6-b861-005b9f274220",
        "board_id": "866ea35e-3a90-427a-8ade-9a93454e43aa",
        "title": "integrate websocket to project",
        "description": "Selesaikan Cbackend UD task dengan soft delete",
        "priority": "high",
        "due_date": "2025-12-31",
        "assignee_id": "c8a1bda1-7d4s-48a5-8810-18e2516d7b65",
        "created_by": "4d6a1748-2f27-4432-9f0f-87969c1fc660",
        "created_at": "2025-10-28T11:28:57.844025+07:00",
        "updated_at": "2025-10-28T11:28:57.844025+07:00"
    }
}
```

    5. Connect WebSocket for Notifications

        ```bash
        // In browser console
        const ws = new WebSocket("ws://localhost:8080/v1/api/ws/notifications?  user_id=<your-user-id>");
        ws.onmessage = (event) => console.log("Notification:", JSON.parse(event.data));
        ```

## License
This project is licensed under the [MIT License](./LICENSE).  
© 2025 Muhammad Farhaan — All rights reserved.