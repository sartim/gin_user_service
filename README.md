# Shop Admin

### Setup

###### Enable uuid extension for postgres

    $ CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

###### Make migrations
    $ go run main.go --action=make_migrations 

###### Drop tables
    $ go run main.go --action=drop_tables 

###### Create super user

    $ go run main.go --action=create_super_user 

###### Running unittests

    $ go test -v ./tests

###### Run server
    $ go run main.go --action=run_server 


### Install requirments
    $ go mod init gin-shop-api
    $ go mod tidy
    $ go get github.com/gin-gonic/gin
    $ go get gorm.io/gorm
    $ go get gorm.io/driver/postgres
    $ go get github.com/joho/godotenv
    $ go get github.com/golang-jwt/jwt/v4
    $ go get github.com/google/uuid
    $ go get github.com/stretchr/testify/assert
