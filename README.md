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
