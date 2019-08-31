# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Adding Binary
ADD lunchmore /app/

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./lunchmore"]
