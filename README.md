# TeamSync

TeamSync is a robust backend service built with **Go (Golang)** for managing collaborative team workflows. It features a complete system for users, teams, memberships, tasks, and comments, backed by a **PostgreSQL** database for relational data and **Redis** for high-performance caching/session management.

The project follows a clean, modular architecture, utilizing `gorilla/mux` for routing and separating concerns into controllers, models, and data access layers (store).

---

## Features

*   **Authentication:** Secure registration and login with JWT-based session management.
*   **User Management:** Create, view, update, and delete user profiles.
*   **Team & Membership:** Create teams, assign leaders, and manage team members with specific roles.
*   **Task Management:** Full lifecycle management for tasks (Create, Read, Update, Delete) with statuses (Todo, In Progress, Done) and priorities.
*   **Comments:** Collaboration features allowing users to add comments to specific tasks.
*   **Performance:** Redis integration for optimized data handling.

---

## Prerequisites

Before running the project, ensure you have the following installed:

*   **[Go](https://go.dev/doc/install)** (version 1.22 or later)
*   **[PostgreSQL](https://www.postgresql.org/download/)**
*   **[Redis](https://redis.io/docs/install/install-redis/)**

---

## Installation & Setup

1.  **Clone the Repository**
    ```bash
    git clone https://github.com/drumilbhati/teamsync.git
    cd teamsync
    ```

2.  **Install Dependencies**
    ```bash
    go mod tidy
    ```

3.  **Database Configuration**
    *   **PostgreSQL:** Create a database named `teamsync` and run your schema migration scripts.
        ```sql
        CREATE DATABASE teamsync;
        ```
    *   **Redis:** Ensure your Redis server is running (default: `localhost:6379`).

4.  **Environment Variables**
    Create a `.env` file in the root directory and populate it with your configuration:

    ```ini
    # Server Configuration
    PORT=8080

    # PostgreSQL Configuration
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=your_postgres_user
    DB_PASSWORD=your_postgres_password
    DB_NAME=teamsync

    # Redis Configuration
    REDIS_ADDR=localhost:6379
    REDIS_PASSWORD=
    REDIS_DB=0
    ```

---

## Running the Application

Start the server using:

```bash
go run main.go
```

The server will start (default port: `8080`).

---

## Deployment (AWS & Docker)

For production deployment, this project is optimized to run using **Docker Compose**.

### Local Production Build
```bash
docker-compose up -d --build
```

### AWS EC2 Deployment
For detailed AWS setup (Security Groups, Swap Space, and Nginx proxying), refer to the [AWS Deployment Guide](./AWS_DEPLOYMENT_GUIDE.md).

Quick start on EC2 (Ubuntu 24.04):
1. **Prepare Server:** Install Docker and Docker Compose.
2. **Transfer Code:** Use `rsync` to upload the project.
3. **Configure Environment:** Create a production `.env` file on the server.
4. **Launch:**
   ```bash
   sudo docker-compose up -d --build
   ```

---

## API Endpoints

### Authentication
*   `POST /auth/register` - Register a new user
*   `POST /auth/login` - Login and receive a JWT
*   `POST /auth/verify` - Verify user email

### Users (Protected)
*   `GET    /api/users` - List all users
*   `GET    /api/user/{id}` - Get user details
*   `PUT    /api/user/{id}` - Update user details
*   `DELETE /api/user/{id}` - Delete a user

### Teams (Protected)
*   `POST   /api/team` - Create a new team
*   `GET    /api/team?user_id={id}` - Get teams for a specific user
*   `GET    /api/team/{id}` - Get specific team details
*   `PUT    /api/team/{id}` - Update a team
*   `DELETE /api/team/{id}` - Delete a team

### Members (Protected)
*   `POST   /api/member` - Add a member to a team
*   `GET    /api/member?team_id={id}` - Get all members of a team
*   `GET    /api/member/{id}` - Get specific membership details
*   `PUT    /api/member/{id}` - Update membership role
*   `DELETE /api/member/{id}` - Remove a member

### Tasks (Protected)
*   `POST   /api/task` - Create a new task
*   `GET    /api/task?team_id={id}` - Get all tasks for a team
*   `GET    /api/task/{id}` - Get specific task details
*   `PUT    /api/task/{id}` - Update a task (status, assignee, etc.)
*   `DELETE /api/task/{id}` - Delete a task

### Comments (Protected)
*   `POST   /api/comment` - Add a comment to a task
*   `GET    /api/comment/{task_id}` - Get all comments for a specific task