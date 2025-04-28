# Lingvogramm Backend

This project is a telegram web apps application for a foreign language tutor that allows managing schedules, students, lessons, and other aspects of the tutor's work.

## Technologies Stack 
* <b>Golang:</b> Programming language.
* <b>Fiber v3:</b> Fast and lightweight web framework for Golang.
* <b>Redis:</b> In-memory database for caching and session management.
* <b>PostgreSQL:</b> Relational SQL database for storing data.
* <b>Docker:</b> Containerization for simplifying deployment and management of the application.
* <b>Kubernetes (K8s):</b> Container orchestration for managing deployment and scaling.
* <b>Golangci-lint:</b> Tool for static code analysis in Golang.
* <b>Unit Tests:</b> Unit tests for verifying individual components of the application.
* <b>Integration Tests:</b> Integration tests for verifying the interaction of various components of the application.
* <b>Mockery:</b> Used for generating mocks for testing.

## Installation and Setup 
### Prerequisites
* Docker 
* Docker Compose
* Kubernetes (if used for deployment)

### Installation Steps
1. <b>Clone the repository:</b> 
```
git clone <repository-url>
cd <repository-directory>
```

2. <b>Create and configure the .env file:</b> <br>
Create a <b>.env</b> file in the root directory of the project and add the necessary environment variables:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
REDIS_HOST=localhost
REDIS_PORT=6379
```

3. <b>Start the containers using Docker Compose:</b>
```
docker-compose up --build
```

4. Configure Kubernetes (if used): <br>
Create and apply the Kubernetes manifests for deploying the application:
```
kubectl apply -f k8s/
```

### Testing
#### Unit Tests
To run unit tests, execute: 
```
go test ./... -v
```

#### Integration Tests
To run integration tests, execute: 
```
go test ./tests -v
```

### Linting
To check the code for compliance with standards, use Golangci-lint: 
```
golangci-lint run
```
or
```
make lint
```