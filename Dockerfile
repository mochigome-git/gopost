# Stage 1: Build the Go program
FROM golang:1.21.1-alpine AS build
WORKDIR /opt/go/post

# Copy the project files and build the program
COPY . .
RUN apk --no-cache add gcc musl-dev
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o post main.go

# Stage 2: Copy the built Go program into a minimal container
FROM alpine:3.18
RUN apk --no-cache add ca-certificates

# Copy the Go binary from the first stage
COPY --from=build /opt/go/post/post /app/post

RUN chmod +x /app/post

CMD ["/app/post"]


# Build Image with command
# docker build -t post:${version} .