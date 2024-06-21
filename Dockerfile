# Dockerfile

# Use the official Golang image as the base image
FROM golang:1.22.4

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code to the working directory
COPY . .

# Build the Go application
RUN make build

# Expose the port the application will run on
EXPOSE 8080

# Command to run the application
CMD ["go build -o keryx && ./keryx"]
