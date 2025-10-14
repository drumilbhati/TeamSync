🧩 TeamSync — Collaborative Task Management API

TeamSync is a production-ready backend service for team-based project and task management.
It’s built to simulate the real-world challenges of modern SaaS backends — user auth, role-based access, async notifications, file uploads, and more.

🚀 Features

✅ User Authentication & Authorization

JWT + Refresh token flow

Role-based permissions (admin, manager, member)

Password reset & email verification

✅ Project & Task Management

CRUD operations for projects, tasks, and subtasks

Task assignment, status tracking (todo, in progress, done)

Comment system and file attachments

✅ Team Collaboration

Invite members via secure email token

Manage roles and permissions per project

✅ Notifications & Activity

Real-time updates (WebSockets or Server-Sent Events)

Email or in-app notifications for task events

Activity logs (who changed what, when)

✅ Scalable & Production-Ready

PostgreSQL database with ORM

Redis for caching and sessions

Dockerized environment for easy deployment

CI/CD pipeline (GitHub Actions)

🧱 Tech Stack
Layer	Technology
Language	Node.js (Express / Fastify) or Python (FastAPI / Django REST)
Database	PostgreSQL + Prisma / SQLAlchemy
Cache	Redis
Queue / Messaging	RabbitMQ or Kafka (optional)
Storage	AWS S3 or MinIO
Auth	JWT + Refresh Tokens
Testing	Jest (Node) / Pytest (Python)
Deployment	Docker, Render / Fly.io / Heroku
CI/CD	GitHub Actions
📂 Project Structure
teamsync/
├── src/
│   ├── config/          # environment & DB configs
│   ├── auth/            # JWT, middleware, guards
│   ├── users/           # user routes, controllers, models
│   ├── projects/        # project + task logic
│   ├── notifications/   # WebSocket & email services
│   ├── utils/           # helpers, error handling
│   ├── tests/           # test cases
│   └── app.js           # main entrypoint
├── prisma/ or migrations/
├── docker-compose.yml
├── Dockerfile
├── .env.example
└── README.md

⚙️ Setup & Installation
1. Clone the repo
git clone https://github.com/yourusername/teamsync-backend.git
cd teamsync-backend

2. Copy and configure environment variables
cp .env.example .env


Edit .env with your credentials:

DATABASE_URL=postgresql://user:password@localhost:5432/teamsync
REDIS_URL=redis://localhost:6379
JWT_SECRET=your_jwt_secret
S3_BUCKET=teamsync-files

3. Run using Docker
docker-compose up --build

4. Run locally (without Docker)
npm install
npm run dev


or

pip install -r requirements.txt
uvicorn src.app:app --reload

🧪 Testing
npm test


or

pytest

📘 API Overview
Method	Endpoint	Description	Auth
GET	/projects	Get user’s projects	✅
POST	/projects	Create a project	✅
POST	/tasks	Create a task	✅
PATCH	/tasks/:id/status	Update task status	✅
WS	/notifications	Real-time task updates	✅

(Full API docs via Swagger/OpenAPI coming soon.)

🧠 Advanced Extensions

Full-text search (PostgreSQL or Elasticsearch)

Background jobs (BullMQ / Celery)

Rate limiting and analytics

GraphQL endpoint

Admin dashboard (optional frontend)

🛠️ Development Roadmap

Phase 1: Core CRUD + JWT Auth
Phase 2: Roles & Permissions
Phase 3: Notifications + File uploads
Phase 4: Real-time collaboration
Phase 5: Docker + CI/CD Deployment

🤝 Contributing

Pull requests are welcome!
Please open an issue first to discuss major changes.

📄 License

MIT License © 2025 Drumil Bhati