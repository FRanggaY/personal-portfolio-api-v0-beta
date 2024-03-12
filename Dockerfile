# Base image
FROM golang:1.21

# Set working directory
WORKDIR /app

# copy
COPY . .

# Set the command to run your app
CMD go run main.go
