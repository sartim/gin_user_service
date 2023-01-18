# Shop Admin

### Setup

    $ CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

    $ go mod init go-shop-api
    $ go mod tidy
    $ go get github.com/gin-gonic/gin
    $ go get gorm.io/gorm
    $ go get gorm.io/driver/postgres
    $ go get github.com/joho/godotenv
    $ go get github.com/golang-jwt/jwt/v4
    $ go get github.com/google/uuid

### Migration

    $ go run app/migrations/migrate.go
