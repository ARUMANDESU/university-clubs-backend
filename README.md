# UCMS API Gateway 

## Overview
This service is part of the Univercity Club Management application, focusing on connecting api endpoints with other microservices. 


## Technologies Used
- Go
- gRPC
- Gin
- Docker
- [Taskfile](https://taskfile.dev/)

## Getting Started
### Prerequisites
- Go version 1.21.4
- Docker 4.26.1

### Connected microservices
- [User](https://github.com/ARUMANDESU/uniclubs-user-service)
- [Notification](https://github.com/ARUMANDESU/uniclubs-notification-service)
- [Protofiles](https://github.com/ARUMANDESU/uniclubs-protos)

### Installation
Clone the repository:
   ```bash
   git clone https://github.com/ARUMANDESU/university-clubs-backend.git
   cd university-clubs-backend
   go mod download
   ```

   
### Configuration
The Uniclubs API Gateway  requires a configuration file to specify various settings like service-specific parameters, other microsrevice address and etc. . 
Depending on your environment (development, test, or production), different configurations may be needed.

#### Configuration Files
- `dev.yaml`: Contains configuration for the development environment.
- `test.yaml`: Used for the test environment.
- `local.yaml`: Configuration for local development.

#### Setting Up Configuration
1. Choose the appropriate configuration file based on your environment.
2. Update the file with your specific settings, such as database connection strings, port numbers, and any third-party service credentials.
3. Ensure the application has access to this configuration file at runtime, either by placing it in the expected directory or setting an environment variable to its path.

#### Example Configuration
Here's an example of what the configuration file might look like (refer to `dev.yaml`, `test.yaml`, or `local.yaml` for full details):

```yaml
# Example configuration snippet
env: "local"
shutdown_timeout: "10s"
http_server:
  address: "localhost:5000"
  timeout: "4s"
  idle_timeout: c
clients:
  user:
    address: "localhost:44044"
    timeout: "3s"
    retries_count: 3
```
### Configuring with environment variables
```
GIN_MODE=
ENV=   //dev | local | prod
SHUTDOWN_TIMEOUT=   //"<int>s" | "10m" | "10h"
HTTP_ADDRESS=   //"localhost:5000"
HTTP_TIMEOUT=   //"<int>s" | "10m" | "10h"
HTTP_IDLE_TIMEOUT=   //"<int>s" | "10m" | "10h"
USER_SERVICE_ADDRESS=   //"localhost:44044"
USER_SERVICE_TIMEOUT=   //"<int>s" | "10m" | "10h"
USER_SERVICE_RETRIES_COUNT=   //<int> | 3
```
## Running the Service
After configuring the service, you can run it as follows:
  ```bash
  go run cmd/main.go --config=<path to the config file>
  //or with env
  go run cmd/main.go
  ```

