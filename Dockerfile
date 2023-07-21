# Use an official Golang runtime as a parent image
FROM golang:1.15.2

# Set the working directory to /app
WORKDIR /go-app

# Copy the current directory contents into the container at /app
COPY . /go-app

# Download and install any required dependencies
RUN go mod download

RUN go env -w GOARCH=wasm GOOS=js
RUN go build -o web/app.wasm ./app

RUN go env -w GOARCH=amd64 GOOS=windows
RUN go build -o hello ./server


# Expose port 8080 for incoming traffic
EXPOSE 8080

# Define the command to run the app when the container starts
CMD ["/go-app/hello"]
