## Implementation of a lightweight URL shortener in Go using:

- http server (using github.com/labstack/echo/v4)
- postgres (using github.com/jackc/pgx/v5)
- mongoDB (using go.mongodb.org/mongo-driver/v2)
- redis (using github.com/redis/go-redis/v9)
- cron (using github.com/robfig/cron/v3)

## App architecture

![high-level system architecture](apparchitecture.png)

## API

### Create short URL - POST

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
    "url": "https://minurl.com/j1ll5xx8"
  }
  ```

### Resolve original URL - GET/{key}

Retrieves and redirects to the original long URL using the short URL key as a path parameter.

#### <u>Request</u>

  ```
  GET https://minurl.com/j1ll5xx8
  ```

## Development

### Setup

```make dep```

### Run with docker:

```docker compose up -d```

### Run without docker:

- Start the app:
```make start```

- Start databases (still run via docker):
```make start-dbs```

- Stop databases (still run via docker):
```make stop-dbs```

### Lint:
- Run linter:
```make lint```

- Fix issues:
```make lint-fix```