## Implementation of a lightweight URL shortener in Go using:

- http server (using github.com/labstack/echo/v4)
- postgres (using github.com/jackc/pgx/v5)
- mongoDB (using go.mongodb.org/mongo-driver/v2)
- redis (using github.com/redis/go-redis/v9)
- cron (using github.com/robfig/cron/v3)

## App architecture

![high-level system architecture](apparchitecture.png)

## API

### Create short URL - POST /create
Creates a new short URL from a provided long URL.
Accepts an optional expiry datetime. Returns the shortened URL containing a unique key.

#### <u>Request</u>

Body schema
  ```
  {
    "url": "string",                          // Required. Full URL to shorten.
    "expiry*": "string (RFC3339 datetime)"    // Optional. Expiration datetime.
  }
  ```

Example
  ```
  {
    "url": "https://api.example.org/v1/users/profile/987654321/details/gb/all",
    "expiry*": "2025-11-08T02:40:25+00:00"
  }
  ```

#### <u>Response</u>

Body schema
  ```
  {
    "url": "string" // The generated short URL
  }
  ```

Example
  ```
  {
    "url": "https://${API_DOMAIN}/j1ll5xx8"
  }
  ```

### Resolve original URL - GET/{key}
Retrieves and redirects to the original long URL using the short URL key as a path parameter.

#### <u>Request</u>

  `GET https://${API_DOMAIN}/j1ll5xx8`

## Development

### Setup

#### 1. Install dependencies:
```make dep```

#### 2. Set up environment files:
Create the following env files (you can use `.env.example` as a reference for the required variables):
- `.env` - your default local environment variables
- `.env.docker` - overrides for Docker (e.g. using @postgres as the DB host instead of localhost)
- `prod/.env` - production environment variables

### Run the app

#### 1. With docker:

`docker compose up -d`

#### 2. Without docker:

- Start the app:
`make start`

- Start databases (still run via docker):
`make start-dbs`

- Stop databases (still run via docker):
`make stop-dbs`

### Test

`make test`

### Lint
- Run linter:
`make lint`

- Fix issues:
`make lint-fix`

## Deployment

Before deploying, add the following secrets to your GitHub repository:

- `DOCKER_USER` – your Docker Hub username
- `DOCKER_PASSWORD` – your Docker Hub password

These secrets are required by the release workflow to push your application’s Docker image to Docker Hub.``