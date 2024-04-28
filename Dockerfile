FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Download dependencies
RUN go mod download

# Add debugging statement
RUN ls -la

# Build the Go application
RUN go build -o Rivr-Engine .

# Set the entry point for the container
CMD ["./Rivr-Engine"]