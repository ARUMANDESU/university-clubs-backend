<a id="readme-top"></a>

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://astanait.edu.kz/">
    <img src="https://static.tildacdn.pro/tild3764-6633-4663-b138-303730646233/aitu-logo__2.png" alt="Logo" height="80">
  </a>
  <h3 align="center">AITU UCMS API Gateway</h3>
  <p align="center">
    This service is part of the University Clubs Management application, focusing on connecting API endpoints with other microservices.
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#about-the-project">About The Project</a></li>
    <li><a href="#other-microservices">Other Microservices</a></li>
    <li><a href="#protofiles">Protofiles</a></li>
    <li><a href="#technologies-used">Technologies Used</a></li>
    <li><a href="#getting-started">Getting Started</a></li>
    <ul>
      <li><a href="#prerequisites">Prerequisites</a></li>
      <li><a href="#installation">Installation</a></li>
      <li><a href="#configuration">Configuration</a></li>
    </ul>
    <li><a href="#running-the-service">Running the Service</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project

This API Gateway service is a critical component of the University Clubs Management System (UCMS) at AITU. It serves as the intermediary, routing requests to various backend microservices.

![High Level Architecture](https://github.com/user-attachments/assets/dda8fb46-1233-4d9f-a816-38d9ef5050fd)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<details>
  <summary><h2>Demo Images<h2></summary>
  <img width="5256" height="2586" alt="image" src="https://github.com/user-attachments/assets/24594266-8621-4c11-ab68-7355feb107bd" />
<img width="3751" height="2586" alt="image" src="https://github.com/user-attachments/assets/ed5dca7b-97ce-4169-a5a4-4303479efe2a" />

  <img width="5262" height="2586" alt="image" src="https://github.com/user-attachments/assets/2803043c-0ad6-4caf-9730-bde7d74d7534" />
  <img width="5262" height="2586" alt="image" src="https://github.com/user-attachments/assets/3041c824-aafe-44b7-8c4a-68586178a71b" />
  <img width="5265" height="2588" alt="image" src="https://github.com/user-attachments/assets/6e098f42-8d19-4a51-9610-fc2c8e780d35" />
  <img width="4273" height="2586" alt="image" src="https://github.com/user-attachments/assets/6f1ee1ed-479a-4834-9987-41d7a900b723" />

  <img width="5265" height="2588" alt="image" src="https://github.com/user-attachments/assets/4b683df9-8bb8-4d03-9136-eea3ad8999cf" />


  <img width="5262" height="2586" alt="image" src="https://github.com/user-attachments/assets/738dd941-36a6-4379-8bdb-23bafa0ebbb7" />
  <img width="5262" height="2586" alt="image" src="https://github.com/user-attachments/assets/5fc3ebb1-3556-4fd5-9af8-99e22170a6d1" />
  <img width="5262" height="2586" alt="image" src="https://github.com/user-attachments/assets/5e6ede2b-42fb-4540-b674-153dd8a9734c" />
  <img width="5262" height="2586" alt="image" src="https://github.com/user-attachments/assets/dc7c8414-7387-4a3e-b44c-85a68bfe23db" />

  <img width="5964" height="3146" alt="image" src="https://github.com/user-attachments/assets/c8b10825-43e5-48b8-852e-0271b75cb79f" />
  <img width="4923" height="2587" alt="image" src="https://github.com/user-attachments/assets/5d7df489-93a3-459e-89b2-a496294ba332" />
  
</details>

<!-- OTHER MICROSERVICES -->
## Other Microservices

* [User Service][user-service-url]
* [Club Service][club-service-url]
* [Posts Service][posts-service-url]
* [Comments Service][comments-service-url]
* [Notification Service][notification-service-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- PROTOFILES -->
## Protofiles

* [Protofiles Repository][protofiles-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- TECHNOLOGIES USED -->
## Technologies Used

* [![Go][go-shield]][go-url]
* [![Gin][gin-shield]][gin-url]
* [![gRPC][grpc-shield]][go-url]
* [![Docker][docker-shield]][docker-url]
* [![Docker Compose][docker-compose-shield]][docker-compose-url]
* [![Taskfile][tasks-shield]][tasks-url]


<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
## Getting Started

### Prerequisites

* Go version 1.22
* Docker 4.29.0
* [Libvips](https://www.libvips.org) (required for [bimg](https://github.com/h2non/bimg))

  ```sh
  go version
  docker --version
  vips --version
  ```

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/ARUMANDESU/university-clubs-backend.git
   cd university-clubs-backend
   go mod download
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Configuration

The UCMS API Gateway requires a configuration file to specify various settings such as service-specific parameters and other microservices addresses.

Create .env file:
```bash
touch .env
```

#### Example Configuration Snippet

```dotenv
# Example configuration snippet
ENV=dev
SHUTDOWN_TIMEOUT=10s
JWT_SECRET=

# HTTP
HTTP_ADDRESS=
HTTP_TIMEOUT=5s
HTTP_IDLE_TIMEOUT=3s

# Microservices
USER_SERVICE_ADDRESS=
USER_SERVICE_TIMEOUT=10s
USER_SERVICE_RETRIES_COUNT=

CLUB_SERVICE_ADDRESS=
CLUB_SERVICE_TIMEOUT=10s
CLUB_SERVICE_RETRIES_COUNT=

EVENT_SERVICE_ADDRESS=localhost:44046
EVENT_SERVICE_TIMEOUT=10s
EVENT_SERVICE_RETRIES_COUNT=2

COMMENT_SERVICE_ADDRESS=localhost:44047
COMMENT_SERVICE_TIMEOUT=10s
COMMENT_SERVICE_RETRIES_COUNT=2

# OIDC
MICROSOFT_OIDC_SECRET=
MICROSOFT_OIDC_AUTHORITY=
MICROSOFT_OIDC_CLIENT_ID=

# S3
AWS_REGION=
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- RUNNING THE SERVICE -->
## Running the Service

After setting up the database and configuring the service, you can run it as follows:

```bash
go run cmd/user-server/main.go
```

Or use the provided Taskfile to run the service:

```bash
task run:environment
```
or
```bash
task env
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[aitu-url]: https://astanait.edu.kz/
[aitu-ucms-url]: https://www.ucms.space/

<!-- Other Microservices -->
[user-service-url]: https://github.com/ARUMANDESU/uniclubs-user-service
[club-service-url]: https://github.com/ARUMANDESU/uniclubs-club-service
[posts-service-url]: https://github.com/ARUMANDESU/uniclubs-posts-service
[comments-service-url]: https://github.com/ARUMANDESU/uniclubs-comments-service
[notification-service-url]: https://github.com/ARUMANDESU/uniclubs-notification-service
[protofiles-url]: https://github.com/ARUMANDESU/uniclubs-protos

[go-url]: https://golang.org/
[docker-url]: https://www.docker.com/
[docker-compose-url]: https://docs.docker.com/compose/
[grpc-url]: https://grpc.io/
[gin-url]: https://gin-gonic.com/
[tasks-url]: https://taskfile.dev/

[go-shield]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[docker-shield]: https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white
[docker-compose-shield]: https://img.shields.io/badge/Docker_Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white
[grpc-shield]: https://img.shields.io/badge/gRPC-008FC7?style=for-the-badge&logo=google&logoColor=white
[gin-shield]: https://img.shields.io/badge/Gin-00ADD8?style=for-the-badge&logo=go&logoColor=white
[tasks-shield]: https://img.shields.io/badge/Taskfile-00ADD8?style=for-the-badge&logo=go&logoColor=white


