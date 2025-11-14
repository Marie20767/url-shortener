## Implementation of a lightweight URL shortener in Go using:

- http server (using github.com/labstack/echo/v4)
- postgres (using github.com/jackc/pgx/v5)
- mongoDB (using go.mongodb.org/mongo-driver/v2)
- redis (using github.com/redis/go-redis/v9)
- cron (using github.com/robfig/cron/v3)

#### App architecture

![high-level system architecture](apparchitecture.png)

#### Setup

```make dep```

#### Run with docker:

- Start the app:
```docker-compose up -d```

#### Run without docker (DBs are still run via docker):

- Start the app:
```make start```

- Start the database:
```make start-db```

- Stop the database:
```make stop-db```

#### Lint:
- Run linter:
```make lint```

- Fix issues:
```make lint-fix```