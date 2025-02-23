# Liive Application

A modern microservices-based application with React frontend and Go backend services.

## Project Structure

```
.
├── backend/
│   ├── shared/              # Shared Go libraries
│   ├── liive-ws-api/        # WebSocket API service
│   ├── liive-rest-api/      # REST API service
│   ├── liive-message-storer/# Message storage service
│   ├── liive-auth/          # Authentication service
│   └── liive-user-manager/  # User management service
├── frontend/
│   └── liive-app/          # React frontend application
└── docker-compose.yml      # Docker composition file
```

## Prerequisites

- Docker and Docker Compose
- Go 1.21 or later
- Node.js 18 or later
- VSCode with Remote Containers extension (for development)

## Getting Started

### Development Environment Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/liive.git
   cd liive
   ```

2. Start the development environment using VS Code Remote Containers:
   - Open VS Code
   - Press F1 and select "Remote-Containers: Open Folder in Container"
   - Select the project directory

3. Start the services:
   ```bash
   docker-compose up -d
   ```

### Frontend Development

1. Navigate to the frontend directory:
   ```bash
   cd frontend/liive-app
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm start
   ```

The frontend will be available at http://localhost:3000

### Backend Development

Each microservice can be developed independently. To work on a service:

1. Navigate to the service directory:
   ```bash
   cd backend/<service-name>
   ```

2. Install Go dependencies:
   ```bash
   go mod tidy
   ```

3. Run the service:
   ```bash
   go run cmd/server/main.go
   ```

### Database Migrations

The application uses GORM for database operations and automatic migrations.
Migrations are handled automatically when the services start.

To create a new migration:

1. Add new fields or modify existing ones in the model structs
2. The changes will be automatically applied when the service restarts

### Testing

To run tests for backend services:

```bash
# Run tests for a specific service
cd backend/<service-name>
go test ./...

# Run tests for all services
cd backend
go test ./...
```

To run frontend tests:

```bash
cd frontend/liive-app
npm test
```

## API Documentation

### Authentication Service (Port 8082)
- POST /auth/login - Login endpoint
- POST /auth/register - Registration endpoint
- GET /auth/verify - Token verification endpoint

### User Manager Service (Port 8083)
- GET /users - List users
- GET /users/{id} - Get user details
- PUT /users/{id} - Update user
- DELETE /users/{id} - Delete user

### REST API Service (Port 8081)
- Various REST endpoints for application functionality

### WebSocket API Service (Port 8080)
- WS /ws - WebSocket connection endpoint

## Architecture

The application follows a microservices architecture with the following components:

- **Frontend**: React application for the user interface
- **Backend Services**:
  - **Auth Service**: Handles authentication and authorization
  - **User Manager**: Manages user accounts and profiles
  - **REST API**: Provides REST endpoints for application functionality
  - **WebSocket API**: Handles real-time communication
  - **Message Storer**: Manages message persistence and retrieval

## Contributing

1. Create a feature branch
2. Make your changes
3. Run tests
4. Submit a pull request

## License

[Your License Here] 