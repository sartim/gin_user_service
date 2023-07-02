FROM ubuntu:22.04

ARG DB_URL

ENV DB_URL=$DB_URL

# Set the timezone
ENV TZ=America/New_York
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Update and install necessary packages
RUN apt update
RUN apt install -y redis-server wget
RUN apt install -y postgresql postgresql-contrib
RUN apt-get install -y  bison mercurial
RUN wget https://go.dev/dl/go1.19.10.linux-amd64.tar.gz
RUN tar -C . -xzf go1.19.10.linux-amd64.tar.gz

# Copy the application code
WORKDIR /app
COPY . /app

RUN export GOROOT=$HOME/go
RUN export PATH=$PATH:$GOROOT/bin
RUN /go/bin/go mod init gin-shop-api
RUN /go/bin/go mod tidy
# Build the app to binary
RUN mkdir build && /go/bin/go build -o ./build/user-service ./cmd/main.go   

# Expose port 8000 for the app
EXPOSE 8000

# Start the app
CMD ["/app/build/user-service", "--action=run-server"]
