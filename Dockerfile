# Stage 1: Build the Go program
FROM golang:1.20-alpine AS build
WORKDIR /opt/go/post

# Copy the project files and build the program
COPY . .
RUN apk --no-cache add gcc musl-dev
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o post main.go

# Stage 2: Copy the built Go program into a minimal container
FROM alpine:3.14
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /opt/go/post /app/
COPY .env /app/.env

RUN chmod +x /app/post

CMD ["/app/post"]

# Build Image with command
# docker build -t nk2-msp:${version} .