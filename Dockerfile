# Use the official Go image as the base image
FROM golang:1.17-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the application code into the container
COPY . .

# Build the Go application
RUN go build -o app

# Expose the port on which your Echo application runs
EXPOSE 8083

# Command to run your application
CMD ["./app"]
