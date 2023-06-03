# User Service

User service running on Gin Framework.

[![Language](https://img.shields.io/badge/language-go-blue.svg)](https://go.dev)
[![Build Status](https://github.com/sartim/gin_user_service/workflows/build/badge.svg)](https://github.com/sartim/gin_user_service/actions/workflows/master.yml)

### Setup

###### Make migrations
    $ go run ./cmd main.go --action=create-tables 

###### Drop tables
    $ go run ./cmd main.go --action=drop-tables 

###### Create super user

    $ go run ./cmd main.go --action=create-super_user 

###### Running unittests

    $ go test -v ./tests

##### Running coverage

    $ go test -v ./tests -coverprofile=coverage.out
    $ go tool cover -html=coverage.out

###### Run server
    $ go run main.go --action=run-server 


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
