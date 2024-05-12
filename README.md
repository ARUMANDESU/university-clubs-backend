# UCMS API Gateway 

## Overview
This service is part of the University Clubs Management application, focusing on connecting api endpoints with other microservices. 


## Technologies Used
- Go
- gRPC
- Gin
- Docker
- [Taskfile](https://taskfile.dev/)

## Getting Started
### Prerequisites
- Go version 1.22
- Docker 4.29.0

### Other microservices
- [User](https://github.com/ARUMANDESU/uniclubs-user-service)
- [Club](https://github.com/ARUMANDESU/uniclubs-club-service)
- [Posts](https://github.com/ARUMANDESU/uniclubs-posts-service)
- [Notification](https://github.com/ARUMANDESU/uniclubs-notification-service)

### Protofiles
- [Protofiles](https://github.com/ARUMANDESU/uniclubs-protos)

### Installation
Clone the repository:
   ```bash
   git clone https://github.com/ARUMANDESU/university-clubs-backend.git
   cd university-clubs-backend
   go mod download
   ```

   
### Configuration
The Uniclubs API Gateway  requires a configuration file to specify various settings like service-specific parameters, other microservices address etc. 
Depending on your environment (development, test, or production), different configurations may be needed.

#### Configuration Files
- `dev.yaml`: Contains configuration for the development environment.
- `test.yaml`: Used for the test environment.
- `local.yaml`: Configuration for local development.

#### Setting Up Configuration
```dotenv
# Example configuration snippet
ENV=dev
SHUTDOWN_TIMEOUT=10s
JWT_SECRET=hart_secret_key

HTTP_ADDRESS=localhost:5000
HTTP_TIMEOUT=5s
HTTP_IDLE_TIMEOUT=3s

USER_SERVICE_ADDRESS=localhost:44044
USER_SERVICE_TIMEOUT=10s
USER_SERVICE_RETRIES_COUNT=2
CLUB_SERVICE_ADDRESS=localhost:44045
CLUB_SERVICE_TIMEOUT=10s
CLUB_SERVICE_RETRIES_COUNT=2

MICROSOFT_OIDC_SECRET=
MICROSOFT_OIDC_AUTHORITY=
MICROSOFT_OIDC_CLIENT_ID=

AWS_REGION=ap-northeast-1
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
```

### Running the Service
After setting up the database and configuring the service, you can run it as follows:
```bash
go run cmd/user-server/main.go
```

Or use the provided Taskfile to run the service:
```bash
task run:enviroment
 ```
#or
```bash
task env
 ```

