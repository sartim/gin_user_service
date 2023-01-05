# Shop Admin

### Setup

    $ go mod init go-shop-api
    $ go mod tidy
    $ go get github.com/gin-gonic/gin
    $ go get gorm.io/gorm
    $ go get gorm.io/driver/postgres
    $ go get github.com/joho/godotenv

### Migration

    $ go run app/migrations/migrate.go
