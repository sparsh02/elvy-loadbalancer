## Features

- **Rate Limiting:** Control the rate of requests to your backend servers.
- **Dynamic Configuration:** Add and remove backend servers at runtime without restarting the load balancer.

## Getting Started

To get started with Elvy, you need to configure it properly and run it using Docker.

### Docker Setup

1. **Build the Docker Image**

   Build the Docker image using the provided `Dockerfile`:

   ```sh
   docker build -t elvy:latest .

    Publish the Docker Image

    Tag and push the image to Docker Hub:

```

docker tag elvy:latest yourusername/elvy:latest
docker push yourusername/elvy:latest
```

 **Run Elvy with Docker**

You can run Elvy and provide a config.yaml file at runtime using Docker Compose or Docker run commands.

Using Docker Compose

Create a docker-compose.yml file:

```

yaml

version: '3.8'

services:
  elvy:
    image: yourusername/elvy:latest
    ports:
      - "8080:8080"
    volumes:
      - ./config.yaml:/root/config.yaml  # Mount the configuration file
```
Run Docker Compose:

```
docker-compose up
```

** Using Docker Run Command**

Run the container and mount the configuration file:

```

    docker run -p 8080:8080 -v /path/to/config.yaml:/root/config.yaml sparsh02/elvy:latest
```

*Configuration*

The configuration for Elvy is managed through a config.yaml file. Below is an explanation of the available configuration options.
Example config.yaml

```
type: round_robin
port: 8080
sticky_session: false
backends:
  - address: http://localhost:8081
    alive: true
    rate_limit:
      enabled: true
      requests_per_minute: 2
  - address: http://localhost:8082
    alive: true
    rate_limit:
      enabled: true
      requests_per_minute: 2
```

Configuration Options

    type: Load balancing algorithm to use. Options are round_robin, least_conn, and ip_hash. Defaults to round_robin.
    port: The port on which Elvy will listen. Defaults to 8080.
    sticky_session: Whether to use sticky sessions. Set to true or false. Defaults to false.
    backends: List of backend servers.

Backend Server Options

    address: The address of the backend server. E.g., http://localhost:8081.
    alive: Indicates if the backend server is considered alive. Set to true or false.
    rate_limit: Configuration for rate limiting.
        enabled: Whether rate limiting is enabled. Set to true or false.
        requests_per_minute: Number of allowed requests per minute.

API Endpoints

Elvy provides API endpoints for dynamically managing backend servers.
Add Server

Endpoint: POST /servers

Request Body:
```
{
  "address": "http://localhost:8083",
  "alive": true,
  "rate_limit": {
    "enabled": true,
    "requests_per_minute": 5
  }
}
```

Remove Server

Endpoint: DELETE /servers

Request Body:
```
{
  "address": "http://localhost:8083"
}
```

