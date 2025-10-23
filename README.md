# TeamSync

TeamSync is a backend service built with Go (Golang) for managing users, teams, and team memberships. It provides a RESTful API for all standard CRUD (Create, Read, Update, Delete) operations, built on a robust relational model with a PostgreSQL database.

This service is built using a clean, modular architecture with `gorilla/mux` for routing, `pq` as the PostgreSQL driver, and a clear separation of concerns between `controllers`, `models`, and the `store` (data access layer).

---

## Features

* **User Management:** Full CRUD operations for platform users.
* **Team Management:** Full CRUD operations for teams.
* **Membership Management:** A relational-driven system to associate users with teams and define their roles.
* **Data Integrity:** Utilizes foreign key constraints in PostgreSQL to ensure relationships (like `team_leader` and `members`) are valid.
* **Configuration-driven:** Uses a `.env` file for all database and server configurations.

---

## Prerequisites

Before you begin, ensure you have the following installed on your local machine:
* [Go](https://go.dev/doc/install) (version 1.18 or later)
* [PostgreSQL](https://www.postgresql.org/download/)
* A tool to run SQL scripts (like `psql` or a GUI like DBeaver)

---

## Installation and Setup

1.  **Clone the Repository**
    ```bash
    git clone [https://github.com/drumilbhati/teamsync.git](https://github.com/drumilbhati/teamsync.git)
    cd teamsync
    ```

2.  **Install Dependencies**
    This will download `gorilla/mux`, `godotenv`, and the `pq` driver.
    ```bash
    go mod tidy
    ```

3.  **Database Setup**
    * Ensure your PostgreSQL server is running.
    * Create a database for the project.
        ```sql
        CREATE DATABASE teamsync;
        ```
    * **Important:** Apply your database schema. You must run your `schema.sql` file (which contains your `CREATE TABLE` statements) against the new database.
        ```bash
        # Example command:
        psql -d teamsync -f path/to/your/schema.sql
        ```

4.  **Configure Environment Variables**
    Create a `.env` file in the root of the project:
    ```bash
    touch .env
    ```
    Add your database and server configuration to the `.env` file.

    **.env Example:**
    ```ini
    # Server Port
    PORT=8080

    # PostgreSQL Database Configuration
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=your_postgres_user
    DB_PASSWORD=your_postgres_password
    DB_NAME=teamsync
    ```

---

## Running the Application

Once your database is running and your `.env` file is configured, you can start the server:

```bash
go run main.go
