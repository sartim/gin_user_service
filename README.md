# User Service

User service running on Gin Framework.

[![Language](https://img.shields.io/badge/language-go-green.svg)](https://github.com/sartim/gin_shop_api)
![build](https://github.com/sartim/gin_shop_api/actions/workflows/master.yml/badge.svg)

### Setup

###### Make migrations
    $ go run ./cmd main.go --action=create_tables 

###### Drop tables
    $ go run ./cmd main.go --action=drop_tables 

###### Create super user

    $ go run ./cmd main.go --action=create_super_user 

###### Running unittests

    $ go test -v ./tests

##### Running coverage

    $ go test -v ./tests -coverprofile=coverage.out
    $ go tool cover -html=coverage.out

###### Run server
    $ go run main.go --action=run_server 


### Install requirements
    $ go mod init gin-shop-api
    $ go mod tidy
    $ go get github.com/gin-gonic/gin
    $ go get gorm.io/gorm
    $ go get gorm.io/driver/postgres
    $ go get github.com/joho/godotenv
    $ go get github.com/golang-jwt/jwt/v4
    $ go get github.com/google/uuid
    $ go get github.com/stretchr/testify/assert
